package rbac

import (
	"context"
	"log"

	rbacv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/rbac/v1"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

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

func (h *RoleGrpcHandler) GetRoleByID(ctx context.Context, req *rbacv1.GetRoleByIDRequest) (*rbacv1.GetRoleByIDResponse, error) {

	// Convert RoleID from string to pgtype.uuid
	parsedID, err := uuid.Parse(req.RoleId)
	if err != nil {
		log.Printf("invalid role uuid': %w", err)
		return nil, apperrors.ErrInvalidUUIDFormat
	}

	roleID := pgtype.UUID{Bytes: parsedID, Valid: true}

	role, err := h.RoleSrv.GetRoleById(ctx, roleID)

	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	return &rbacv1.GetRoleByIDResponse{
		Role: MapRoleToPb(role),
	}, nil

}

func (h *RoleGrpcHandler) UpdateRole(ctx context.Context, req *rbacv1.UpdateRoleRequest) (*rbacv1.UpdateRoleResponse, error) {

	newIds, err := MapUpdateRoleRequest(req)

	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	dtoReq := &dto.UpdateRoleRequestDTO{
		RoleID:        newIds.RoleID,
		Name:          req.Name,
		Description:   req.Description,
		PermissionIDs: newIds.PermissionIDs,
	}

	role, err := h.RoleSrv.UpdateRole(ctx, dtoReq.RoleID, dtoReq)

	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	return &rbacv1.UpdateRoleResponse{
		Success: true,
		Role:    MapRoleToPb(role),
	}, nil
}
