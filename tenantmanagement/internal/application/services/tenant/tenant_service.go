package tenant

import (
	"context"
	"fmt"
	"log"
	"time"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	authv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/auth/v1"
	membershipv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/membership/v1"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/application/dtos"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/application/mappers"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/repository"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/clients"
)

// TenantService implements the business logic for tenant lifecycle management.
type TenantService struct {
	tenantRepo repository.ITenantRepository
	themeRepo  repository.IThemeRepository
	userSvc    *clients.UserServiceClient
}

// NewTenantService creates a new instance of the TenantService.
func NewTenantService(tr repository.ITenantRepository, th repository.IThemeRepository, us *clients.UserServiceClient) *TenantService {
	return &TenantService{
		tenantRepo: tr,
		themeRepo:  th,
		userSvc:    us,
	}
}

// CreateTenant orchestrates the registration of a new store.
// It fetches the theme template, clones the settings, and persists everything atomically.
func (s *TenantService) CreateTenant(ctx context.Context, req dtos.CreateTenantRequestDTO) (*tenantv1.CreateTenantResponse, error) {
	const op = "TenantService.CreateTenant"

	log.Printf("theme id : %v", req)

	themeID, _ := value_object.NewThemeID(req.ThemeID)

	theme, err := s.themeRepo.GetThemeById(ctx, themeID)

	if err != nil {
		return nil, fmt.Errorf("[%s]: %w", op, err)
	}

	if theme == nil {
		return nil, fmt.Errorf("[%s]: theme not found", op)
	}

	log.Printf("feteched theme: %s", theme.Name)
	// Generate slug from store_name
	slug := GenerateSlug(req.StoreName)

	// Generate sub_domain from slug + domain
	subDomain, err := GenerateSubdomain(slug)

	if err != nil {
		return nil, fmt.Errorf("an error occurred when generating a sub domain: %w", err)
	}

	domain := FullDomain(subDomain)

	newTenant := &entity.Tenant{
		StoreName:    req.StoreName,
		Slug:         slug,
		BusinessType: req.BusinessType,
		SubDomain:    subDomain,
		Domain:       domain,
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

	// Link tenant to user
	// Use timeout to ensure tenant svc does not hang if linking fails
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	linkUserReq := &membershipv1.LinkUserToTenantRequest{
		UserId:   req.UserId,
		TenantId: newTenant.ID.String(),
		Role:     "owner",
	}

	resp, err := s.userSvc.LinkUserToTenant(ctx, linkUserReq)

	log.Printf("link user to tenant result %v", resp)

	if err != nil {
		log.Printf("[%s]: failed to link user %s to tenant %s: %v", op, linkUserReq.UserId, newTenant.ID, err)

		// For now, return an error so the user knows something went wrong
		return nil, fmt.Errorf("store created but ownership link failed: %w", err)
	}

	// Return new JWT Token after user has been linked to their store, incase linking fails, new JWT should not be issued
	updateActiveSessionReq := &authv1.UpdateSessionContextRequest{
		UserId:   req.UserId,
		TenantId: newTenant.ID.String(),
		Role:     "owner",
	}

	newToken, err := s.userSvc.UpdateActiveSessionContext(ctx, updateActiveSessionReq)

	tokenDTO := &dtos.TokenInfoDTO{
		AccessToken: newToken.Token.AccessToken,
		TokenType:   newToken.Token.TokenType,
		ExpiresIn:   newToken.Token.ExpiresIn,
		ExpiresAt:   time.Now().Add(time.Duration(newToken.Token.ExpiresIn) * time.Second),
		IssuedAt:    time.Now(),
	}

	return mappers.ToProtoCreateTenantResponse(&dtos.CreateTenantResponseDTO{
		Tenant:         newTenant,
		Theme:          theme,
		TenantSettings: newTenant.Settings,
		Token:          tokenDTO,
	}), nil
}
