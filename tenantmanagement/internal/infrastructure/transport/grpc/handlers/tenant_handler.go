package handlers

import (
	"context"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/application/dtos"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/application/services/tenant"
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
	// Map Proto Request -> Application DTO
	dtoReq := dtos.CreateTenantRequestDTO{
		StoreName:    req.StoreName,
		Slug:         req.Slug,
		BusinessType: req.BusinessType,
		SubDomain:    req.SubDomain,
		ThemeID:      req.ThemeId,
	}

	result, err := h.TenantService.CreateTenant(ctx, dtoReq)

	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	return result, nil
}
