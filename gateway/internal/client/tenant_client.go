package client

import (
	"context"
	"log"
	"time"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TenantClient struct {
	Service tenantv1.TenantServiceClient
}

func NewTenantClient(addr string) *TenantClient {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("did not connect to Tenant Service: %v", err)
	}

	return &TenantClient{
		Service: tenantv1.NewTenantServiceClient(conn),
	}
}

func (c *TenantClient) CreateTenant(ctx context.Context, req *tenantv1.CreateTenantRequest) (*tenantv1.CreateTenantResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.CreateTenant(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
