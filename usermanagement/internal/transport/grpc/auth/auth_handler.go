package auth

import (
	"context"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/auth"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/transport/grpc/grpc_errors"
	authv1 "github.com/MartinMurithi/storeforge/usermanagement/protos/auth/v1"
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
		return nil, grpc_errors.MapGrpcError(err)
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
		return nil, grpc_errors.MapGrpcError(err)
	}

	return ToProtoLoginResponse(user, token), nil

}
