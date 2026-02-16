package client

import (
	"context"
	"log"
	"time"

	userv1 "github.com/MartinMurithi/storeforge/api/protos/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	Service userv1.UserServiceClient
}

func NewUserClient(addr string) *UserClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to Auth Service: %v", err)
	}

	return &UserClient{
		Service: userv1.NewUserServiceClient(conn),
	}
}

func (c *UserClient) GetCurrentUser(ctx context.Context, req *userv1.GetCurrentUserRequest) (*userv1.GetCurrentUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.GetCurrentUser(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *UserClient) GetAllUsers(ctx context.Context, req *userv1.GetAllUsersRequest) (*userv1.GetAllUsersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.GetAllUsers(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *UserClient) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.UpdateUser(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *UserClient) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.DeleteUser(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
