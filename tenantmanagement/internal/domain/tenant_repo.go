package domain

import (
	"context"
	"fmt"
	"log"
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
func (r *TenantRepository) CreateTenant(ctx context.Context, t *entity.Tenant) error {
    const op = "TenantRepository.CreateTenant"

    // 5-second deadline to prevent hanging database connections
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // Start the transaction; if this fails, no rows are touched
    tx, err := r.DB.Tx(ctx)
    if err != nil {
        log.Printf("[%s]: transaction start error: %v", op, err)
        return fmt.Errorf("[%s]: could not start transaction: %w", op, err)
    }

    // Ensure cleanup: if the function returns before tx.Commit(), 
    // the transaction is automatically aborted (rolled back).
    defer tx.Rollback(ctx)

    // Insert the core Tenant record. 
    tenantQuery := `
        INSERT INTO tenants (store_name, business_type, slug, sub_domain)
        VALUES($1, $2, $3, $4) 
        RETURNING id, store_name`

    err = tx.QueryRow(ctx, tenantQuery, 
        t.StoreName, 
        t.BusinessType, 
        t.Slug, 
        t.SubDomain,
    ).Scan(&t.ID, &t.StoreName)

    if err != nil {
        log.Printf("[%s]: tenant insert error: %v", op, err)
        return fmt.Errorf("[%s]: failed to insert tenant: %w", op, err)
    }

    // Insert the Tenant Settings.
    // The 'config' field (map[string]any) is serialized as JSONB by the driver.
    // This links the new Tenant ID to their selected Theme ID.
    settingsQuery := `
        INSERT INTO tenant_settings (tenant_id, theme_id, config, version)
        VALUES($1, $2, $3, $4)
        RETURNING updated_at`

    err = tx.QueryRow(ctx, settingsQuery, 
        t.ID, 
        t.Settings.ThemeID, 
        t.Settings.Config, 
        t.Settings.Version,
    ).Scan(&t.Settings.UpdatedAt)

    if err != nil {
        log.Printf("[%s]: settings insert error: %v", op, err)
        return fmt.Errorf("[%s]: failed to attach settings: %w", op, err)
    }

    // Finalize the transaction. 
    // Both the Tenant and their Settings are now permanently stored.
    if err := tx.Commit(ctx); err != nil {
        log.Printf("[%s]: commit error: %v", op, err)
        return fmt.Errorf("[%s]: failed to commit transaction: %w", op, err)
    }

    return nil
}
