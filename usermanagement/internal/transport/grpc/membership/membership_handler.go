package membership

import (
	"context"
	"log"

	membershipv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/membership/v1"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/membership"
	"github.com/google/uuid"
)

type MembershipGrpcHnadler struct {
	MembershipSrv *membership.MembershipService
	membershipv1.UnimplementedMembershipServiceServer
}

func NewMembershipGrpcHandler(m *membership.MembershipService) *MembershipGrpcHnadler {
	return &MembershipGrpcHnadler{
		MembershipSrv: m,
	}
}

func (h *MembershipGrpcHnadler) LinkUserToTenant(ctx context.Context, req *membershipv1.LinkUserToTenantRequest) (*membershipv1.LinkUserToTenantResponse, error) {
	uID, _ := uuid.Parse(req.UserId)
	tID, _ := uuid.Parse(req.TenantId)

	log.Printf("user id in user service %s", uID)
	log.Printf("tenant id user service %s", tID)

	err := h.MembershipSrv.LinkUserToTenant(ctx, uID, tID, req.Role)
	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	return &membershipv1.LinkUserToTenantResponse{
		Success: true,
		Message: "user linked to their store successfully",
	}, nil
}
