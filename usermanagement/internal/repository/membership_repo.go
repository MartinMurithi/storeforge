package repository

import (
	"context"
	"log"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/database"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database/postgres"
	"github.com/google/uuid"
)

type MembershipRepository struct {
	DB database.DB
}

type IMembershipRepository interface {
	AddTenantMembership(ctx context.Context, userId uuid.UUID, tenantId uuid.UUID, roleName string) error
}

func NewMembershipRepository(db database.DB) *MembershipRepository {
	return &MembershipRepository{
		DB: db,
	}
}

func (r *MembershipRepository) AddTenantMembership(ctx context.Context, userID, tenantID uuid.UUID, roleName string) error {
	const op = "UserRepository.AddTenantMembership"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	getRoleQuery := `SELECT id FROM ROLES WHERE name = $1`

	// Get the Role ID
	var roleID uuid.UUID
	err := r.DB.QueryRow(ctx, getRoleQuery, roleName).Scan(&roleID)

	if err != nil {
		log.Printf("[%s]: error fetching role ID %v", op, err)

		// Map DB/driver errors → infra errors
		infraErr := postgres.MapPostgresError(err)

		// Translate infra → domain errors
		return TranslateRoleRepoError(infraErr)
	}

	// Added ON CONFLICT to make the operation idempotent (safe to retry)
	query := `INSERT INTO users_tenants (user_id, tenant_id, role_id)
				VALUES ($1, $2, $3) 
				ON CONFLICT (user_id, tenant_id) DO NOTHING
				RETURNING user_id, tenant_id`

	err = r.DB.QueryRow(ctx, query, userID, tenantID, roleID).Scan(&userID, &tenantID)

	if err != nil {
		log.Printf("[%s]: error inserting membership %v", op, err)

		// Map DB/driver errors → infra errors
		infraErr := postgres.MapPostgresError(err)

		// Translate infra → domain errors
		return TranslateRoleRepoError(infraErr)
	}

	return err
}
