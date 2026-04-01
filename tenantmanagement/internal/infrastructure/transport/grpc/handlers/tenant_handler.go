package handlers

import (
	"context"
	"log"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	"github.com/MartinMurithi/storeforge/pkg/auth"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/application/dtos"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/application/services/tenant"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TenantGrpcHandler struct {
	tenantv1.UnimplementedTenantServiceServer
	TenantService *tenant.TenantService
}

// NewTenantGrpcHandler initializes the handler with the required application service.
func NewTenantGrpcHandler(s *tenant.TenantService) *TenantGrpcHandler {
	if s == nil {
		panic("NewTenantGrpcHandler: service is nil")
	}
	return &TenantGrpcHandler{
		TenantService: s,
	}
}

// CreateTenant converts the protobuf request into an internal DTO,
// executes the creation logic, and returns a mapped protobuf response.
func (h *TenantGrpcHandler) CreateTenant(ctx context.Context, req *tenantv1.CreateTenantRequest) (*tenantv1.CreateTenantResponse, error) {

	userID, err := auth.GetUserIDFromMetadata(ctx)

	if err != nil {
		log.Printf("failed to extract user id from metadata %s", err)
		return nil, status.Errorf(codes.Unauthenticated, "missing user identity: %v", err)
	}

	log.Printf("extracted user id %s", userID)

	dtoReq := dtos.CreateTenantRequestDTO{
		StoreName:    req.StoreName,
		BusinessType: req.BusinessType,
		ThemeID:      req.ThemeId,
		UserId:       userID,
	}

	result, err := h.TenantService.CreateTenant(ctx, dtoReq)

	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	return result, nil
}

func (h *TenantGrpcHandler) GetTenantContext(ctx context.Context, req *tenantv1.GetTenantContextRequest) (*tenantv1.GetTenantContextResponse, error) {
	// Extract authenticated user from metadata
	userID, err := auth.GetUserIDFromMetadata(ctx)
	if err != nil {
		log.Printf("[%s] failed to extract user id from metadata: %v", "GetTenantContext", err)
		return nil, status.Errorf(codes.Unauthenticated, "missing user identity: %v", err)
	}

	log.Printf("extracted user id %s", userID)

	tenantD, err := auth.GetTenantIDFromMetadata(ctx)
	if err != nil {
		log.Printf("[%s] failed to extract tenant id from metadata: %v", "GetTenantContext", err)
		return nil, status.Errorf(codes.Unauthenticated, "missing tenant identity: %v", err)
	}

	log.Printf("extracted tenant id %s", tenantD)

	dtoReq := dtos.GetTenantContextRequestDTO{
		UserId:   userID,
		TenantId: tenantD,
	}

	tenantCtxResp, err := h.TenantService.GetTenantContext(ctx, dtoReq)
	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	return tenantCtxResp, nil
}

func (h *TenantGrpcHandler) UpdateTenant(ctx context.Context, req *tenantv1.UpdateTenantRequest) (*tenantv1.GetTenantContextResponse, error) {
	const op = "TenantGrpcHandler.UpdateTenant"

	userID, err := auth.GetUserIDFromMetadata(ctx)
	if err != nil {
		log.Printf("[%s] failed to extract user id: %v", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing user identity: %v", err)
	}

	tenantID, err := auth.GetTenantIDFromMetadata(ctx)
	if err != nil {
		log.Printf("[%s] failed to extract tenant id: %v", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing tenant identity: %v", err)
	}

	dtoReq := &dtos.UpdateTenantRequestDTO{
		TenantID: tenantID,
		UserID:   userID,
	}

	if req.Settings != nil {
		dtoReq.Settings = &dtos.SettingsUpdateDTO{
			Config: req.Settings.ThemeConfig.AsMap(),
		}
	}

	log.Printf("[%s] tenant id : %v", op, dtoReq.TenantID)

	resp, err := h.TenantService.UpdateTenant(ctx, dtoReq)
	if err != nil {
		log.Printf("[%s] service error: %v", op, err)
		return nil, errconv.ToGrpcError(err)
	}

	return resp, nil
}
