package rbac

import (
	"context"

	rbacv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/rbac/v1"
	"github.com/MartinMurithi/storeforge/pkg/errconv"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/rbac"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
)

type RoleGrpcHandler struct {
	RoleSrv *rbac.RoleService
	rbacv1.UnimplementedRbacServiceServer
}

func NewRoleGrpcHandler(r *rbac.RoleService) *RoleGrpcHandler {
	return &RoleGrpcHandler{
		RoleSrv: r,
	}
}

func (h *RoleGrpcHandler) CreateRole(ctx context.Context, req *rbacv1.CreateRoleRequest) (*rbacv1.CreateRoleResponse, error) {

	newIds, err := MapCreateRoleRequest(req)

	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	dtoReq := &dto.CreateRoleRequestDTO{
		Name:          req.Name,
		Slug:          req.Slug,
		Description:   req.Description,
		PermissionIDs: newIds.PermissionIDs,
	}

	role, err := h.RoleSrv.CreateRole(ctx, dtoReq)

	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	return &rbacv1.CreateRoleResponse{
		Success: true,
		Role:    MapRoleToPb(role),
	}, nil
}
