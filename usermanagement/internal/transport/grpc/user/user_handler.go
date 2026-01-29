package user

import (

	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/user"
	"github.com/MartinMurithi/storeforge/usermanagement/proto/user/v1"
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

// func (u *UserGrpcHandler) GetCurrentUser(ctx context.Context, req *userv1.GetUserByIdRequest) (userv1.GetUserByIdResponse, error) {


// }
