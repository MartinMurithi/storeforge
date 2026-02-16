package auth

import (
	"context"

	authv1 "github.com/MartinMurithi/storeforge/api/protos/auth/v1"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/auth"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
)

type AuthGrpcHandler struct {
	AuthService *auth.AuthService
	authv1.UnimplementedAuthServiceServer
}

func NewAuthGrpcHandler(a *auth.AuthService) *AuthGrpcHandler {
	return &AuthGrpcHandler{
		AuthService: a,
	}
}

func (h *AuthGrpcHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {

	dtoReq := &dto.RegisterUserRequestDTO{
		FullName:     req.FullName,
		Email:        req.Email,
		Phone:        req.Phone,
		Password:     req.Password,
		BusinessType: req.BusinessType,
		BusinessName: req.BusinessName,
	}

	user, err := h.AuthService.RegisterUser(ctx, dtoReq)

	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	return ToProtoRegisterResponse(user), nil

}

func (a *AuthGrpcHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {

	dtoReq := &dto.LoginUserRequestDTO{
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := a.AuthService.LoginUser(ctx, dtoReq)

	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	return ToProtoLoginResponse(user, token), nil

}
