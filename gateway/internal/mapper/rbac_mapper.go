package mapper

import (
	rbacv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/rbac/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
)

// MapCreateRoleDTOToProto converts the Gin request DTO to a gRPC message
func MapCreateRoleDTOToProto(req dto.CreateRoleRequestDTO) *rbacv1.CreateRoleRequest {
    return &rbacv1.CreateRoleRequest{
        Name:          req.Name,
        Slug:          req.Slug,
        Description:   req.Description,
        PermissionIds: req.PermissionIDs,
    }
}

// MapRoleProtoToDTO converts a single Role gRPC message to a JSON DTO
func MapRoleProtoToDTO(pbRole *rbacv1.Role) dto.RoleResponseDTO {
    if pbRole == nil {
        return dto.RoleResponseDTO{}
    }

    perms := make([]dto.PermissionResponseDTO, len(pbRole.Permissions))
    for i, p := range pbRole.Permissions {
        perms[i] = dto.PermissionResponseDTO{
            ID:          p.Id,
            Slug:        p.Slug,
            Description: p.Description,
        }
    }

    return dto.RoleResponseDTO{
        ID:          pbRole.Id,
        Name:        pbRole.Name,
        Slug:        pbRole.Slug,
        Description: pbRole.Description,
        IsSystem:    pbRole.IsSystem,
        Permissions: perms,
        CreatedAt:   pbRole.CreatedAt.AsTime().Format("2006-01-02 15:04:05"),
    }
}

// MapCreateRoleResponseProtoToDTO converts the gRPC response to the final Gateway API response
func MapCreateRoleResponseProtoToDTO(pbRes *rbacv1.CreateRoleResponse) dto.CreateRoleResponseDTO {
    return dto.CreateRoleResponseDTO{
        Role:    MapRoleProtoToDTO(pbRes.Role),
        Message: "Role created successfully",
    }
}