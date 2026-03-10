package client

import (
	"log"

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
