package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	// "log"
	"time"

	"github.com/MartinMurithi/storeforge/auth/internal/database/config"
	"github.com/MartinMurithi/storeforge/auth/internal/dto"
	"github.com/MartinMurithi/storeforge/auth/internal/lib/db"
	"github.com/MartinMurithi/storeforge/auth/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository struct {
	DB *config.Pool
}

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id pgtype.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetAllUsers(ctx context.Context, page, limit int) ([]*models.User, error)
	UpdateUser(ctx context.Context, id pgtype.UUID, user *models.User) error
	DeleteUser(ctx context.Context, id pgtype.UUID) error
}

func NewUserRepository(pool *config.Pool) *UserRepository {
	return &UserRepository{DB: pool}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) error {

	const op = "UserRepository.CreateUser"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (full_name, email, phone, password_hash, business_type, business_name) VALUES($1, $2, $3, $4, $5, $6) returning id, created_at`

	tx, err := repo.DB.Begin(ctx)

	if err != nil {
		return fmt.Errorf("%s: error starting a transaction %w", op, err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			fmt.Printf("%s: rollback failed %s", op, err)
		}
	}()

	err = tx.QueryRow(ctx, query, user.FullName, user.Email, user.Phone, user.PasswordHash, user.BusinessType, user.BusinessName).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return db.WrapDbError(ctx, op, 5*time.Second, err)
	}

	err = tx.Commit(ctx)

	if err != nil {
		return fmt.Errorf("%s: error occurred when committing a transaction %w", op, err)
	}

	return nil
}

func (repo *UserRepository) GetAllUsers(
    ctx context.Context,
    p dto.Pagination,
) ([]*models.User, int, error) {

    const op = "UserRepository.GetAllUsers"

    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // --- Total Users Count ---
    var totalUsers int
    if err := repo.DB.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&totalUsers); err != nil {
        return nil, 0, db.WrapDbError(ctx, op, 5*time.Second, err)
    }

    const maxLimit = 15
    if p.Page < 1 { p.Page = 1 }
    if p.Limit <= 0 || p.Limit > maxLimit { p.Limit = maxLimit }

    offset := (p.Page - 1) * p.Limit

    query := `
    SELECT id, full_name, email, phone, business_type, business_name, created_at, updated_at, deleted_at, is_verified
    FROM users
    ORDER BY created_at DESC
    LIMIT $1 OFFSET $2`

    rows, err := repo.DB.Query(ctx, query, p.Limit, offset)
    if err != nil {
        return nil, 0, db.WrapDbError(ctx, op, 5*time.Second, err)
    }
    defer rows.Close()

    var users []*models.User
    for rows.Next() {
        user := &models.User{}

        err := rows.Scan(
            &user.ID, &user.FullName, &user.Email, &user.Phone,
            &user.BusinessType, &user.BusinessName, &user.CreatedAt,
            &user.UpdatedAt, &user.DeletedAt, &user.IsVerified,
        )
        
        if err != nil {
            fmt.Printf("[%s] SCAN ERROR: %v\n", op, err)
            return nil, 0, db.WrapDbError(ctx, op, 5*time.Second, err)
        }

        users = append(users, user)
    }

    if err := rows.Err(); err != nil {
        fmt.Printf("[%s] ROWS ERROR: %v\n", op, err)
        return nil, 0, db.WrapDbError(ctx, op, 5*time.Second, err)
    }

    fmt.Printf("[%s] Returning %d users\n", op, len(users))
    return users, totalUsers, nil
}

func (repo *UserRepository) GetUserById(ctx context.Context, id pgtype.UUID) (*models.User, error) {
	const op = "UserRepository.GetUserById"

	log.Printf("[REPO]: user id %v", id.Valid)

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	user := &models.User{}

	query := `SELECT id, full_name, email, phone, business_type, business_name, created_at, updated_at, deleted_at, is_verified FROM users
	WHERE id = $1`

	err := repo.DB.QueryRow(ctx, query, id).Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.BusinessType, &user.BusinessName, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.IsVerified)

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

	query := `SELECT id, full_name, email, phone, password_hash, business_type, business_name, created_at, updated_at, is_verified FROM users
	WHERE email = $1`

	err := repo.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.PasswordHash, &user.BusinessType, &user.BusinessName, &user.CreatedAt, &user.UpdatedAt, &user.IsVerified)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("[%s]: user not found: %w", op, err)
		}
		return nil, db.WrapDbError(ctx, op, 5*time.Second, err)
	}
	return user, nil
}

// Remember to check if business name already exists

func (repo *UserRepository) UpdateUser(ctx context.Context, id pgtype.UUID, user *models.User) error {
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
		return fmt.Errorf("%s: user not found id%s", op, id)
	}

	user.ID = id

	err = tx.Commit(ctx)

	if err != nil {
		return fmt.Errorf("%s: error occurred when committing a transaction %w", op, err)
	}

	return nil
}

func (repo *UserRepository) DeleteUser(ctx context.Context, id pgtype.UUID) error {
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
		return fmt.Errorf("%s: user not found id%s", op, id)
	}
	return nil
}
