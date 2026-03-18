package rbac

import (
	rbacv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/rbac/v1"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// MapCreateRoleRequest converts a gRPC request to an Application DTO
func MapCreateRoleRequest(pbReq *rbacv1.CreateRoleRequest) *dto.CreateRoleRequestDTO {
	ids := make([]pgtype.UUID, 0, len(pbReq.PermissionIds))

	for _, id := range pbReq.PermissionIds {
		parsedID := pgtype.UUID{Bytes: uuid.MustParse(id), Valid: true}

		ids = append(ids, parsedID)
	}

	return &dto.CreateRoleRequestDTO{
		Name:          pbReq.Name,
		Slug:          pbReq.Slug,
		Description:   pbReq.Description,
		PermissionIDs: ids,
	}
}

func MapRoleToPb(role *entity.Role) *rbacv1.Role {
	if role == nil {
		return nil
	}

	// Convert permissions slice
	pbPermissions := make([]*rbacv1.Permission, len(role.Permissions))
	for i, p := range role.Permissions {
		pbPermissions[i] = &rbacv1.Permission{
			Id:          p.Id.String(),
			Slug:        p.Slug,
			Description: p.Description,
		}
	}

	return &rbacv1.Role{
		Id:          role.ID.String(),
		Name:        role.Name,
		Slug:        role.Slug,
		IsSystem:    role.IsSystem,
		Permissions: pbPermissions,
	}
}
