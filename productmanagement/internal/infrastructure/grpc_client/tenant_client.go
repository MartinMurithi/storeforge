package grpcclient

import (
	"context"
	"time"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	"google.golang.org/grpc"
)

type TenantSvcClient struct {
	TenantClient tenantv1.TenantServiceClient
}

func NewTenantSvcClient(conn *grpc.ClientConn) *TenantSvcClient {
	return &TenantSvcClient{
		TenantClient: tenantv1.NewTenantServiceClient(conn),
	}
}

func (c *TenantSvcClient) GetTenantContext(ctx context.Context, req *tenantv1.GetTenantContextRequest) (*tenantv1.GetTenantContextResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return c.GetTenantContext(ctx, req)
}
