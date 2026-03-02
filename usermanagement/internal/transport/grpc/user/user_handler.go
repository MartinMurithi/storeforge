package user

import (
	"context"
	"fmt"

	userv1 "github.com/MartinMurithi/storeforge/api/protos/user/v1"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/user"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserGrpcHandler struct {
	UserService *user.UserService
	userv1.UnimplementedUserServiceServer
}

func NewUserGrpcHandler(userService *user.UserService) *UserGrpcHandler {
	return &UserGrpcHandler{
		UserService: userService,
	}
}

func (h *UserGrpcHandler) GetAllUsers(ctx context.Context, req *userv1.GetAllUsersRequest) (*userv1.GetAllUsersResponse, error) {

	pagination := &dto.Pagination{
		Page:  int(req.Page),
		Limit: int(req.Limit),
	}

	users, meta, err := h.UserService.FetchAllUsers(ctx, *pagination)

	if err != nil {
		fmt.Printf("[USERGRPCHANDLER]: failed to fetch users %s\n", err)
		return nil, errconv.ToGrpcError(err)
	}

	// map service layer users → proto Users
	var protoUsers []*entity.User

	for _, u := range users {
		protoUsers = append(protoUsers, u)
	}

	return ToProtoFetchAllUsersResponse(protoUsers, &meta), nil

}

func (h *UserGrpcHandler) GetCurrentUser(ctx context.Context, req *userv1.GetCurrentUserRequest) (*userv1.GetCurrentUserResponse, error) {

	uuid := pgtype.UUID{}
	if err := uuid.Scan(req.Id); err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	user, err := h.UserService.GetCurrentUserById(ctx, uuid)

	if err != nil {
		fmt.Printf("[USERGRPCHANDLER]: failed to fetch user %s\n", err)
		return nil, errconv.ToGrpcError(err)
	}

	return &userv1.GetCurrentUserResponse{
		User: ToProtoUser(user),
	}, nil

}

func (h *UserGrpcHandler) UpdateUser(
	ctx context.Context,
	req *userv1.UpdateUserRequest,
) (*userv1.UpdateUserResponse, error) {

	uuid := pgtype.UUID{}
	if err := uuid.Scan(req.Id); err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	input := &user.PatchUserInput{
		Id: uuid,
	}

	if req.Email != nil {
		input.Email = req.Email
	}

	if req.Phone != nil {
		input.Phone = req.Phone
	}

	updatedUser, err := h.UserService.UpdateCurrentUser(ctx, input)

	if err != nil {
		fmt.Printf("[UserGrpcHandler.UpdateUser] failed: %v\n", err)
		return nil, errconv.ToGrpcError(err)
	}

	return ToProtoUpdateUserResponse(updatedUser), nil
}

// DeleteUser, allows an admin to soft delete a user
func (h *UserGrpcHandler) DeleteUser(
	ctx context.Context,
	req *userv1.DeleteUserRequest,
) (*userv1.DeleteUserResponse, error) {

	uuid := pgtype.UUID{}
	if err := uuid.Scan(req.Id); err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	err := h.UserService.SoftDeleteUser(ctx, uuid)

	if err != nil {
		fmt.Printf("[UserGrpcHandler.SoftDeleteUser] failed: %v\n", err)
		return nil, errconv.ToGrpcError(err)
	}

	return ToProtoSoftDeleteUserResponse(), nil
}
