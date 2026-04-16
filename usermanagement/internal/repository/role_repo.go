package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database/postgres"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type RoleRepository struct {
	DB database.DB
}

type IRoleRepository interface {
	CreateRole(ctx context.Context, role *entity.Role) error
	GetRoleBySlug(ctx context.Context, slug string) (*entity.Role, error)
	GetRoleByID(ctx context.Context, roleID pgtype.UUID) (*entity.Role, error)
	UpdateRole(ctx context.Context, input *entity.Role) error
}

func NewRoleRepository(db database.DB) IRoleRepository {
	return &RoleRepository{DB: db}
}

func (repo *RoleRepository) CreateRole(ctx context.Context, role *entity.Role) error {
	const op = "role_repository.CreateRole"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := repo.DB.Tx(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("%s: rollback failed: %v", op, err)
		}
	}()

	const query = `
		INSERT INTO roles (name, slug, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	if err := tx.QueryRow(ctx, query, role.Name, role.Slug, role.Description).Scan(&role.ID, &role.CreatedAt); err != nil {
		log.Printf("%s: create role: %v", op, err)
		return fmt.Errorf("%w", TranslateRoleRepoError(postgres.MapPostgresError(err)))
	}

	const linkQuery = `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)
	`

	seen := make(map[pgtype.UUID]struct{})

	for _, perm := range role.Permissions {
		if _, ok := seen[perm.Id]; ok {
			continue
		}
		seen[perm.Id] = struct{}{}

		if _, err := tx.Exec(ctx, linkQuery, role.ID, perm.Id); err != nil {
			log.Printf("%s: link permission %s to role %s: %v", op, perm.Id, role.ID, err)
			return fmt.Errorf("%w", TranslateRoleRepoError(postgres.MapPostgresError(err)))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit tx: %w", op, err)
	}

	return nil
}

func (repo *RoleRepository) GetRoleBySlug(ctx context.Context, slug string) (*entity.Role, error) {
	const op = "role_repo.GetRoleBySlug"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT id, name, slug, description, is_system
		FROM roles
		WHERE slug = $1
		`

	role := &entity.Role{}

	err := repo.DB.QueryRow(ctx, query, slug).Scan(
		&role.ID,
		&role.Name,
		&role.Slug,
		&role.Description,
		&role.IsSystem,
		&role.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrRoleNotFound
		}
		return nil, TranslateRoleRepoError(postgres.MapPostgresError(err))
	}

	return role, nil
}

func (repo *RoleRepository) GetRoleByID(ctx context.Context, roleID pgtype.UUID) (*entity.Role, error) {
	const op = "role_repo.GetRoleByID"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT
			r.id,
			r.name,
			r.slug,
			r.description,
			r.is_system,
			r.created_at,
			p.id,
			p.slug,
			p.description
		FROM roles r
		LEFT JOIN role_permissions rp ON r.id = rp.role_id
		LEFT JOIN permissions p ON rp.permission_id = p.id
		WHERE r.id = $1
	`

	rows, err := repo.DB.Query(ctx, query, roleID)
	if err != nil {
		return nil, TranslateRoleRepoError(postgres.MapPostgresError(err))
	}
	defer rows.Close()

	role := &entity.Role{
		Permissions: []*entity.Permission{},
	}

	found := false

	for rows.Next() {
		found = true

		var (
			permID          pgtype.UUID
			permSlug        *string
			permDescription *string
		)

		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Slug,
			&role.Description,
			&role.IsSystem,
			&role.CreatedAt,
			&permID,
			&permSlug,
			&permDescription,
		)
		if err != nil {
			return nil, TranslateRoleRepoError(postgres.MapPostgresError(err))
		}

		// Only append if permission exists (LEFT JOIN → can be NULL)
		if permID.Valid {
			role.Permissions = append(role.Permissions, &entity.Permission{
				Id:          permID,
				Slug:        *permSlug,
				Description: *permDescription,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, TranslateRoleRepoError(postgres.MapPostgresError(err))
	}

	if !found {
		return nil, apperrors.ErrRoleNotFound
	}

	return role, nil
}

// UpdateRole performs an atomic update of a role's core fields and its permissions.
// It updates the role's name and description in the 'roles' table, then replaces
// all entries in 'role_permissions' with the provided final permission set.
func (r *RoleRepository) UpdateRole(ctx context.Context, input *entity.Role) error {
	const op = "RoleRepository.UpdateRole"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := r.DB.Tx(ctx)
	if err != nil {
		return fmt.Errorf("[%s]: %w", op, err)
	}
	defer tx.Rollback(ctx)

	roleQuery := `
		UPDATE roles
        SET name = COALESCE(NULLIF($1, ''), name),
            description = COALESCE(NULLIF($2, ''), description)
        WHERE id = $3 AND is_system = false
        RETURNING id, name, slug, description, is_system, created_at
 	`

	var role entity.Role

	err = tx.QueryRow(ctx, roleQuery, input.Name, input.Description, input.ID).
		Scan(
			&role.ID,
			&role.Name,
			&role.Slug,
			&role.Description,
			&role.IsSystem,
			&role.CreatedAt,
		)
	if err != nil {
		log.Printf("[%s error]: %v", op, err)
		return fmt.Errorf("%w", TranslateRoleRepoError(postgres.MapPostgresError(err)))
	}

	// DIFF LOGIC FOR PERMISSIONS
	// If Permissions slice is nil, skip this part to leave them as they are.
	// If it's empty [], it means "remove all".
	if input.Permissions != nil {
		clearQuery := `DELETE FROM role_permissions WHERE role_id = $1`
		if _, err = tx.Exec(ctx, clearQuery, input.ID); err != nil {
			return fmt.Errorf("clear perms: %w", err)
		}

		if len(input.Permissions) > 0 {
			insertQuery := `INSERT INTO role_permissions (permission_id, role_id) VALUES ($1, $2)`
			for _, p := range input.Permissions {
				if _, err = tx.Exec(ctx, insertQuery, p.Id, input.ID); err != nil {
					return fmt.Errorf("insert perm %s: %w", p.Id, err)
				}
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: commit tx: %w", op, err)
	}
	return nil
}
