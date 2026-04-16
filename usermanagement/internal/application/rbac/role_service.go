package rbac

import (
	"context"
	"fmt"
	"log"

	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
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

func (srv *RoleService) GetRoleById(ctx context.Context, id pgtype.UUID) (*entity.Role, error) {
	const op = "role_service.GetRoleByID"

	log.Printf("role id %v", id)

	role, err := srv.RoleRepo.GetRoleByID(ctx, id)

	log.Printf("role role %v", role)

	if err != nil {
		log.Printf("[%s] error %v", op, err)
		return nil, err
	}

	return role, nil
}

func (srv *RoleService) UpdateRole(ctx context.Context, roleID pgtype.UUID, input *dto.UpdateRoleRequestDTO) (*entity.Role, error) {
	const op = "RoleService.UpdateRole"

	// input.Normalize()

	log.Printf("role ID: %s", roleID)

	role, err := srv.RoleRepo.GetRoleByID(ctx, roleID)
	if err != nil {
		log.Printf("[%s] failed to fetch role: %v", op, err)
		return nil, err
	}
	if role == nil {
		return nil, apperrors.ErrRoleNotFound
	}

	// Prevent updates to system roles eg owner
	if role.IsSystem {
		return nil, fmt.Errorf("%s: system roles cannot be updated", op)
	}

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

	updatedRole := &entity.Role{
		ID:          role.ID,
		Name:        input.Name,
		Description: input.Description,
		Slug:        role.Slug, // slug stays immutable
		IsSystem:    role.IsSystem,
		CreatedAt:   role.CreatedAt,
		Permissions: permissions, // list of permission IDs for repo to replace
	}

	if err := srv.RoleRepo.UpdateRole(ctx, updatedRole); err != nil {
		log.Printf("[%s] update role: %v", op, err)
		return nil, err
	}

	updatedRole.Permissions = permissions

	log.Printf("[%s] successfully updated role: %s (ID: %s)", op, updatedRole.Name, updatedRole.ID)

	newRole, err := srv.RoleRepo.GetRoleByID(ctx, updatedRole.ID)

	if err != nil {
			log.Printf("[%s] failed to fetch error: %v", op, err)
			return nil, err
		}

	return newRole, nil
}
