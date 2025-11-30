package repository

import (
	"context"
	"time"

	"github.com/MartinMurithi/storeforge/auth/internal/database/config"
	"github.com/MartinMurithi/storeforge/auth/internal/lib/db"
	"github.com/MartinMurithi/storeforge/auth/internal/models"
)

type UserRepository struct {
	DB *config.Pool
}

type IUserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, id string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, id string, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id string) error
}

func NewUserRepository(pool *config.Pool) *UserRepository {
	return &UserRepository{DB: pool}
}

func (repo *UserRepository) CreateUser(ctx context.Context, user *models.User) error {

	const op = "UserRepository.CreateUser"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (full_name, email, phone, password_hash, business_type, business_name, is_verified) VALUES($1, $2, $3, $4, $5, $6, $7) returning id`

	err := repo.DB.QueryRow(ctx, query, user.FullName, user.Email, user.Phone, user.PasswordHash, user.BusinessType, user.BusinessName, user.IsVerified).Scan(&user.ID)

	if err != nil {
		return db.WrapDbError(ctx, op, 5*time.Second, err)
	}

	return nil
}

func (repo *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	const op = "UserRepository.GetAllUsers"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	//add limit and offset for pagination
	query := `SELECT id, full_name, email, phone, business_type, business_name, created_at, is_verified, roles FROM users
	ORDER BY created_at DESC
	`

	rows, err := repo.DB.Query(ctx, query)

	if err != nil {
		return nil, db.WrapDbError(ctx, op, 5*time.Second, err)
	}

	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		user := &models.User{}

		err := rows.Scan(&user.ID, &user.FullName, &user.Email, &user.Phone, &user.BusinessType, &user.BusinessName, &user.CreatedAt, &user.IsVerified, &user.Roles)

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
