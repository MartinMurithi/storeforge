package tests

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/apperrors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/auth"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/user"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/token"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/utils"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestUpdateCurrentUser(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := user.NewUserService(mockRepo)

	ctx := context.Background()
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	input := &user.PatchUserInput{
		Id:           id,
		BusinessName: ptr("New Name"),
		BusinessType: ptr("New Type"),
	}

	patchInput := &repository.UpdateUserInput{
		BusinessName: input.BusinessName,
		BusinessType: input.BusinessType,
	}

	updated := &entity.User{FullName: "Alice"}
	mockRepo.On("PatchUser", ctx, id, patchInput).Return(updated, nil)

	res, err := srv.UpdateCurrentUser(ctx, input)
	assert.NoError(t, err)
	assert.Equal(t, updated, res)
}

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

// ====================== AUTH TESTS ===========================
func TestRegisterUser_Success(t *testing.T) {
	repo := new(MockRepository)
	jwtMaker := &token.JWTMaker{}
	srv := auth.NewAuthService(repo, jwtMaker)

	ctx := context.Background()
	input := &dto.RegisterUserRequestDTO{
		FullName:     "Alice Doe",
		Email:        "alice@example.com",
		Phone:        "+254700000000",
		Password:     "password123",
		BusinessType: "Retail",
		BusinessName: "Alice Shop",
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
	srv := auth.NewAuthService(repo, &token.JWTMaker{})

	ctx := context.Background()

	input := &dto.RegisterUserRequestDTO{
		FullName:     "Alice Doe",
		Email:        "alice@example.com",
		Phone:        "+254700000000",
		Password:     "password123",
		BusinessType: "Retail",
		BusinessName: "Alice Shop",
	}

	repo.On("GetUserByEmail", ctx, input.Email).
		Return(&entity.User{Email: input.Email}, nil)

	user, err := srv.RegisterUser(ctx, input)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, apperrors.ErrUserAlreadyExists)
}

func TestRegisterUser_PhoneAlreadyExists(t *testing.T) {
	repo := new(MockRepository)
	srv := auth.NewAuthService(repo, &token.JWTMaker{})

	ctx := context.Background()

	input := &dto.RegisterUserRequestDTO{
		FullName:     "Alice Doe",
		Email:        "alice@example.com",
		Phone:        "+254700000000",
		Password:     "password123",
		BusinessType: "Retail",
		BusinessName: "Alice Shop",
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
	srv := auth.NewAuthService(repo, &token.JWTMaker{})

	input := &dto.RegisterUserRequestDTO{
		FullName:     "Alice Doe",
		Email:        "",
		Phone:        "+254700000000",
		Password:     "password123",
		BusinessType: "Retail",
		BusinessName: "Alice Shop",
	}

	user, err := srv.RegisterUser(context.Background(), input)

	assert.Nil(t, user)
	assert.ErrorIs(t, err, apperrors.ErrEmailRequired)
}

func TestLoginUser_Success(t *testing.T) {
	repo := new(MockRepository)
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	jwtMaker, _ := token.NewJWTMaker(privateKey)
	srv := auth.NewAuthService(repo, jwtMaker)

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
	srv := auth.NewAuthService(repo, jwtMaker)

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
	srv := auth.NewAuthService(repo, jwtMaker)

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
