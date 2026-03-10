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
