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

	log.Printf("connected to tenant svc successfully")

	return &TenantClient{
		Service: tenantv1.NewTenantServiceClient(conn),
	}

}

func (c *TenantClient) CreateTenant(ctx context.Context, req *tenantv1.CreateTenantRequest) (*tenantv1.CreateTenantResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.CreateTenant(ctx, req)

	log.Printf("tenant created successfully")

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *TenantClient) GetTenantContext(ctx context.Context, req *tenantv1.GetTenantContextRequest) (*tenantv1.GetTenantContextResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := c.Service.GetTenantContext(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *TenantClient) UpdateTenant(ctx context.Context, req *tenantv1.UpdateTenantRequest) (*tenantv1.GetTenantContextResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := c.Service.UpdateTenant(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}
