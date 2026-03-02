package auth

import (
	"context"
	"fmt"

	authv1 "github.com/MartinMurithi/storeforge/api/protos/auth/v1"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/auth"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (a *AuthGrpcHandler) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	const op = "AuthGrpcHandler.RefreshToken"

	if req.RefreshToken == "" {
		return nil, errconv.ToGrpcError(fmt.Errorf("%s: refresh token is required", op))
	}

	token, err := a.AuthService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	return &authv1.RefreshTokenResponse{
		Token: &authv1.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
			ExpiresAt:    timestamppb.New(token.ExpiresAt),
			ExpiresIn:    token.ExpiresIn,
			IssuedAt:     timestamppb.New(token.IssuedAt),
			TokenType:    token.TokenType,
		},
	}, nil
}

func (a *AuthGrpcHandler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	const op = "AuthGrpcHandler.Logout"

	if req.RefreshToken == "" {
		return nil, errconv.ToGrpcError(fmt.Errorf("%s: refresh token is required", op))
	}

	err := a.AuthService.Logout(ctx, req.RefreshToken)
	
	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	return &authv1.LogoutResponse{
		Success: true,
	}, nil
}
