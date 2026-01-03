package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MartinMurithi/storeforge/auth/internal/apperrors"
	"github.com/MartinMurithi/storeforge/auth/internal/dto"
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

func (srv *UserService) RegisterUser(ctx context.Context, input *dto.RegisterUserRequestDTO) (*models.User, error) {
	const op = "UserService.RegisterUser"

	// Normalize user input
	input.Normalize()

	checks := []struct {
		FieldName string
		Value     string
		Err       error
	}{
		{"FullName", input.FullName, apperrors.ErrFullNameRequired},
		{"Email", input.Email, apperrors.ErrEmailRequired},
		{"Phone", input.Phone, apperrors.ErrPhoneRequired},
		{"Password", input.Password, apperrors.ErrPasswordRequired},
		{"BusinessType", input.BusinessType, apperrors.ErrBusinessTypeRequired},
		{"BusinessName", input.BusinessName, apperrors.ErrBusinessNameRequired},
	}

	for _, check := range checks {
		if check.Value == "" {
			log.Printf("[%s] missing required field '%s':", op, check.FieldName)
			return nil, check.Err
		}
	}

	if err := lib.ValidateEmail(input.Email); err != nil {
		log.Printf("[%s] error validating email '%s': ", op, input.Email)
		return nil, err
	}

	_, err := lib.ValidatePhone(input.Phone)

	if err != nil {
		log.Printf("[%s] error validating phone number '%s': ", op, input.Phone)
		return nil, err
	}

	//check if user already exists
	existingUser, err := srv.repo.GetUserByEmail(ctx, input.Email)

	if existingUser != nil {
		log.Printf("[%s] user with email %s is already registered ", op, input.Email)
		return nil, apperrors.ErrUserAlreadyExists
	}

	//hashpassword
	hashedPassword, err := utils.Hashpassword(input.Password)

	if err != nil {
		log.Printf("error hashing password %s", err)
		return nil, fmt.Errorf("internal server error")
	}

	newUser := &models.User{
		FullName:     input.FullName,
		Email:        input.Email,
		Phone:        input.Phone,
		PasswordHash: hashedPassword,
		BusinessType: input.BusinessType,
		BusinessName: input.BusinessName,
	}

	//save user to db
	err = srv.repo.CreateUser(ctx, newUser)

	if err != nil {
		log.Printf("%s: error occurred when registering user %v", op, err)
		return nil, fmt.Errorf("internal server error")
	}

	return newUser, nil
}

func (srv *UserService) LoginUser(ctx context.Context, input *dto.LoginUserRequestDTO) (*models.User, string, error) {
	const op = "UserService.LoginUser"

	input.Normalize()

	if input.Email == "" || input.Password == "" {
		return nil, "", fmt.Errorf("%s:both email and password are required ", op)
	}

	if err := lib.ValidateEmail(input.Email); err != nil {
		log.Printf("[%s] error validating email '%s': ", op, input.Email)
		return nil, "", err
	}

	//check if user already exists
	existingUser, err := srv.repo.GetUserByEmail(ctx, input.Email)

	fmt.Println("exisiting user", existingUser)

	if err != nil || existingUser == nil {
		return nil, "", fmt.Errorf("invalid email or password %w", err)
	}

	//verify password
	err = utils.VerifyPassword(input.Password, existingUser.PasswordHash)

	if err != nil {
		return nil, "", fmt.Errorf("invalid email or password %w", err)
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
