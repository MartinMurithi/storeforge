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

	"github.com/jackc/pgx/v5/pgtype"
)

type TenantRepository struct {
	DB database.DB
}

type ITenantRepository interface {
	CreateTenant(ctx context.Context, tenant *entity.Tenant) error
	GetTenantContext(ctx context.Context, tenantID value_object.TenantID, userID value_object.UserID) (*TenantContext, error)
	UpdateTenantSettings(ctx context.Context, tenantID value_object.TenantID, config entity.ThemeConfig) (*entity.Tenant, error)
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
	RoleId pgtype.UUID
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
            ut.role_id
        FROM tenants t
        INNER JOIN users_tenants ut ON t.id = ut.tenant_id
        LEFT JOIN tenant_settings ts ON t.id = ts.tenant_id
        WHERE t.id = $1 AND ut.user_id = $2;
    `

	t := &entity.Tenant{
		Settings: &entity.Settings{
			Config: make(entity.ThemeConfig), // default to empty map
		},
	}

	var statusDB string
	var updatedAtDB, deletedAtDB, settingsUpdatedAtDB *time.Time
	var themeIDDB value_object.ThemeID
	var versionDB int32
	var cfg entity.ThemeConfig
	var roleIDDB pgtype.UUID

	// pgx v5 automatically unmarshals JSONB into Go structs/maps
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
		&cfg, // automatic JSONB → ThemeConfig
		&versionDB,
		&settingsUpdatedAtDB,
		&roleIDDB,
	)
	if err != nil {
		return nil, fmt.Errorf("[%s]: %w", op, domain.TranslateTenantRepoError(postgres.MapPostgresError(err)))
	}

	// Map tenant fields
	t.Status = string(entity.TenantProvisioning)
	t.UpdatedAt = updatedAtDB
	t.DeletedAt = deletedAtDB

	// Map settings
	t.Settings.ThemeID = themeIDDB
	t.Settings.TenantID = t.ID
	t.Settings.Version = int(versionDB)
	t.Settings.UpdatedAt = time.Now()
	if cfg != nil {
		t.Settings.Config = cfg
	}

	return &TenantContext{
		Tenant: t,
		RoleId: roleIDDB,
	}, nil
}

// UpdateTenant updates a tenant's fields and/or settings partially.
// Only updates fields provided in the DTO; others remain unchanged.
// UpdateTenantSettings partially updates a tenant's theme configuration and increments the version.
// It uses a transaction to ensure atomic updates to the configuration and audit timestamps.
//
// SECURITY: This method assumes tenantID validation has occurred at the Service/Handler layer.
func (r *TenantRepository) UpdateTenantSettings(ctx context.Context, tenantID value_object.TenantID, config entity.ThemeConfig) (*entity.Tenant, error) {
	const op = "TenantRepository.UpdateTenantSettings"

	tx, err := r.DB.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// Using the JSONB merge operator (||) to allow partial config updates
	// and COALESCE to handle nil inputs gracefully.
	query := `
        UPDATE tenant_settings
        SET
            config     = config || COALESCE($1, '{}'::jsonb),
            version    = version + 1,
            updated_at = NOW()
        WHERE tenant_id = $2
        RETURNING theme_id, config, version, updated_at
    `

	tn := &entity.Tenant{}

	err = tx.QueryRow(ctx, query, config, tenantID).Scan(
		&tn.Settings.ThemeID,
		&tn.Settings.Config,
		&tn.Settings.Version,
		&tn.Settings.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to update settings: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("[%s]: failed to commit tx: %w", op, err)
	}

	return tn, nil
}
