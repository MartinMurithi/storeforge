package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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
func NewUserService(repo *repository.UserRepository, jwtMaker *token.JWTMaker) *UserService {

	if jwtMaker == nil {
		panic("jwt maker must not be nil")
	}

	return &UserService{
		repo:     repo,
		jwtMaker: jwtMaker,
	}
}

type RegisterUserInput struct {
	FullName     string
	Email        string
	Phone        string
	Password     string
	BusinessType string
	BusinessName string
}

type LoginUserInput struct {
	Email    string
	Password string
}

func (srv *UserService) RegisterUser(ctx context.Context, user *RegisterUserInput) (*models.User, error) {
	const op = "UserService.RegisterUser"

	if user.FullName == "" || user.Email == "" || user.Phone == "" || user.Password == "" || user.BusinessType == "" || user.BusinessName == "" {
		return nil, fmt.Errorf("%s:all fields are required ", op)
	}

	fullName := strings.TrimSpace(user.FullName)
	email := strings.ToLower(strings.TrimSpace(user.Email))
	phone := strings.TrimSpace(user.Phone)
	password := strings.TrimSpace(user.Password)
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

	if existingUser != nil {
		return nil, fmt.Errorf("%s: user with that email already exists", op)
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

func (srv *UserService) LoginUser(input *LoginUserInput, ctx context.Context) (*models.User, string, error) {
	const op = "UserService.LoginUser"

	if input.Email == "" || input.Password == "" {
		return nil, "", fmt.Errorf("%s:both email and password are required ", op)
	}

	sanitizedEmail := strings.ToLower(strings.TrimSpace(input.Email))
	sanitizedPassword := strings.TrimSpace(input.Password)

	if !lib.IsEmailValid(sanitizedEmail) {
		return nil, "", fmt.Errorf("%s: invalid email or password", op)
	}

	//check if user already exists
	existingUser, err := srv.repo.GetUserByEmail(ctx, sanitizedEmail)

	fmt.Println("exisiting user", existingUser)

	if err != nil || existingUser == nil {
		return nil, "", fmt.Errorf("%s: invalid email or password", op)
	}

	//verify password
	err = utils.VerifyPassword(sanitizedPassword, existingUser.PasswordHash)

	if err != nil {
		return nil, "", fmt.Errorf("invalid email or password %s", err)
	}

	// Before issuing JWT, create a tenant first(this will issue role to the user as owner), will revisit this later

	// Generate JWT
	token, _, err := srv.jwtMaker.CreateToken(existingUser.ID, existingUser.ID, existingUser.Email, "owner", 1*time.Hour)

	if err != nil {
		log.Printf("%s: error creating token %s", op, err)
		return nil, "", fmt.Errorf("failed to issue token %w", err)
	}

	return existingUser, token, nil
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
