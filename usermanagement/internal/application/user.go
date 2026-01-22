package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/apperrors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/utils"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/token"


	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
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

type PatchUserInput struct {
	Id           pgtype.UUID
	BusinessName *string
	BusinessType *string
}

func (srv *UserService) RegisterUser(ctx context.Context, input *dto.RegisterUserRequestDTO) (*entity.User, error) {
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

	if err := utils.ValidateEmail(input.Email); err != nil {
		log.Printf("[%s] error validating email '%s': ", op, input.Email)
		return nil, err
	}

	_, err := utils.ValidatePhone(input.Phone)

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

func (srv *UserService) LoginUser(ctx context.Context, input *dto.LoginUserRequestDTO) (*entity.User, *token.Token, error) {
	const op = "UserService.LoginUser"

	input.Normalize()

	if input.Email == "" || input.Password == "" {
		return nil, nil, fmt.Errorf("%s:both email and password are required ", op)
	}

	if err := utils.ValidateEmail(input.Email); err != nil {
		log.Printf("[%s] error validating email '%s': ", op, input.Email)
		return nil, nil, err
	}

	//check if user already exists
	existingUser, err := srv.repo.GetUserByEmail(ctx, input.Email)

	fmt.Println("exisiting user", existingUser)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			log.Printf("[%s] user not found '%s': ", op, err)
			return nil, nil, apperrors.ErrUserNotFound
		}
		log.Printf("[%s] get user by email failed '%s': ", op, err)
		return nil, nil, err
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

func (srv *UserService) FetchAllUsers(ctx context.Context, p dto.Pagination) ([]*entity.User, dto.PaginationMeta, error) {
	const op = "UserService.FetchAllUsers"

	users, total, err := srv.repo.GetAllUsers(ctx, p)

	if err != nil {
		return nil, dto.PaginationMeta{}, fmt.Errorf("%s: error fetching users %w", op, err)
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + p.Limit - 1) / p.Limit
	}

	meta := dto.PaginationMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    p.Page < totalPages,
		HasPrev:    p.Page > 1,
	}

	return users, meta, nil
}

func (srv *UserService) GetCurrentUserById(ctx context.Context, id pgtype.UUID) (*entity.User, error) {
	const op = "UserService.FetchUserById"

	log.Printf("user id %v", id.Valid)

	user, err := srv.repo.GetUserById(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("%s: error fetching user %w", op, err)
	}

	return user, nil
}

func (srv *UserService) UpdateCurrentUser(ctx context.Context, input *PatchUserInput) (*entity.User, error) {
	const op = "UserService.UpdateCurrentUser"

	log.Printf("user id %v", input.Id.Valid)

	patch := &repository.UpdateUserInput{
		BusinessName: input.BusinessName,
		BusinessType: input.BusinessType,
	}

	updatedUser, err := srv.repo.PatchUser(ctx, input.Id, patch)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, fmt.Errorf("[%s]: %w", op, apperrors.ErrUserNotFound)
		}
		return nil, fmt.Errorf("[%s]: [%w]", op, err)
	}

	return updatedUser, nil
}

func (srv *UserService) SoftDeleteUser(ctx context.Context, id pgtype.UUID) error {
	const op = "UserService.SoftDeleteUser"

	log.Printf("user id %v", id.Valid)

	err := srv.repo.DeleteUser(ctx, id)

	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return fmt.Errorf("[%s]: %w", op, apperrors.ErrUserNotFound)
		}
		return fmt.Errorf("[%s]: [%w]", op, err)
	}

	return nil
}
