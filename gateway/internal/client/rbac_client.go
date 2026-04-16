package client

import (
	"context"
	"log"
	"time"

	rbacv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/rbac/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RbacClient struct {
	Service rbacv1.RbacServiceClient
}

func NewRbacClient(addr string) *RbacClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to RBAC Service: %v", err)
	}

	return &RbacClient{
		Service: rbacv1.NewRbacServiceClient(conn),
	}
}

func (c *RbacClient) CreateRole(ctx context.Context, req *rbacv1.CreateRoleRequest) (*rbacv1.CreateRoleResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.CreateRole(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *RbacClient) GetRoleById(ctx context.Context, req *rbacv1.GetRoleByIDRequest) (*rbacv1.GetRoleByIDResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.GetRoleByID(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *RbacClient) UpdateRole(ctx context.Context, req *rbacv1.UpdateRoleRequest) (*rbacv1.UpdateRoleResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.UpdateRole(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
