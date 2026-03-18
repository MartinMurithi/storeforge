package rbac

import (
	"context"
	"fmt"
	"log"

	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"
)

type RoleService struct {
	RoleRepo       repository.IRoleRepository
	PermissionRepo repository.IPermissionRepository
}

// create a factory function to initialize my service with repo
func NewRoleService(roleRepo repository.IRoleRepository, permissionRepo repository.IPermissionRepository) *RoleService {

	return &RoleService{
		RoleRepo:       roleRepo,
		PermissionRepo: permissionRepo,
	}
}

func (srv *RoleService) CreateRole(ctx context.Context, input *dto.CreateRoleRequestDTO) (*entity.Role, error) {
	const op = "RoleService.CreateRole"

	input.Normalize()
	if input.Name == "" || input.Slug == "" {
		return nil, fmt.Errorf("%s: name and slug are required", op)
	}

	// Check if slug is already taken early before starting a heavy transaction
	existingRole, err := srv.RoleRepo.GetRoleBySlug(ctx, input.Slug)
	if err == nil && existingRole != nil {
		return nil, apperrors.ErrRoleAlreadyExists
	}

	// Check if all Permission IDs provided actually exist
	var permissions []*entity.Permission
	if len(input.PermissionIDs) > 0 {

		permissions, err = srv.PermissionRepo.GetPermissionsById(ctx, input.PermissionIDs)
		if err != nil {
			log.Printf("[%s] failed to verify permissions: %v", op, err)
			return nil, err
		}

		if len(permissions) != len(input.PermissionIDs) {
			return nil, fmt.Errorf("%s: one or more permission IDs are invalid", op)
		}
	}

	newRole := &entity.Role{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Permissions: permissions,
	}

	if err := srv.RoleRepo.CreateRole(ctx, newRole); err != nil {
		log.Printf("[%s] create role: %v", op, err)
		return nil, err
	}

	log.Printf("[%s] successfully created role: %s (ID: %s)", op, newRole.Slug, newRole.ID)

	return newRole, nil
}
