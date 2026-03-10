package clients

import (
	"context"
	"time"

	membershipv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/membership/v1"
	"google.golang.org/grpc"
)

type UserServiceClient struct {
	client membershipv1.MembershipServiceClient
}

func NewUserServiceClient(conn *grpc.ClientConn) *UserServiceClient {
	return &UserServiceClient{
		client: membershipv1.NewMembershipServiceClient(conn),
	}
}

func (c *UserServiceClient) LinkUserToTenant(ctx context.Context, req *membershipv1.LinkUserToTenantRequest) (*membershipv1.LinkUserToTenantResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return c.client.LinkUserToTenant(ctx, req)
}
