package user

import (
	"context"
	"errors"
	"testing"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/apperrors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/stretchr/testify/assert"
)

func TestFetchAllUsers_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := NewUserService(mockRepo)

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
	srv := NewUserService(mockRepo)

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
	srv := NewUserService(mockRepo)

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
	srv := NewUserService(mockRepo)

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
	srv := NewUserService(mockRepo)

	ctx := context.Background()
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	mockRepo.On("GetUserById", ctx, id).Return(&entity.User{}, apperrors.ErrUserNotFound)

	res, err := srv.GetCurrentUserById(ctx, id)
	assert.Nil(t, res)
	assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
}

func TestUpdateCurrentUser(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := NewUserService(mockRepo)

	ctx := context.Background()
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	input := &PatchUserInput{
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
	srv := NewUserService(mockRepo)

	ctx := context.Background()
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}
	input := &PatchUserInput{Id: id}

	mockRepo.On("PatchUser", ctx, id, &repository.UpdateUserInput{}).Return(&entity.User{}, apperrors.ErrUserNotFound)

	res, err := srv.UpdateCurrentUser(ctx, input)
	assert.Nil(t, res)
	assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
}

func TestSoftDeleteUser(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := NewUserService(mockRepo)

	ctx := context.Background()
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}

	mockRepo.On("DeleteUser", ctx, id).Return(nil)
	err := srv.SoftDeleteUser(ctx, id)
	assert.NoError(t, err)
}

func TestSoftDeleteUser_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	srv := NewUserService(mockRepo)

	ctx := context.Background()
	id := pgtype.UUID{Bytes: [16]byte{1}, Valid: true}

	mockRepo.On("DeleteUser", ctx, id).Return(apperrors.ErrUserNotFound)
	err := srv.SoftDeleteUser(ctx, id)
	assert.ErrorIs(t, err, apperrors.ErrUserNotFound)
}

// helper for pointer literals
func ptr[T any](v T) *T { return &v }
