package tests

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/auth"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/user"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/token"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFetchAllUsers_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := user.NewUserService(mockRepo)

	ctx := context.Background()
	p := dto.Pagination{Page: 1, Limit: 2}

	users := []*entity.User{
		{FullName: "Alice"},
		{FullName: "Bob"},
	}

	mockRepo.On("GetAllUsers", ctx, p).Return(users, 8, nil)

	// Internally, FetchAllUsers calls "GetAllUsers"
	result, meta, err := srv.FetchAllUsers(ctx, p)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 4, meta.TotalPages) // 5 users, limit 2 -> 4 pages
	assert.True(t, meta.HasNext)
	assert.False(t, meta.HasPrev)
}

func TestFetchAllUsers_Empty(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := user.NewUserService(mockRepo)

	ctx := context.Background()
	p := dto.Pagination{Page: 1, Limit: 10}

	mockRepo.On("GetAllUsers", ctx, p).Return([]*entity.User{}, 0, nil)

	result, meta, err := srv.FetchAllUsers(ctx, p)

	assert.NoError(t, err)
	assert.Len(t, result, 0)

	assert.Equal(t, 1, meta.Page)
	assert.Equal(t, 10, meta.Limit)
	assert.Equal(t, 0, meta.Total)
	assert.Equal(t, 0, meta.TotalPages)
	assert.False(t, meta.HasNext)
	assert.False(t, meta.HasPrev)
}

func TestFetchAllUsers_RepoError(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := user.NewUserService(mockRepo)

	ctx := context.Background()
	p := dto.Pagination{Page: 1, Limit: 10}

	mockRepo.On("GetAllUsers", ctx, p).Return(nil, 0, errors.New("db failure"))

	result, meta, err := srv.FetchAllUsers(ctx, p)

	assert.Nil(t, result)
	assert.Equal(t, dto.PaginationMeta{}, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db failure")
}

func TestGetCurrentUserById(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := user.NewUserService(mockRepo)

	ctx := context.Background()
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	user := &entity.User{FullName: "Alice"}
	mockRepo.On("GetUserById", ctx, id).Return(user, nil)

	res, err := srv.GetCurrentUserById(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", res.FullName)
}

func TestGetCurrentUserById_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := user.NewUserService(mockRepo)

	ctx := context.Background()
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	mockRepo.On("GetUserById", ctx, id).Return(&entity.User{}, apperrors.ErrUserNotFound)

	res, err := srv.GetCurrentUserById(ctx, id)
	assert.Nil(t, res)
	assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
}

// func TestUpdateCurrentUser(t *testing.T) {
// 	mockRepo := new(MockRepository)
// 	srv := user.NewUserService(mockRepo)

// 	ctx := context.Background()
// 	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
// 	input := &user.PatchUserInput{
// 		Id:           id,
// 		Email: ptr("New Name"),
// 		Phone: ptr("New Type"),
// 	}

// 	patchInput := &repository.UpdateUserInput{
// 		Email: input.,
// 		Phone: input.BusinessType,
// 	}

// 	updated := &entity.User{FullName: "Alice"}
// 	mockRepo.On("PatchUser", ctx, id, patchInput).Return(updated, nil)

// 	res, err := srv.UpdateCurrentUser(ctx, input)
// 	assert.NoError(t, err)
// 	assert.Equal(t, updated, res)
// }

func TestUpdateCurrentUser_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := user.NewUserService(mockRepo)

	ctx := context.Background()
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	input := &user.PatchUserInput{Id: id}

	mockRepo.On("PatchUser", ctx, id, &repository.UpdateUserInput{}).Return(&entity.User{}, apperrors.ErrUserNotFound)

	res, err := srv.UpdateCurrentUser(ctx, input)
	assert.Nil(t, res)
	assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
}

func TestSoftDeleteUser(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := user.NewUserService(mockRepo)

	ctx := context.Background()
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}

	mockRepo.On("DeleteUser", ctx, id).Return(nil)
	err := srv.SoftDeleteUser(ctx, id)
	assert.NoError(t, err)
}

func TestSoftDeleteUser_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := user.NewUserService(mockRepo)

	ctx := context.Background()
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}

	mockRepo.On("DeleteUser", ctx, id).Return(apperrors.ErrUserNotFound)
	err := srv.SoftDeleteUser(ctx, id)
	assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
}

func TestGetUserByEmailIncludingDeleted(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	deletedUser := &entity.User{
		ID:        pgtype.UUID{Bytes: uuid.New()},
		Email:     "deleted@example.com",
		DeletedAt: &time.Time{},
	}

	mockRepo.On("GetUserByEmailIncludingDeleted", ctx, "deleted@example.com").
		Return(deletedUser, nil)

	user, err := mockRepo.GetUserByEmailIncludingDeleted(ctx, "deleted@example.com")
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, user.DeletedAt)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByIdIncludingDeleted(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	id := pgtype.UUID{Bytes: uuid.New()}
	deletedUser := &entity.User{
		ID:        id,
		DeletedAt: &time.Time{},
	}

	mockRepo.On("GetUserByIdIncludingDeleted", ctx, id).
		Return(deletedUser, nil)

	user, err := mockRepo.GetUserByIdIncludingDeleted(ctx, id)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, user.DeletedAt)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByPhoneIncludingDeleted(t *testing.T) {
	mockRepo := new(MockRepository)
	ctx := context.Background()
	phone := "+254700123456"
	deletedUser := &entity.User{
		Phone:     phone,
		DeletedAt: &time.Time{},
	}

	mockRepo.On("GetUserByPhoneIncludingDeleted", ctx, phone).
		Return(deletedUser, nil)

	user, err := mockRepo.GetUserByPhoneIncludingDeleted(ctx, phone)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, user.DeletedAt)

	mockRepo.AssertExpectations(t)
}

// ====================== AUTH TESTS ===========================
func TestRegisterUser_Success(t *testing.T) {
	repo := new(MockRepository)
	jwtMaker := &token.JWTMaker{}
	srv := auth.NewAuthService(&repository.UserRepository{}, &repository.AuthRepository{}, jwtMaker)

	ctx := context.Background()
	input := &dto.RegisterUserRequestDTO{
		FullName:     "Alice Doe",
		Email:        "alice@example.com",
		Phone:        "+254700000000",
		Password:     "password123",
	}

	repo.On("GetUserByEmail", ctx, input.Email).Return(nil, nil)
	repo.On("GetUserByPhone", ctx, input.Phone).Return(nil, nil)
	repo.On("CreateUser", ctx, mock.AnythingOfType("*entity.User")).Return(nil)

	user, err := srv.RegisterUser(ctx, input)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, input.Email, user.Email)
	repo.AssertExpectations(t)
}

func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	repo := new(MockRepository)
	srv := auth.NewAuthService(&repository.UserRepository{}, &repository.AuthRepository{}, &token.JWTMaker{})

	ctx := context.Background()

	input := &dto.RegisterUserRequestDTO{
		FullName:     "Alice Doe",
		Email:        "alice@example.com",
		Phone:        "+254700000000",
		Password:     "password123",
	}

	repo.On("GetUserByEmail", ctx, input.Email).
		Return(&entity.User{Email: input.Email}, nil)

	user, err := srv.RegisterUser(ctx, input)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, apperrors.ErrUserAlreadyExists)
}

func TestRegisterUser_PhoneAlreadyExists(t *testing.T) {
	repo := new(MockRepository)
	srv := auth.NewAuthService(&repository.UserRepository{}, &repository.AuthRepository{}, &token.JWTMaker{})

	ctx := context.Background()

	input := &dto.RegisterUserRequestDTO{
		FullName:     "Alice Doe",
		Email:        "alice@example.com",
		Phone:        "+254700000000",
		Password:     "password123",
	}

	repo.On("GetUserByEmail", ctx, input.Email).Return(nil, nil)
	repo.On("GetUserByPhone", ctx, input.Phone).
		Return(&entity.User{Phone: input.Phone}, nil)

	user, err := srv.RegisterUser(ctx, input)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, apperrors.ErrUserAlreadyExists)
}

func TestRegisterUser_MissingEmail(t *testing.T) {
	repo := new(MockRepository)
	srv := auth.NewAuthService(&repository.UserRepository{}, &repository.AuthRepository{}, &token.JWTMaker{})

	input := &dto.RegisterUserRequestDTO{
		FullName:     "Alice Doe",
		Email:        "",
		Phone:        "+254700000000",
		Password:     "password123",
	}

	repo.On("GetUserByEmail", mock.Anything, input.Email).
		Return(input, nil)

	user, err := srv.RegisterUser(context.Background(), input)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, apperrors.ErrEmailRequired)
}

func TestLoginUser_Success(t *testing.T) {
	repo := new(MockRepository)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	jwtMaker, _ := token.NewJWTMaker(privateKey)
	srv := auth.NewAuthService(&repository.UserRepository{}, &repository.AuthRepository{}, jwtMaker)

	password := "password123"
	hashed, _ := utils.Hashpassword(password)

	user := &entity.User{
		ID:           pgtype.UUID{Valid: true},
		Email:        "alice@example.com",
		PasswordHash: hashed,
	}

	input := &dto.LoginUserRequestDTO{
		Email:    user.Email,
		Password: password,
	}

	repo.On("GetUserByEmail", mock.Anything, input.Email).
		Return(user, nil)

	u, tok, err := srv.LoginUser(context.Background(), input)

	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.NotNil(t, tok)
}

func TestLoginUser_InvalidPassword(t *testing.T) {
	repo := new(MockRepository)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	jwtMaker, _ := token.NewJWTMaker(privateKey)
	srv := auth.NewAuthService(&repository.UserRepository{}, &repository.AuthRepository{}, jwtMaker)

	hashed, _ := utils.Hashpassword("correct-password")

	user := &entity.User{
		Email:        "alice@example.com",
		PasswordHash: hashed,
	}

	repo.On("GetUserByEmail", mock.Anything, user.Email).
		Return(user, nil)

	input := &dto.LoginUserRequestDTO{
		Email:    user.Email,
		Password: "wrong-password",
	}

	u, tok, err := srv.LoginUser(context.Background(), input)

	assert.Nil(t, u)
	assert.Nil(t, tok)
	assert.ErrorIs(t, err, apperrors.ErrInvalidCredentials)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	repo := new(MockRepository)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	jwtMaker, _ := token.NewJWTMaker(privateKey)
	srv := auth.NewAuthService(&repository.UserRepository{}, &repository.AuthRepository{}, jwtMaker)

	repo.On("GetUserByEmail", mock.Anything, "missing@example.com").
		Return(nil, apperrors.ErrUserNotFound)

	input := &dto.LoginUserRequestDTO{
		Email:    "missing@example.com",
		Password: "password",
	}

	u, tok, err := srv.LoginUser(context.Background(), input)

	assert.Nil(t, u)
	assert.Nil(t, tok)
	assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
}

// helper for pointer literals
func ptr[T any](v T) *T { return &v }


// func TestUpdateSessionContext_Success(t *testing.T) {
// 	// 1. Setup
// 	mockAuthRepo := new(MockRepository)
// 	mockUserRepo := new(MockRepository) // Assume similar mock for User

// 	cfg, err := config.Load()
// 	// Use your real JWTMaker to verify the outcome
// 	jwtMaker := *&token.JWTMaker{PrivateKey: cfg.JWT.PrivateKeyPath}

// 	srv := NewAuthService(mockAuthRepo, mockUserRepo, jwtMaker)

// 	// 2. Data
// 	uID := uuid.New()
// 	tID := uuid.New()
// 	userEmail := "martin@storeforge.com"

// 	targetUserID := pgtype.UUID{Bytes: uID, Valid: true}
// 	targetTenantID := pgtype.UUID{Bytes: tID, Valid: true}

// 	input := &dto.UpdateActiveSessionContextRequestDTO{
// 		UserId:   targetUserID,
// 		TenantId: targetTenantID,
// 		Role:     "owner",
// 	}

// 	// 3. Set Expectations
// 	// When the service calls the repo, the repo returns this "fake" result
// 	mockAuthRepo.On("UpdateActiveSessionContext", mock.Anything, targetUserID, targetTenantID, "owner").
// 		Return(&entity.RefreshToken{
// 			UserId:       targetUserID,
// 			LastTenantId: targetTenantID,
// 			LastRole:     "owner",
// 		}, nil)

// 	mockUserRepo.On("GetActiveUserById", mock.Anything, targetUserID).
// 		Return(&entity.User{
// 			ID:    uID,
// 			Email: userEmail,
// 		}, nil)

// 	// 4. Execute
// 	token, err := srv.UpdateSessionContext(context.Background(), input)

// 	// 5. Assertions
// 	require.NoError(t, err)
// 	require.NotNil(t, token)

// 	// Verify the JWT actually has the data!
// 	payload, err := jwtMaker.VerifyToken(token.AccessToken)
// 	require.NoError(t, err)

// 	assert.Equal(t, uID, payload.UserID)
// 	assert.Equal(t, tID.String(), payload.TenantID) // Proof that promotion worked
// 	assert.Equal(t, "owner", payload.Role)

// 	// Ensure the mocks were actually used
// 	mockAuthRepo.AssertExpectations(t)
// 	mockUserRepo.AssertExpectations(t)
// }