package tenant

import (
	"context"
	"fmt"
	"time"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/application/dtos"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/application/mappers"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"
)

// TenantService implements the business logic for tenant lifecycle management.
type TenantService struct {
	tenantRepo domain.ITenantRepository
	themeRepo  domain.IThemeRepository
}

// NewTenantService creates a new instance of the TenantService.
func NewTenantService(tr domain.ITenantRepository, th domain.IThemeRepository) *TenantService {
	return &TenantService{
		tenantRepo: tr,
		themeRepo:  th,
	}
}

// CreateTenant orchestrates the registration of a new store.
// It fetches the theme template, clones the settings, and persists everything atomically.
func (s *TenantService) CreateTenant(ctx context.Context, req dtos.CreateTenantRequestDTO) (*tenantv1.CreateTenantResponse, error) {
	const op = "TenantService.CreateTenant"

	// Fetch the Theme blueprint from the master catalog
	// We need the 'Golden Template' settings stored in this Theme.
	themeID, err := value_object.NewThemeID(req.ThemeID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: invalid theme id: %w", op, err)
	}

	theme, err := s.themeRepo.GetThemeById(ctx, themeID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to fetch theme template: %w", op, err)
	}

	newTenant := &entity.Tenant{
		StoreName:    req.StoreName,
		Slug:         req.Slug,
		BusinessType: req.BusinessType,
		SubDomain:    req.SubDomain,
	}

	// Clone the Theme's DefaultSettings for this specific Tenant
	// This creates a unique 'Settings' instance for the tenant so they can
	// customize their config without affecting the master theme.
	newTenant.Settings = &entity.Settings{
		ThemeID:   theme.ID,
		Config:    theme.DefaultConfig.Config, // The map[string]any snapshot
		Version:   1,                          // Initial version
		UpdatedAt: time.Now(),
	}

	// Persist the Tenant and Settings
	if err := s.tenantRepo.CreateTenant(ctx, newTenant); err != nil {
		return nil, fmt.Errorf("[%s]: failed to persist store: %w", op, err)
	}

	resp := mappers.ToProtoCreateTenantResponse(&dtos.CreateTenantResponseDTO{
		Tenant:         newTenant,
		Theme:          theme,
		TenantSettings: newTenant.Settings,
	})

	return resp, nil
}
