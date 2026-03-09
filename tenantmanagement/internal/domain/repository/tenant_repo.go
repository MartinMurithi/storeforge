package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/database"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/database/postgres"
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
		return fmt.Errorf("[%s]: %w", op, err)
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
		log.Printf("[%s]: error creating tenant: %v", op, err)

		infraErr := postgres.MapPostgresError(err)
		return domain.TranslateUserRepoError(infraErr)
	}

	return tx.Commit(ctx)
}
