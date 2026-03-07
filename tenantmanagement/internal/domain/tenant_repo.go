package domain

import (
	"context"
	"fmt"
	// "log"
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/database"
)

type TenantRepository struct {
	DB database.DB
}

type ITenantRepository interface {
	CreateTenant(ctx context.Context, tenant *entity.Tenant) error
	// GetTenantById(ctx context.Context, id string) (*entity.Tenant, error)
	// GetTenantBySlug(ctx context.Context, slug string) (*entity.Tenant, error)
}

func NewTenantRepository(db database.DB) ITenantRepository {
	return &TenantRepository{DB: db}
}

// CreateTenant performs an atomic dual-insertion into the 'tenants' and 'tenant_settings' tables.
// It uses a database transaction to ensure that a tenant is never created without its
// corresponding theme configuration.
//
// The process involves:
//  1. Initializing a transaction via the custom database adapter.
//  2. Persisting the core tenant identity (StoreName, Slug, etc.).
//  3. Linking the tenant to their chosen ThemeID and storing the initial ThemeConfig (JSONB).
//  4. Committing both records or rolling back entirely if either step fails.
// 
func (r *TenantRepository) CreateTenant(ctx context.Context, t *entity.Tenant) error {
    const op = "TenantRepository.CreateTenant"
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    tx, err := r.DB.Tx(ctx)
    if err != nil {
        return fmt.Errorf("[%s]: %w", op, err)
    }
    defer tx.Rollback(ctx)

    // 1. Insert Tenant (Returning 3 fields)
    tenantQuery := `INSERT INTO tenants (store_name, business_type, slug, sub_domain)
                    VALUES($1, $2, $3, $4) RETURNING id, store_name, status, created_at`

    // Fixed: Added &t.Status to match the 3 returned columns
    err = tx.QueryRow(ctx, tenantQuery, t.StoreName, t.BusinessType, t.Slug, t.SubDomain).
        Scan(&t.ID, &t.StoreName, &t.Status, &t.CreatedAt)
    if err != nil {
        return fmt.Errorf("[%s]: %w", op, err)
    }

    // 2. Link IDs in memory for the Mapper
    t.Settings.TenantID = t.ID

    // 3. Insert Settings
    settingsQuery := `INSERT INTO tenant_settings (tenant_id, theme_id, config, version)
                      VALUES($1, $2, $3, $4) RETURNING updated_at`

    configToPersist := t.Settings.Config
    if configToPersist == nil {
        configToPersist = make(entity.ThemeConfig)
    }

    err = tx.QueryRow(ctx, settingsQuery, t.ID, t.Settings.ThemeID, configToPersist, t.Settings.Version).
        Scan(&t.Settings.UpdatedAt)
    if err != nil {
        return fmt.Errorf("[%s]: %w", op, err)
    }

    return tx.Commit(ctx)
}
