package tenant

import (
	"context"
	"fmt"
	"log"
	"time"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	authv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/auth/v1"
	membershipv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/membership/v1"
	"github.com/MartinMurithi/storeforge/pkg/rbac"
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
		return nil, err
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

	log.Printf("user id in tenant service %s", req.UserId)
	log.Printf("tenant id  in tenant service %s", newTenant.ID.String())

	resp, err := s.userSvc.LinkUserToTenant(ctx, linkUserReq)

	log.Printf("link user to tenant result %v", resp)

	if err != nil {
		log.Printf("[%s]: store created but ownership link failed: %v", op, err)

		// For now, return an error so the user knows something went wrong
		// To Do: Delete the created tenant(store) to avoid orphaned stores in the database
		return nil, err
	}

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

// GetTenantContext retrieves the full tenant context for a given user and tenant.
// Includes tenant info, settings, and the caller’s role.
func (s *TenantService) GetTenantContext(ctx context.Context, req dtos.GetTenantContextRequestDTO) (*tenantv1.GetTenantContextResponse, error) {
	const op = "TenantService.GetTenantContext"

	tenantId, err := value_object.NewTenantID(req.TenantId)
	if err != nil {
		return nil, fmt.Errorf("[%s]: invalid tenant id: %w", op, err)
	}

	userID, err := value_object.NewUserID(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("[%s]: invalid user id: %w", op, err)
	}

	// Fetch tenant context from repository (includes tenant + settings + role)
	tenantCtx, err := s.tenantRepo.GetTenantContext(ctx, tenantId, userID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to fetch tenant context: %w", op, err)
	}

	if tenantCtx == nil || tenantCtx.Tenant == nil {
		return nil, fmt.Errorf("[%s]: tenant context not found", op)
	}

	tenant := tenantCtx.Tenant

	// Defensive: settings always exist but double-check
	settings := tenant.Settings
	if settings == nil {
		settings = &entity.Settings{
			Config: make(entity.ThemeConfig),
		}
	}

	// role info if needed for caller context
	role := tenantCtx.Role

	return &tenantv1.GetTenantContextResponse{
		Tenant:   mappers.ToProtoTenant(tenant),
		Settings: mappers.ToProtoSettings(settings),
		Role:     role,
	}, nil
}

func (s *TenantService) UpdateTenant(ctx context.Context, req *dtos.UpdateTenantRequestDTO) (*tenantv1.GetTenantContextResponse, error) {

	const op = "TenantService.UpdateTenant"

	tenantID, err := value_object.NewTenantID(req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: invalid tenant id: %w", op, err)
	}

	log.Printf("[%s] service tenant id: %v", op, tenantID)

	userID, err := value_object.NewUserID(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: invalid user id: %w", op, err)
	}

	tenantCtx, err := s.tenantRepo.GetTenantContext(ctx, tenantID, userID)
	if err != nil {
		log.Printf("[%s] failed to fetch tenant context: %v",op, err)
		return nil, err
	}

	if tenantCtx.Tenant == nil {
		return nil, fmt.Errorf("[%s]: tenant not found", op)
	}

	// RBAC
	if req.Settings == nil {
		return nil, fmt.Errorf("[%s]: settings payload is required", op)
	}

	if tenantCtx.Role != rbac.RoleOwner && tenantCtx.Role != rbac.RoleAdmin {
		return nil, fmt.Errorf("unauthorized, only admin and owner are allowed to customize theme")
	}

	config := req.Settings.Config
	if config == nil {
		config = make(entity.ThemeConfig)
	}

	updatedCtx, err := s.tenantRepo.UpdateTenantSettings(ctx, tenantID, userID, config)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to update settings: %w", op, err)
	}

	return &tenantv1.GetTenantContextResponse{
		Tenant:   mappers.ToProtoTenant(updatedCtx.Tenant),
		Settings: mappers.ToProtoSettings(updatedCtx.Tenant.Settings),
		Role:     updatedCtx.Role,
	}, nil
}
