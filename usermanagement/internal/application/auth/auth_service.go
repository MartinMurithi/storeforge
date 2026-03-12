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
	"github.com/jackc/pgx/v5/pgtype"

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

	// To issue a token WITHOUT a tenant context:
	var emptyTenantID pgtype.UUID
	emptyTenantID.Valid = false

	// Generate JWT access token
	accessToken, claims, err := srv.jwtMaker.CreateToken(existingUser.ID, emptyTenantID, existingUser.Email, "", 25*time.Minute)
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
		UserId:       existingUser.ID,
		TokenHash:    auth.HashToken(newRefreshToken),
		ExpiresAt:    time.Now().Add(refreshTokenExpiry),
		Revoked:      false,
		LastRole:     "",
		LastTenantId: emptyTenantID,
	}
	if err := srv.authRepo.CreateRefreshToken(ctx, newTokenEntity); err != nil {
		log.Printf("[%s] failed to store refresh token: %v", op, err)
		return nil, nil, err
	}

	token := &entity.Token{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Local().Add(25 * time.Minute),
		ExpiresIn:    int64((25 * time.Minute)),
		IssuedAt:     time.Now().Local(),
		TokenType:    "Bearer",
	}

	log.Printf("issued access token with no role and tenant id %s %s", *claims.TenantId, claims.Role)

	return existingUser, token, nil
}

func (srv *AuthService) RefreshToken(ctx context.Context, incomingRefresh string) (*entity.Token, error) {
	const op = "AuthService.RefreshUserSession"

	if incomingRefresh == "" {
		return nil, apperrors.ErrInvalidRefreshToken
	}

	// Hash incoming token to compare with DB
	hash := auth.HashToken(incomingRefresh)

	oldToken, err := srv.authRepo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		log.Printf("[%s] refresh token lookup failed: %v", op, err)
		return nil, apperrors.ErrInvalidRefreshToken
	}

	// Validate token status and expiration
	if oldToken.Revoked || time.Now().After(oldToken.ExpiresAt) {
		log.Printf("[%s] access denied: token %s was revoked or expired", op, hash)
		return nil, apperrors.ErrInvalidRefreshToken
	}

	user, err := srv.userRepo.GetActiveUserById(ctx, oldToken.UserId)
	if err != nil {
		log.Printf("[%s] user lookup failed during refresh: %v", op, err)
		return nil, apperrors.ErrUserNotFound
	}

	// Revoke the old refresh token (Rotation)
	_, err = srv.authRepo.RevokeRefreshToken(ctx, oldToken.Id.String())
	if err != nil {
		return nil, err
	}

	// Generate a new rotating refresh token
	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	newTokenEntity := &entity.RefreshToken{
		UserId:       oldToken.UserId,
		TokenHash:    auth.HashToken(newRefreshToken),
		ExpiresAt:    time.Now().Add(refreshTokenExpiry),
		Revoked:      false,
		LastRole:     oldToken.LastRole,
		LastTenantId: oldToken.LastTenantId,
	}

	if err = srv.authRepo.CreateRefreshToken(ctx, newTokenEntity); err != nil {
		return nil, err
	}

	accessToken, _, err := srv.jwtMaker.CreateToken(
		user.ID,
		oldToken.LastTenantId,
		user.Email,
		oldToken.LastRole,
		30*time.Minute,
	)
	if err != nil {
		return nil, err
	}

	return &entity.Token{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: newRefreshToken, // raw token to client
		ExpiresAt:    time.Now().Add(30 * time.Minute),
		ExpiresIn:    int64((30 * time.Minute).Seconds()),
		IssuedAt:     time.Now(),
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

// Logout revokes the refresh token stored in the DB
func (srv *AuthService) Logout(ctx context.Context, refreshToken string) error {
	const op = "AuthService.Logout"

	if refreshToken == "" {
		return apperrors.ErrInvalidRefreshToken
	}

	log.Printf("[%s] incoming refresh token: %v", op, refreshToken)

	// Hash incoming token to compare with DB
	hash := auth.HashToken(refreshToken)
	log.Printf("[%s] hashed refresh token: %v", op, hash)

	log.Printf("[%s] hashed refresh token from DB: %v", op, hash)

	// Revoke old token
	_, err := srv.authRepo.RevokeRefreshTokenByHash(ctx, hash)

	if err != nil {
		return err
	}

	return nil
}

func (srv *AuthService) UpdateSessionContext(ctx context.Context, input *dto.UpdateActiveSessionContextRequestDTO) (*entity.Token, error) {
	const op = "AuthService.UpdateSessionContext"

	// We update the Refresh Token record so the session "remembers" this store.
	dbSession, err := srv.authRepo.UpdateActiveSessionContext(ctx, input.UserId, input.TenantId, input.Role)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to update active session context in DB: %w", op, err)
	}

	// We need the user's email to put it into the new JWT claims.
	user, err := srv.userRepo.GetActiveUserById(ctx, input.UserId)

	if err != nil {
		return nil, fmt.Errorf("%s: user lookup failed: %w", op, err)
	}

	// ISSUE THE NEW JWT
	// This creates a BRAND NEW Access Token with the TenantID and Role inside it.
	accessToken, claims, err := srv.jwtMaker.CreateToken(
		user.ID,
		dbSession.LastTenantId,
		user.Email,
		dbSession.LastRole,
		30*time.Minute,
	)

	log.Printf("claims from updated JWT token %s, %s", claims.TenantId, *claims.Role)

	if err != nil {
		return nil, fmt.Errorf("%s: token signing failed: %w", op, err)
	}

	token := &entity.Token{
		AccessToken: accessToken.AccessToken,
		ExpiresAt:   time.Now().Local().Add(25 * time.Minute),
		ExpiresIn:   int64((25 * time.Minute)),
		IssuedAt:    time.Now().Local(),
		TokenType:   "Bearer",
	}

	return token, nil
}
