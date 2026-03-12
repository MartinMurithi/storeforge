package clients

import (
	"context"
	"time"

	authv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/auth/v1"
	membershipv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/membership/v1"
	"google.golang.org/grpc"
)

type UserServiceClient struct {
	// To link user to their store
	MembershipClient          membershipv1.MembershipServiceClient
	UpdateActiveSessionClient authv1.AuthServiceClient
}

func NewUserServiceClient(conn *grpc.ClientConn) *UserServiceClient {
	return &UserServiceClient{
		MembershipClient:          membershipv1.NewMembershipServiceClient(conn),
		UpdateActiveSessionClient: authv1.NewAuthServiceClient(conn),
	}
}

func (c *UserServiceClient) LinkUserToTenant(ctx context.Context, req *membershipv1.LinkUserToTenantRequest) (*membershipv1.LinkUserToTenantResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return c.MembershipClient.LinkUserToTenant(ctx, req)
}

func (c *UserServiceClient) UpdateActiveSessionContext(ctx context.Context, req *authv1.UpdateSessionContextRequest) (*authv1.UpdateSessionContextResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return c.UpdateActiveSessionClient.UpdateSessionContext(ctx, req)
}
