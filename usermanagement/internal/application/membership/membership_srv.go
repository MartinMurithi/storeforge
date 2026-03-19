package membership

import (
	"context"
	"fmt"
	"log"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/repository"
	"github.com/google/uuid"
)

type MembershipService struct {
	membershipRepo repository.IMembershipRepository
}

func NewMembershipService(r repository.IMembershipRepository) *MembershipService {
	return &MembershipService{
		membershipRepo: r,
	}
}


func (s *MembershipService) LinkUserToTenant(ctx context.Context, userId, tenantId uuid.UUID, roleName string) error{

	const op = "MembershipService.LinkUserToTenant"

	roleName = entity.RoleOwner

	// Ensure only "owner role" is assigned
	if roleName != "owner"{
		return fmt.Errorf("invalid initial role assignment")
	}

	err := s.membershipRepo.AddTenantMembership(ctx, userId, tenantId, roleName)

	if err != nil{
		log.Printf("[%s]: an error occurred when linking user to tenant %v", op, err)
		return err
	}
	return nil
}