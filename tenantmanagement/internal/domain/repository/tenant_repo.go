package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/database"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/database/postgres"
)

type TenantRepository struct {
	DB database.DB
}

type ITenantRepository interface {
	CreateTenant(ctx context.Context, tenant *entity.Tenant) error
	GetTenantContext(ctx context.Context, tenantID value_object.TenantID, userID value_object.UserID) (*TenantContext, error)
	UpdateTenantSettings(ctx context.Context, tenantID value_object.TenantID, userID value_object.UserID, incoming entity.ThemeConfig,
	) (*TenantContext, error)
	// GetTenantById(ctx context.Context, id string) (*entity.Tenant, error)
	// GetTenantBySlug(ctx context.Context, slug string) (*entity.Tenant, error)
}

func NewTenantRepository(db database.DB) ITenantRepository {
	return &TenantRepository{DB: db}
}

// CreateTenant performs an atomic dual-insertion into the 'tenants' and 'tenant_settings' tables.
// It uses a database transaction to ensure that a tenant is never created without its
// corresponding theme configuration.
func (r *TenantRepository) CreateTenant(ctx context.Context, t *entity.Tenant) error {
	const op = "TenantRepository.CreateTenant"
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := r.DB.Tx(ctx)
	if err != nil {
		return fmt.Errorf("[%s]: %w", op, err)
	}
	defer tx.Rollback(ctx)

	tenantQuery := `INSERT INTO tenants (store_name, business_type, slug, sub_domain, domain)
                    VALUES($1, $2, $3, $4, $5) RETURNING id, store_name, sub_domain, domain, status, created_at`

	err = tx.QueryRow(ctx, tenantQuery, t.StoreName, t.BusinessType, t.Slug, t.SubDomain, t.Domain).
		Scan(&t.ID,
			&t.StoreName,
			&t.SubDomain,
			&t.Domain,
			&t.Status,
			&t.CreatedAt,
		)
	if err != nil {
		log.Printf("[%s error]: %v", op, err)
		return fmt.Errorf("%w", domain.TranslateTenantRepoError(postgres.MapPostgresError(err)))
	}

	t.Settings.TenantID = t.ID

	settingsQuery := `INSERT INTO tenant_settings (tenant_id, theme_id, config, version)
                      VALUES($1, $2, $3, $4) RETURNING updated_at`

	configToPersist := t.Settings.Config
	if configToPersist == nil {
		configToPersist = make(entity.ThemeConfig)
	}

	err = tx.QueryRow(ctx, settingsQuery, t.ID, t.Settings.ThemeID, configToPersist, t.Settings.Version).
		Scan(&t.Settings.UpdatedAt)
	if err != nil {
		log.Printf("%s error: %v", op, err)

		return fmt.Errorf("%w", domain.TranslateTenantRepoError(postgres.MapPostgresError(err)))
	}

	return tx.Commit(ctx)
}

type TenantContext struct {
	Tenant *entity.Tenant
	Role   string
}

// GetTenantContext resolves the full execution context of a user within a specific tenant.
//
// PURPOSE:
// Loads the tenant aggregate along with the caller’s role within that tenant. This enables
// downstream authorization and business logic decisions.
//
// SECURITY MODEL:
// Access is enforced at the query level via an INNER JOIN on `users_tenants`. If the
// provided userID is not associated with the tenantID, the query returns no rows. This
// prevents unauthorized access and tenant ID enumeration.
//
// DATA COMPOSITION:
// - tenants: core tenant identity (store metadata, domain, lifecycle status, timestamps)
// - tenant_settings: theme configuration (cloned from the tenant’s theme at creation and always present)
// - users_tenants: resolves the caller’s role within the tenant
//
// RETURNS:
// - TenantContext: contains the hydrated Tenant entity and the caller’s RoleID
//
// NOTE:
// - Settings.Config is guaranteed to be present (cloned on tenant creation).
// - JSONB from Postgres is scanned automatically into ThemeConfig via pgx v5.
func (r *TenantRepository) GetTenantContext(ctx context.Context, tenantID value_object.TenantID, userID value_object.UserID) (*TenantContext, error) {
	const op = "TenantRepository.GetTenantContext"

	query := `
        SELECT 
            t.id, t.store_name, t.business_type, t.slug, t.sub_domain, t.domain,
            t.status, t.created_at, t.updated_at, t.deleted_at,
            ts.theme_id, ts.config, ts.version, ts.updated_at,
            r.name
        FROM tenants t
        INNER JOIN users_tenants ut ON t.id = ut.tenant_id
        INNER JOIN roles r ON ut.role_id = r.id
        LEFT JOIN tenant_settings ts ON t.id = ts.tenant_id
        WHERE t.id = $1 AND ut.user_id = $2
        LIMIT 1;
    `

	t := &entity.Tenant{
		Settings: &entity.Settings{
			Config: make(entity.ThemeConfig), // always initialized
		},
	}

	var statusDB string
	var updatedAtDB, deletedAtDB *time.Time
	var settingsUpdatedAtDB *time.Time
	var themeIDDB value_object.ThemeID
	var versionDB int32
	var cfg entity.ThemeConfig
	var roleNameDB string

	err := r.DB.QueryRow(ctx, query, tenantID, userID).Scan(
		&t.ID,
		&t.StoreName,
		&t.BusinessType,
		&t.Slug,
		&t.SubDomain,
		&t.Domain,
		&statusDB,
		&t.CreatedAt,
		&updatedAtDB,
		&deletedAtDB,
		&themeIDDB,
		&cfg, // JSONB → map[string]any (pgx v5 auto-unmarshal)
		&versionDB,
		&settingsUpdatedAtDB,
		&roleNameDB,
	)
	if err != nil {
		return nil, fmt.Errorf("[%s]: %w", op, domain.TranslateTenantRepoError(postgres.MapPostgresError(err)))
	}

	t.Status = statusDB
	t.UpdatedAt = updatedAtDB
	t.DeletedAt = deletedAtDB

	// Map settings (always exists in your system, but still defensive)
	t.Settings.ThemeID = themeIDDB
	t.Settings.TenantID = t.ID
	t.Settings.Version = int(versionDB)

	if settingsUpdatedAtDB != nil {
		t.Settings.UpdatedAt = *settingsUpdatedAtDB
	}

	if cfg != nil {
		t.Settings.Config = cfg
	}

	// Return tenant context with role name
	return &TenantContext{
		Tenant: t,
		Role:   roleNameDB,
	}, nil
}

// UpdateTenantSettings performs an atomic, deep-merge update of a tenant's theme configuration.
//
// PURPOSE:
//   - Allows partial updates to tenant theme settings (colors, layout, typography, etc.)
//     without overwriting the entire config.
//
// HOW IT WORKS:
//  1. Begins a database transaction.
//  2. Fetches the current config from `tenant_settings` using `FOR UPDATE`
//     to lock the row and prevent concurrent modifications.
//  3. Performs a **deep merge** of the incoming config into the existing config.
//  4. Persists the merged result back to the database, incrementing the version
//     and updating the timestamp.
//  5. Commits the transaction.
//  6. Returns the **fully hydrated tenant context** (tenant + settings + role).
//
// CONCURRENCY:
// - Uses row-level locking (`FOR UPDATE`) to ensure safe concurrent updates.
// - Prevents lost updates when multiple users modify the same config.
//
// DATA GUARANTEES:
// - Config is never null (always at least an empty map).
// - Only provided fields are updated; all others are preserved.
// - Version is incremented on every successful update.
//
// RETURNS:
//   - Full TenantContext (not partial data), ensuring the caller receives
//     the latest consistent state after update.
//
// NOTE:
//   - This approach replaces PostgreSQL JSONB shallow merge (`||`) with
//     application-level deep merge for correct nested updates.
func (r *TenantRepository) UpdateTenantSettings(ctx context.Context, tenantID value_object.TenantID, userID value_object.UserID, incoming entity.ThemeConfig,
) (*TenantContext, error) {

	const op = "TenantRepository.UpdateTenantSettings"

	tx, err := r.DB.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// get current config
	var currentConfig entity.ThemeConfig

	err = tx.QueryRow(ctx, `
        SELECT config
        FROM tenant_settings
        WHERE tenant_id = $1
        FOR UPDATE
    `, tenantID).Scan(&currentConfig)

	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to fetch current config: %w", op, err)
	}

	// Ensure maps are initialized
	if currentConfig == nil {
		currentConfig = make(entity.ThemeConfig)
	}
	if incoming == nil {
		incoming = make(entity.ThemeConfig)
	}

	// Deep merge
	merged := deepMerge(currentConfig, incoming)

	_, err = tx.Exec(ctx, `
        UPDATE tenant_settings
        SET 
            config = $1,
            version = version + 1,
            updated_at = NOW()
        WHERE tenant_id = $2
    `, merged, tenantID)

	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to update config: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("[%s]: failed to commit tx: %w", op, err)
	}

	// Return FULL tenant context
	return r.GetTenantContext(ctx, tenantID, userID)
}

// deepMerge recursively merges two JSON-like maps (map[string]any).
//
// BEHAVIOR:
// - Performs a **deep (recursive) merge** of `src` into `dst`.
// - For each key in `src`:
//   - If the value is a nested map, it merges recursively.
//   - Otherwise, it overwrites the value in `dst`.
//
// - Keys not present in `src` remain unchanged in `dst`.
//
// EXAMPLE:
// dst:
//
//	{
//	  "colors": { "primary": "#000", "secondary": "#fff" },
//	  "layout": { "header": "sticky" }
//	}
//
// src:
//
//	{
//	  "colors": { "primary": "#f00" }
//	}
//
// result:
//
//	{
//	  "colors": { "primary": "#f00", "secondary": "#fff" },
//	  "layout": { "header": "sticky" }
//	}
//
// NOTES:
//   - This is used to support **partial theme updates** without overwriting
//     the entire configuration.
//   - Unlike PostgreSQL JSONB `||`, this performs a true deep merge.
//   - Mutates and returns `dst` (not immutable).
func deepMerge(dst, src map[string]any) map[string]any {
	for k, v := range src {
		if vMap, ok := v.(map[string]any); ok {
			if dstMap, ok := dst[k].(map[string]any); ok {
				dst[k] = deepMerge(dstMap, vMap)
			} else {
				dst[k] = deepMerge(make(map[string]any), vMap)
			}
		} else {
			dst[k] = v
		}
	}
	return dst
}
