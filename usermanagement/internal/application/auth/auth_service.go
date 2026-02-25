package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MartinMurithi/storeforge/pkg/auth"
	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/token"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.IUserRepository
	authRepo repository.IAuthRepository
	jwtMaker *token.JWTMaker
}

const refreshTokenExpiry = 30 * 24 * time.Hour // 30 days

// create a factory function to initialize my service with repo
func NewAuthService(userRepo repository.IUserRepository, authRepo repository.IAuthRepository, jwtMaker *token.JWTMaker) *AuthService {

	if jwtMaker == nil {
		panic("jwt maker must not be nil")
	}

	return &AuthService{
		userRepo: userRepo,
		authRepo: authRepo,
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
	existingUser, err := srv.userRepo.GetActiveUserByEmail(ctx, input.Email)

	if existingUser != nil {
		log.Printf("[%s] user with email %s is already registered ", op, input.Email)
		return nil, apperrors.ErrUserEmailAlreadyExists
	}

	//check if phone already exists
	existingPhone, err := srv.userRepo.GetActiveUserByPhone(ctx, input.Phone)

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
	err = srv.userRepo.CreateUser(ctx, newUser)

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
		return nil, nil, fmt.Errorf("%s: both email and password are required", op)
	}

	if err := utils.ValidateEmail(input.Email); err != nil {
		log.Printf("[%s] error validating email '%s': %v", op, input.Email, err)
		return nil, nil, err
	}

	// Check if user exists
	existingUser, err := srv.userRepo.GetUserByEmailIncludingDeleted(ctx, input.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			log.Printf("[%s] user not found: %v", op, err)
			return nil, nil, apperrors.ErrUserNotFound
		}
		log.Printf("[%s] get user by email failed: %v", op, err)
		return nil, nil, err
	}

	if existingUser.DeletedAt != nil {
		log.Printf("[%s] account deactivated: %s", op, existingUser.Email)
		return nil, nil, apperrors.ErrAccountDeactivated
	}

	// Verify password
	if err := utils.VerifyPassword(input.Password, existingUser.PasswordHash); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, nil, apperrors.ErrInvalidCredentials
		}
		return nil, nil, err // unexpected crypto failure
	}

	// Generate JWT access token
	accessToken, _, err := srv.jwtMaker.CreateToken(existingUser.ID, existingUser.ID, existingUser.Email, "owner", 25 *time.Minute)
	if err != nil {
		log.Printf("[%s] failed to create JWT: %v", op, err)
		return nil, nil, fmt.Errorf("%s: failed to issue access token: %w", op, err)
	}

	// Generate the very first Refresh Token
	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		log.Printf("[%s] failed to generate refresh token: %v", op, err)
		return nil, nil, err
	}
	log.Printf("[%s] refresh token created: %v", op, newRefreshToken)

	// Store hashed refresh token in DB
	newTokenEntity := &entity.RefreshToken{
		UserId:    existingUser.ID,
		TokenHash: auth.HashToken(newRefreshToken),
		ExpiresAt: time.Now().Add(refreshTokenExpiry),
		Revoked:   false,
	}
	if err := srv.authRepo.CreateRefreshToken(ctx, newTokenEntity); err != nil {
		log.Printf("[%s] failed to store refresh token: %v", op, err)
		return nil, nil, err
	}

	token := &entity.Token{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(25 * time.Minute),
		ExpiresIn:    int64((25 * time.Minute)),
		IssuedAt:     time.Now(),
		TokenType:    "Bearer",
	}

	return existingUser, token, nil
}

func (srv *AuthService) RefreshToken(ctx context.Context, incomingRefresh string) (*entity.Token, error) {
	const op = "AuthService.RefreshUserSession"

	if incomingRefresh == "" {
		return nil, apperrors.ErrInvalidRefreshToken
	}

	log.Printf("[%s] incoming refresh token: %v", op, incomingRefresh)

	// Hash incoming token to compare with DB
	hash := auth.HashToken(incomingRefresh)
	log.Printf("[%s] hashed refresh token: %v", op, hash)

	// Lookup token
	tokenEntity, err := srv.authRepo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		log.Printf("[%s] refresh token lookup failed: %v", op, err)
		return nil, apperrors.ErrInvalidRefreshToken
	}

	// Validate token
	if tokenEntity.Revoked || time.Now().After(tokenEntity.ExpiresAt) {
		log.Printf("Access denied: Token %s was previously revoked", hash)
		return nil, apperrors.ErrInvalidRefreshToken
	}

	// Revoke old token
	_, err = srv.authRepo.RevokeRefreshToken(ctx, tokenEntity.Id.String())
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newTokenEntity := &entity.RefreshToken{
		UserId:    tokenEntity.UserId,
		TokenHash: auth.HashToken(newRefreshToken),
		ExpiresAt: time.Now().Add(refreshTokenExpiry),
		Revoked:   false,
	}

	if err = srv.authRepo.CreateRefreshToken(ctx, newTokenEntity); err != nil {
		return nil, err
	}

	// Create new JWT access token
	accessToken, _, err := srv.jwtMaker.CreateToken(
		tokenEntity.UserId,
		tokenEntity.UserId,
		"",
		"owner",
		30*time.Minute,
	)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Africa/Nairobi")
	now := time.Now().In(loc)

	return &entity.Token{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: newRefreshToken, // raw token to client
		ExpiresAt:    now.Add(30 * time.Minute),
		ExpiresIn:    int64((30 * time.Minute).Seconds()),
		IssuedAt:     now,
		TokenType:    "Bearer",
	}, nil
}

// generateRefreshToken creates a new random refresh token string
func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32) // 256-bit random
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// base64 URL-safe encoding
	return base64.RawURLEncoding.EncodeToString(b), nil
}
