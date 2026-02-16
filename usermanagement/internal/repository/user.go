package repository

import (
	"context"
	"errors"
	"log"
	"time"

	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database/postgres"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	DB database.DB
}

type IUserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetActiveUserById(ctx context.Context, id pgtype.UUID) (*entity.User, error)
	GetUserByIdIcludingDeleted(ctx context.Context, id pgtype.UUID) (*entity.User, error)
	GetActiveUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByEmailIncludingDeleted(ctx context.Context, email string) (*entity.User, error)
	GetActiveUserByPhone(ctx context.Context, phone string) (*entity.User, error)
	GetUserByPhoneIncludingDeleted(ctx context.Context, phone string) (*entity.User, error)
	GetAllUsers(ctx context.Context, p dto.Pagination) ([]*entity.User, int, error)
	PatchUser(ctx context.Context, id pgtype.UUID, input *UpdateUserInput) (*entity.User, error)
	DeleteUser(ctx context.Context, id pgtype.UUID) error
}

func NewUserRepository(db database.DB) IUserRepository {
	return &UserRepository{DB: db}
}

type UpdateUserInput struct {
	BusinessName *string
	BusinessType *string
}

type userLookupMode int

const (
	onlyActive userLookupMode = iota
	includeDeleted
)

func (repo *UserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	const op = "UserRepository.CreateUser"

	// Set a timeout for the DB operation
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (
			full_name, email, phone, password_hash, business_type, business_name
		) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	// Execute query
	err := repo.DB.QueryRow(
		ctx,
		query,
		user.FullName,
		user.Email,
		user.Phone,
		user.PasswordHash,
		user.BusinessType,
		user.BusinessName,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		log.Printf("[%s]: error creating user: %v", op, err)

		// Map DB/driver errors → infra errors
		infraErr := postgres.MapPostgresError(err)

		// Translate infra → domain errors
		return TranslateUserRepoError(infraErr)
	}

	return nil
}

func (repo *UserRepository) getUser(
	ctx context.Context,
	where string,
	arg any,
	mode userLookupMode,
) (*entity.User, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		SELECT id, full_name, email, phone, password_hash,
		       business_type, business_name,
		       created_at, updated_at, deleted_at, is_verified
		FROM users
		WHERE ` + where

	if mode == onlyActive {
		query += " AND deleted_at IS NULL"
	}

	user := &entity.User{}

	err := repo.DB.QueryRow(ctx, query, arg).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.BusinessType,
		&user.BusinessName,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.IsVerified,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, TranslateUserRepoError(postgres.MapPostgresError(err))
	}

	return user, nil
}

func (repo *UserRepository) GetActiveUserById(ctx context.Context, id pgtype.UUID) (*entity.User, error) {
	const op = "UserRepository.GetActiveById"
	return repo.getUser(ctx, "id=$1", id, onlyActive)
}

func (repo *UserRepository) GetUserByIdIcludingDeleted(ctx context.Context, id pgtype.UUID) (*entity.User, error) {
	const op = "UserRepository.GetByIdIcludingDeleted"
	return repo.getUser(ctx, "id=$1", id, includeDeleted)
}

func (repo *UserRepository) GetActiveUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	const op = "UserRepository.GetActiveByEmail"
	return repo.getUser(ctx, "email=$1", email, onlyActive)
}

func (repo *UserRepository) GetUserByEmailIncludingDeleted(ctx context.Context, email string) (*entity.User, error) {
	const op = "UserRepository.GetByEmailIncludingDeleted"
	return repo.getUser(ctx, "email=$1", email, includeDeleted)
}

func (repo *UserRepository) GetActiveUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	const op = "UserRepository.GetActiveByPhone"
	return repo.getUser(ctx, "phone=$1", phone, onlyActive)
}

func (repo *UserRepository) GetUserByPhoneIncludingDeleted(ctx context.Context, phone string) (*entity.User, error) {
	const op = "UserRepository.GetByPhoneIncludingDeleted"
	return repo.getUser(ctx, "phone=$1", phone, includeDeleted)
}

func (repo *UserRepository) GetAllUsers(ctx context.Context, p dto.Pagination) ([]*entity.User, int, error) {
	const op = "UserRepository.GetAllUsers"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// --- Total Users Count ---
	var totalUsers int
	if err := repo.DB.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&totalUsers); err != nil {
		return nil, 0, TranslateUserRepoError(postgres.MapPostgresError(err))
	}

	const maxLimit = 15
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > maxLimit {
		p.Limit = maxLimit
	}

	offset := (p.Page - 1) * p.Limit

	query := `
		SELECT id, full_name, email, phone, business_type, business_name, created_at, updated_at, deleted_at, is_verified
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := repo.DB.Query(ctx, query, p.Limit, offset)
	if err != nil {
		return nil, 0, TranslateUserRepoError(postgres.MapPostgresError(err))
	}
	defer rows.Close()

	users := make([]*entity.User, 0, p.Limit)
	for rows.Next() {
		user := &entity.User{}
		if err := rows.Scan(
			&user.ID,
			&user.FullName,
			&user.Email,
			&user.Phone,
			&user.BusinessType,
			&user.BusinessName,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
			&user.IsVerified,
		); err != nil {
			return nil, 0, TranslateUserRepoError(postgres.MapPostgresError(err))
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, TranslateUserRepoError(postgres.MapPostgresError(err))
	}

	return users, totalUsers, nil
}

// Remember to check if business name already exists
func (repo *UserRepository) PatchUser(ctx context.Context, id pgtype.UUID, input *UpdateUserInput) (*entity.User, error) {
	const op = "UserRepository.PatchUser"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET
			business_name = COALESCE($1, business_name),
			business_type = COALESCE($2, business_type),
			updated_at = NOW()
		WHERE deleted_at IS NULL AND id = $3
		RETURNING
			id,
			full_name,
			email,
			phone,
			business_type,
			business_name,
			is_verified,
			created_at,
			updated_at
	`

	var user entity.User

	err := repo.DB.QueryRow(
		ctx,
		query,
		input.BusinessName,
		input.BusinessType,
		id,
	).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.BusinessType,
		&user.BusinessName,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}

		return nil, TranslateUserRepoError(postgres.MapPostgresError(err))
	}

	return &user, nil
}

// Soft Delete, these updates the deleted_at column
func (repo *UserRepository) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	const op = "UserRepository.DeleteUser"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id
	`

	var returnedID pgtype.UUID
	err := repo.DB.QueryRow(ctx, query, id).Scan(&returnedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrUserNotFound
		}

		return TranslateUserRepoError(postgres.MapPostgresError(err))
	}

	return nil
}
