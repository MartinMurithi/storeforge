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
	roleRepo repository.IRoleRepository
	jwtMaker *token.JWTMaker
}

const refreshTokenExpiry = 30 * 24 * time.Hour // 30 days

// create a factory function to initialize my service with repo
func NewAuthService(userRepo repository.IUserRepository, authRepo repository.IAuthRepository, roleRepo repository.IRoleRepository, jwtMaker *token.JWTMaker) *AuthService {

	if jwtMaker == nil {
		panic("jwt maker must not be nil")
	}

	return &AuthService{
		userRepo: userRepo,
		authRepo: authRepo,
		roleRepo: roleRepo,
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

	// 1. Initialize empty/global context defaults
	var tenantID pgtype.UUID
	tenantID.Valid = false
	roleName := ""

	// 2. Fetch last active tenant + roleID (Only call this once)
	lastTenant, lastRoleID, err := srv.userRepo.GetLastActiveTenant(ctx, existingUser.ID)
	if err != nil {
		// We log but don't block login; the user just gets a global context
		log.Printf("[%s] failed to fetch last tenant: %v", op, err)
	}

	// 3. Handle Branching: If a tenant context exists, resolve it
	if lastTenant != nil && lastTenant.Valid {
		tenantID = *lastTenant

		// Resolve Role ID → Role Name string for the JWT
		if lastRoleID != nil && lastRoleID.Valid {
			roleEntity, err := srv.roleRepo.GetRoleByID(ctx, *lastRoleID)
			if err != nil {
				log.Printf("[%s] failed to fetch role name for ID %v: %v", op, *lastRoleID, err)
			} else {
				roleName = roleEntity.Name
			}
		}
	}

	// 4. Generate JWT with the resolved context (Tenant may be invalid/empty)
	accessToken, _, err := srv.jwtMaker.CreateToken(
		existingUser.ID,
		tenantID,
		existingUser.Email,
		roleName,
		25*time.Minute,
	)
	if err != nil {
		log.Printf("[%s] failed to create JWT: %v", op, err)
		return nil, nil, fmt.Errorf("%s: failed to issue access token: %w", op, err)
	}

	// 5. Generate the raw Refresh Token string
	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		log.Printf("[%s] failed to generate refresh token: %v", op, err)
		return nil, nil, err
	}

	// 6. Store tenant + role context in the Refresh Token record
	// This ensures that when the user "refreshes", they stay in the same store
	newTokenEntity := &entity.RefreshToken{
		UserId:       existingUser.ID,
		TokenHash:    auth.HashToken(newRefreshToken),
		ExpiresAt:    time.Now().Add(refreshTokenExpiry),
		Revoked:      false,
		LastRole:     roleName,
		LastTenantId: tenantID, // pgtype.UUID handles the "Invalid" state automatically
	}

	if err := srv.authRepo.CreateRefreshToken(ctx, newTokenEntity); err != nil {
		log.Printf("[%s] failed to store refresh token: %v", op, err)
		return nil, nil, err
	}

	// 7. Build the final response DTO
	token := &entity.Token{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(25 * time.Minute),
		ExpiresIn:    int64((25 * time.Minute).Seconds()),
		IssuedAt:     time.Now(),
		TokenType:    "Bearer",
	}

	log.Printf("[%s] issued token: user=%s, tenant=%v, role=%s", op, existingUser.Email, tenantID.Valid, roleName)

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

	dbSession, err := srv.authRepo.UpdateActiveSessionContext(ctx, input.UserId, input.TenantId, input.Role)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to update active session context in DB: %w", op, err)
	}

	user, err := srv.userRepo.GetActiveUserById(ctx, input.UserId)

	if err != nil {
		return nil, fmt.Errorf("%s: user lookup failed: %w", op, err)
	}

	// ISSUE THE NEW JWT
	// This creates a BRAND NEW Access Token with the TenantID and Role inside it.
	accessToken, _, err := srv.jwtMaker.CreateToken(
		user.ID,
		dbSession.LastTenantId,
		user.Email,
		dbSession.LastRole,
		30*time.Minute,
	)

	// log.Printf("claims from updated JWT token %s", claims.TenantId)

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
