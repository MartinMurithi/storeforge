package user

import (
	"context"
	"fmt"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/user"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/transport/grpc/grpc_errors"
	"github.com/MartinMurithi/storeforge/usermanagement/proto/user/v1"

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
		return nil, grpc_errors.MapGrpcError(err)
	}

	// map service layer users → proto Users
	var protoUsers []*entity.User

	for _, u := range users {
		protoUsers = append(protoUsers, u)
	}

	return ToProtoFetchAllUsersResponse(protoUsers, &meta), nil

}

func (h *UserGrpcHandler) UpdateUser(
	ctx context.Context,
	req *userv1.UpdateUserRequest,
) (*userv1.UpdateUserResponse, error) {

	uuid := pgtype.UUID{}
	if err := uuid.Scan(req.Id); err != nil {
		return nil, grpc_errors.MapGrpcError(err)
	}

	input := &user.PatchUserInput{
		Id: uuid,
	}

	if req.BusinessName != nil {
		input.BusinessName = req.BusinessName
	}

	if req.BusinessType != nil {
		input.BusinessType = req.BusinessType
	}

	updatedUser, err := h.UserService.UpdateCurrentUser(ctx, input)

	if err != nil {
		fmt.Printf("[UserGrpcHandler.UpdateUser] failed: %v\n", err)
		return nil, grpc_errors.MapGrpcError(err)
	}

	return ToProtoUpdateUserResponse(updatedUser), nil
}
