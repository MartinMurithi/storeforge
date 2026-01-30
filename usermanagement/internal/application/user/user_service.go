package user

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/apperrors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	repo repository.IUserRepository
}

// create a factory function to initialize my service with repo
func NewUserService(repo repository.IUserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

type PatchUserInput struct {
	Id           pgtype.UUID
	BusinessName *string
	BusinessType *string
}

// Admin role
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

// Admin Role
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
