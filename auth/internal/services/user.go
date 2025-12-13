package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	// "time"

	"github.com/MartinMurithi/storeforge/auth/internal/lib"
	"github.com/MartinMurithi/storeforge/auth/internal/models"
	"github.com/MartinMurithi/storeforge/auth/internal/repository"
	"github.com/MartinMurithi/storeforge/auth/internal/token"
	"github.com/MartinMurithi/storeforge/auth/internal/utils"
)

type UserService struct {
	repo     *repository.UserRepository
	jwtMaker *token.JWTMaker
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

	fullName := strings.TrimSpace(user.FullName)
	email := strings.ToLower(strings.TrimSpace(user.Email))
	phone := strings.TrimSpace(user.Phone)
	password := strings.TrimSpace(user.PasswordHash)
	businessName := strings.TrimSpace(user.BusinessName)
	businessType := strings.TrimSpace(user.BusinessType)

	if !lib.IsEmailValid(email) {
		return nil, fmt.Errorf("%s: invalid email format", op)
	}

	ValidatedPhone, err := lib.ValidatePhone(phone)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//check if user already exists
	existingUser, err := srv.repo.GetUserByEmail(ctx, email)

	if err != nil {
		log.Printf("db error: %s", err)
		return nil, fmt.Errorf("%s: internal error", op)
	}

	if existingUser != nil {
		return nil, fmt.Errorf("%s: user already exists", op)
	}

	//hashpassword
	hashedPassword, err := utils.Hashpassword(password)

	if err != nil {
		log.Printf("error hashing password %s", err)
		return nil, fmt.Errorf("%s: internal error", op)
	}

	newUser := &models.User{
		FullName:     fullName,
		Email:        email,
		Phone:        ValidatedPhone,
		PasswordHash: hashedPassword,
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

func (srv *UserService) LoginUser(email, password string, ctx context.Context) (string, error) {
	const op = "UserService.LoginUser"
	if email == "" || password == "" {
		return "", fmt.Errorf("%s:both email and password are required ", op)
	}

	//verify if email is valid
	sanitizedEmail := strings.ToLower(strings.TrimSpace(email))
	sanitizedPassword := strings.TrimSpace(password)

	if !lib.IsEmailValid(email) {
		return "", fmt.Errorf("%s: invalid email or password", op)
	}

	//check if user already exists
	existingUser, err := srv.repo.GetUserByEmail(ctx, sanitizedEmail)

	if err != nil {
		return "", fmt.Errorf("%s: invalid email or password", op)
	}

	//check password
	err = utils.CheckPassword(sanitizedPassword, existingUser.PasswordHash)

	if err != nil {
		return "", fmt.Errorf("invalid email or password %s", err)
	}

	// Before issuing JWT, create a tenant first, will revisit this later

	// Generate JWT
	// token, _, err := srv.jwtMaker.CreateToken(existingUser.ID, existingUser.Email, existingUser.Role.Name, time.Hour)

	// if err != nil {
	// 	log.Printf("%s: error creating token %s", op, err)
	// 	return "", fmt.Errorf("failed to issue token %w", err)
	// }

	return "", nil
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
