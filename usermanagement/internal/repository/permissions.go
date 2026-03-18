package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/database"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database/postgres"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"

	"github.com/jackc/pgx/v5/pgtype"
)

type PermissionRepository struct {
	DB database.DB
}

type IPermissionRepository interface {
	GetPermissions(ctx context.Context) ([]*entity.Permission, error)
	GetPermissionsById(ctx context.Context, permId []pgtype.UUID) ([]*entity.Permission, error)
}

func NewPermissionRepository(db database.DB) IPermissionRepository {
	return &PermissionRepository{DB: db}
}

func (p *PermissionRepository) GetPermissionsById(ctx context.Context, permIds []pgtype.UUID) ([]*entity.Permission, error) {
	const op = "permission_repository.GetPermissions"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, slug, category, description
		FROM permissions
		WHERE id = ANY($1)
	`

	rows, err := p.DB.Query(ctx, query, permIds)
	if err != nil {
		return nil, TranslateUserRepoError(postgres.MapPostgresError(err))
	}
	defer rows.Close()

	perms := make([]*entity.Permission, 0, len(permIds))
	
	for rows.Next() {
		perm := &entity.Permission{}
		if err := rows.Scan(
			&perm.Id,
			&perm.Slug,
			&perm.Category,
			&perm.Description,
		); err != nil {
			return nil, fmt.Errorf("scan %w", TranslateUserRepoError(postgres.MapPostgresError(err)))
		}
		perms = append(perms, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan %w", TranslateUserRepoError(postgres.MapPostgresError(err)))
	}
	return perms, nil
}

func (p *PermissionRepository) GetPermissions(ctx context.Context) ([]*entity.Permission, error) {
	const op = "permission_repository.GetPermissions"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT id, slug, category, description
		FROM permissions
		ORDER BY category DESC
	`

	rows, err := p.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("scan %w", TranslateUserRepoError(postgres.MapPostgresError(err)))
	}
	defer rows.Close()

	perms := make([]*entity.Permission, 0, 10)
	for rows.Next() {
		perm := &entity.Permission{}
		if err := rows.Scan(
			&perm.Id,
			&perm.Slug,
			&perm.Category,
			&perm.Description,
		); err != nil {
			return nil, fmt.Errorf("scan %w", TranslateUserRepoError(postgres.MapPostgresError(err)))
		}
		perms = append(perms, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan %w", TranslateUserRepoError(postgres.MapPostgresError(err)))
	}
	return perms, nil
}
