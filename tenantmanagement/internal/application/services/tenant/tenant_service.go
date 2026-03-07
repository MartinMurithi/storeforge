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

	themeID, _ := value_object.NewThemeID(req.ThemeID)
	theme, err := s.themeRepo.GetThemeById(ctx, themeID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: %w", op, err)
	}

	newTenant := &entity.Tenant{
		StoreName:    req.StoreName,
		Slug:         req.Slug,
		BusinessType: req.BusinessType,
		SubDomain:    req.SubDomain,
	}

	tenantConfig := theme.DefaultConfig.Config
	if tenantConfig == nil {
		tenantConfig = make(entity.ThemeConfig)
	}

	newTenant.Settings = &entity.Settings{
		ThemeID:   theme.ID,
		Config:    tenantConfig,
		Version:   1,
		UpdatedAt: time.Now(),
	}

	if err := s.tenantRepo.CreateTenant(ctx, newTenant); err != nil {
		return nil, fmt.Errorf("[%s]: %w", op, err)
	}

	return mappers.ToProtoCreateTenantResponse(&dtos.CreateTenantResponseDTO{
		Tenant:         newTenant,
		Theme:          theme,
		TenantSettings: newTenant.Settings,
	}), nil
}
