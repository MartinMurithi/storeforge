package user

import (
	"context"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

// CreateUser implements [repository.IUserRepository].
func (m *MockRepository) CreateUser(ctx context.Context, user *entity.User) error {
	panic("unimplemented")
}

// GetUserByEmail implements [repository.IUserRepository].
func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

// GetUserByPhone implements [repository.IUserRepository].
func (m *MockRepository) GetUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockRepository) GetAllUsers(ctx context.Context, p dto.Pagination) ([]*entity.User, int, error) {
	args := m.Called(ctx, p)

	var users []*entity.User
	if args.Get(0) != nil {
		users = args.Get(0).([]*entity.User)
	}

	return users, args.Int(1), args.Error(2)
}

func (m *MockRepository) GetUserById(ctx context.Context, id pgtype.UUID) (*entity.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockRepository) PatchUser(ctx context.Context, id pgtype.UUID, input *repository.UpdateUserInput) (*entity.User, error) {
	args := m.Called(ctx, id, input)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockRepository) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
