package rbac

import (
	"log"

	rbacv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/rbac/v1"
	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MapCreateRoleRequest converts a gRPC request to an Application DTO

func MapCreateRoleRequest(pbReq *rbacv1.CreateRoleRequest) (*dto.CreateRoleRequestDTO, error) {
	ids := make([]pgtype.UUID, 0, len(pbReq.PermissionIds))

	for _, idStr := range pbReq.PermissionIds {
		parsed, err := uuid.Parse(idStr)
		if err != nil {
			log.Printf("invalid permission uuid '%s': %w", idStr, err)
			return nil, apperrors.ErrInvalidPermissionID
		}

		ids = append(ids, pgtype.UUID{Bytes: parsed, Valid: true})
	}

	return &dto.CreateRoleRequestDTO{
		Name:          pbReq.Name,
		Slug:          pbReq.Slug,
		Description:   pbReq.Description,
		PermissionIDs: ids,
	}, nil
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
			CreatedAt:   timestamppb.New(p.CreatedAt),
		}
	}

	return &rbacv1.Role{
		Id:          role.ID.String(),
		Name:        role.Name,
		Slug:        role.Slug,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		Permissions: pbPermissions,
		CreatedAt:   timestamppb.New(role.CreatedAt),
	}
}
