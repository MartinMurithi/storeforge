package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/apperrors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/token"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo     repository.IUserRepository
	jwtMaker *token.JWTMaker
}

// create a factory function to initialize my service with repo
func NewAuthService(repo repository.IUserRepository, jwtMaker *token.JWTMaker) *AuthService {

	if jwtMaker == nil {
		panic("jwt maker must not be nil")
	}

	return &AuthService{
		repo:     repo,
		jwtMaker: jwtMaker,
	}
}

func (srv *AuthService) RegisterUser(ctx context.Context, input *dto.RegisterUserRequestDTO) (*entity.User, error) {
	const op = "AuthService.RegisterUser"

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

	if err := utils.ValidateEmail(input.Email); err != nil {
		log.Printf("[%s] error validating email '%s': ", op, input.Email)
		return nil, err
	}

	_, err := utils.ValidatePhone(input.Phone)

	if err != nil {
		log.Printf("[%s] error validating phone number '%s': ", op, input.Phone)
		return nil, err
	}

	//check if user already exists by email
	existingUser, err := srv.repo.GetActiveUserByEmail(ctx, input.Email)

	if existingUser != nil {
		log.Printf("[%s] user with email %s is already registered ", op, input.Email)
		return nil, apperrors.ErrUserEmailAlreadyExists
	}

	//check if phone already exists
	existingPhone, err := srv.repo.GetActiveUserByPhone(ctx, input.Phone)

	if existingPhone != nil {
		log.Printf("[%s] user with phone number %s already exists ", op, input.Phone)
		return nil, apperrors.ErrUserMobileExists
	}

	//hashpassword
	hashedPassword, err := utils.Hashpassword(input.Password)

	if err != nil {
		log.Printf("error hashing password %s", err)
		return nil, fmt.Errorf("internal server error")
	}

	newUser := &entity.User{
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

func (srv *AuthService) LoginUser(ctx context.Context, input *dto.LoginUserRequestDTO) (*entity.User, *entity.Token, error) {
	const op = "AuthService.LoginUser"

	input.Normalize()

	if input.Email == "" || input.Password == "" {
		return nil, nil, fmt.Errorf("%s:both email and password are required ", op)
	}

	if err := utils.ValidateEmail(input.Email); err != nil {
		log.Printf("[%s] error validating email '%s': ", op, input.Email)
		return nil, nil, err
	}

	//check if user already exists
	existingUser, err := srv.repo.GetUserByEmailIncludingDeleted(ctx, input.Email)

	fmt.Println("exisiting user", existingUser)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			log.Printf("[%s] user not found '%s': ", op, err)
			return nil, nil, apperrors.ErrUserNotFound
		}
		log.Printf("[%s] get user by email failed '%s': ", op, err)
		return nil, nil, err
	}

	if existingUser.DeletedAt != nil {
		log.Printf("your accoun has been deactivated. please contact support")
		return nil, nil, apperrors.ErrAccountDeactivated
	}

	//verify password
	err = utils.VerifyPassword(input.Password, existingUser.PasswordHash)

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, nil, apperrors.ErrInvalidCredentials
		}
		return nil, nil, err // unexpected crypto failure
	}

	// Before issuing JWT, create a tenant first(this will issue role to the user as owner), will revisit this later

	// Generate JWT
	token, _, err := srv.jwtMaker.CreateToken(existingUser.ID, existingUser.ID, existingUser.Email, "owner", 30*time.Minute)

	if err != nil {
		log.Printf("%s: error creating token %s", op, err)
		return nil, nil, fmt.Errorf("failed to issue token %w", err)
	}

	return existingUser, token, nil
}
