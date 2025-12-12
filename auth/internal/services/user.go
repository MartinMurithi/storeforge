package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/MartinMurithi/storeforge/auth/internal/models"
	"github.com/MartinMurithi/storeforge/auth/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

// create a factory function to initialize my service with repo
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

type RegisterUserDTO struct {
	FullName     string
	Email        string
	Phone        string
	PasswordHash string
	BusinessType string
	BusinessName string
}

func (srv *UserService) RegisterUser(ctx context.Context, user *RegisterUserDTO) (*models.User, error) {
	const op = "UserService.RegisterUser"
	if user.FullName == "" || user.Email == "" || user.Phone == "" || user.PasswordHash == "" || user.BusinessType == "" || user.BusinessName == "" {
		return nil, fmt.Errorf("%s:all fields are required ", op)
	}

	//verify if phone is valid
	//verify if email is valid
	fullName := strings.TrimSpace(user.FullName)
	email := strings.ToLower(strings.TrimSpace(user.Email))
	phone := strings.TrimSpace(user.Phone)
	password := strings.TrimSpace(user.PasswordHash)
	businessName := strings.TrimSpace(user.BusinessName)
	businessType := strings.TrimSpace(user.BusinessType)

	//check if email is already registered
	//check if business name already exists

	existingUser, err := srv.repo.GetUserByEmail(ctx, email)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to get user by email %w", op, err)
	}

	if existingUser != nil {
		return nil, fmt.Errorf("%s: user with email %s already exists %w", op, email, err)
	}

	//hashpassword
	//salt password

	// create user first for testing

	newUser := &models.User{
		FullName:     fullName,
		Email:        email,
		Phone:        phone,
		PasswordHash: password,
		BusinessType: businessType,
		BusinessName: businessName,
	}

	//save user to db
	err = srv.repo.CreateUser(ctx, newUser)

	if err != nil {
		return nil, fmt.Errorf("%s: error occurred when creating user %w", op, err)
	}

	return newUser, nil
}

func (srv *UserService) FetchAllUsers(ctx context.Context, page, limit int) ([]*models.User, error) {
	const op = "UserService.FetchAllUsers"
	users, err := srv.repo.GetAllUsers(ctx, page, limit)

	if err != nil {
		return nil, fmt.Errorf("%s: error fetching users %w", op, err)
	}

	if len(users) == 0 {
		return []*models.User{}, nil
	}

	return users, nil
}
