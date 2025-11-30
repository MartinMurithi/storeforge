package repository

import (
	"context"
	"time"

	"github.com/MartinMurithi/storeforge/internal/database/config"
	"github.com/MartinMurithi/storeforge/pkg/dbhelper"
	"github.com/MartinMurithi/storeforge/internal/models"
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
		return dbhelper.WrapDbError(err)
	}

	return nil
}
