package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MartinMurithi/storeforge/auth/internal/database/config"
	"github.com/MartinMurithi/storeforge/auth/internal/lib/db"
	"github.com/MartinMurithi/storeforge/auth/internal/models"
	"github.com/jackc/pgx/v5"

	"github.com/google/uuid"
)

type UserRepository struct {
	DB *config.Pool
}

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetAllUsers(ctx context.Context, page, limit int) ([]*models.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

func NewUserRepository(pool *config.Pool) *UserRepository {
	return &UserRepository{DB: pool}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) error {

	const op = "UserRepository.CreateUser"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (full_name, email, phone, password_hash, business_type, business_name) VALUES($1, $2, $3, $4, $5, $6) returning id`

	tx, err := repo.DB.Begin(ctx)

	if err != nil {
		return fmt.Errorf("%s: error starting a transaction %w", op, err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			fmt.Printf("%s: rollback failed %s", op, err)
		}
	}()

	err = tx.QueryRow(ctx, query, user.FullName, user.Email, user.Phone, user.PasswordHash, user.BusinessType, user.BusinessName).Scan(&user.ID)

	if err != nil {
		return db.WrapDbError(ctx, op, 5*time.Second, err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		return fmt.Errorf("%s: error occurred when committing a transaction %w", op, err)
	}

	return nil
}

func (repo *UserRepository) GetAllUsers(ctx context.Context, page, limit int) ([]*models.User, error) {
	const op = "UserRepository.GetAllUsers"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	const maxLimit = 15

	if page < 1 {
		page = 1
	}

	if limit == 0 || limit > maxLimit {
		limit = maxLimit
	}

	offset := (page - 1) * limit

	//add limit and offset for pagination
	query := `SELECT id, full_name, email, phone, business_type, business_name, created_at, is_verified, role FROM users
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	rows, err := repo.DB.Query(ctx, query, limit, offset)

	if err != nil {
		return nil, db.WrapDbError(ctx, op, 5*time.Second, err)
	}

	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		user := &models.User{}

		err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.BusinessType, &user.BusinessName, &user.CreatedAt, &user.IsVerified, &user.Role)

		if err != nil {
			return nil, db.WrapDbError(ctx, op, 5*time.Second, err)
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, db.WrapDbError(ctx, op, 5*time.Second, err)
	}

	return users, nil
}

func (repo *UserRepository) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	const op = "UserRepository.GetUserById"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user := &models.User{}

	query := `SELECT id, full_name, email, phone, business_type, business_name, created_at, updated_at, is_verified, roles FROM users
	WHERE id = $1`

	err := repo.DB.QueryRow(ctx, query, id).Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.BusinessType, &user.BusinessName, &user.CreatedAt, &user.UpdatedAt, &user.IsVerified, &user.Role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: user not found %w", op, err)
		}
		return nil, db.WrapDbError(ctx, op, 5*time.Second, err)
	}
	return user, nil
}

func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	const op = "UserRepository.GetUserByEmail"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user := &models.User{}

	query := `SELECT id, full_name, email, phone, business_type, business_name, created_at, updated_at, is_verified, roles FROM users
	WHERE email = $1`

	err := repo.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.BusinessType, &user.BusinessName, &user.CreatedAt, &user.UpdatedAt, &user.IsVerified, &user.Role)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: user not found %w", op, err)
		}
		return nil, db.WrapDbError(ctx, op, 5*time.Second, err)
	}
	return user, nil
}

func (repo *UserRepository) UpdateUser(ctx context.Context, id uuid.UUID, user *models.User) error {
	const op = "UserRepository.UpdateUser"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `UPDATE users 
	SET business_name=$1, business_type=$2, updated_at=$3
	WHERE id=$4`

	tx, err := repo.DB.Begin(ctx)

	if err != nil {
		return fmt.Errorf("%s: error starting a transaction %w", op, err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			fmt.Printf("%s: rollback failed %s", op, err)
		}
	}()

	result, err := tx.Exec(ctx, query, user.BusinessName, user.BusinessType, user.UpdatedAt, id.String())

	if err != nil {
		return db.WrapDbError(ctx, op, 3*time.Second, err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("%s: user not found id%w", op, id)
	}

	user.ID = id

	err = tx.Commit(ctx)

	if err != nil {
		return fmt.Errorf("%s: error occurred when committing a transaction %w", op, err)
	}

	return nil
}

func (repo *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	const op = "UserRepository.DeleteUser"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `DELETE FROM users WHERE id=$1`

	result, err := repo.DB.Exec(ctx, query, id)

	if err != nil {
		return db.WrapDbError(ctx, op, 3*time.Second, err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return fmt.Errorf("%s: user not found id%w", op, id)
	}
	return nil
}
