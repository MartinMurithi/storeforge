package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/apperrors"
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
	GetUserById(ctx context.Context, id pgtype.UUID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*entity.User, error)
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

func (repo *UserRepository) GetUserById(ctx context.Context, id pgtype.UUID) (*entity.User, error) {
	const op = "UserRepository.GetUserById"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user := &entity.User{}

	query := `
		SELECT id, full_name, email, phone, business_type, business_name, created_at, updated_at, deleted_at, is_verified
		FROM users
		WHERE id = $1
	`

	err := repo.DB.QueryRow(ctx, query, id).Scan(
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
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}

		return nil, TranslateUserRepoError(postgres.MapPostgresError(err))
	}

	return user, nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	const op = "UserRepository.GetUserByEmail"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user := &entity.User{}

	query := `
		SELECT id, full_name, email, phone, password_hash, business_type, business_name, created_at, updated_at, is_verified
		FROM users
		WHERE email = $1
	`

	err := repo.DB.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.BusinessType,
		&user.BusinessName,
		&user.CreatedAt,
		&user.UpdatedAt,
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

func (repo *UserRepository) GetUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	const op = "UserRepository.GetUserByPhone"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user := &entity.User{}

	query := `
		SELECT id, full_name, email, phone, password_hash, business_type, business_name, created_at, updated_at, is_verified
		FROM users
		WHERE phone = $1
	`

	err := repo.DB.QueryRow(ctx, query, phone).Scan(
		&user.ID,
		&user.FullName,
		&user.Email,
		&user.Phone,
		&user.PasswordHash,
		&user.BusinessType,
		&user.BusinessName,
		&user.CreatedAt,
		&user.UpdatedAt,
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

func (repo *UserRepository) GetAllUsers(ctx context.Context, p dto.Pagination) ([]*entity.User, int, error) {
	const op = "UserRepository.GetAllUsers"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// --- Total Users Count ---
	var totalUsers int
	if err := repo.DB.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&totalUsers); err != nil {
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
		WHERE id = $3
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
		WHERE id = $1
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
