package client

import (
	"context"
	"log"
	"time"

	authv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	Service authv1.AuthServiceClient
}

func NewAuthClient(addr string) *AuthClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to Auth Service: %v", err)
	}

	return &AuthClient{
		Service: authv1.NewAuthServiceClient(conn),
	}
}

func (c *AuthClient) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.Register(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *AuthClient) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.Login(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *AuthClient) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.RefreshToken(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *AuthClient) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.Logout(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
