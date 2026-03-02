package entity

import (
	"time"

	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"

	"github.com/google/uuid"
)

// NewTenant is a Domain Factory that ensures business rules are met upon creation.
func NewTenant(name, slug, business_type string) (*Tenant, error) {
	if name == "" {
		return nil, apperrors.ErrBusinessNameRequired
	}
	return &Tenant{
		// We use a random UUID for a brand new Tenant
		ID:           value_object.NewTenantIDFromUUID(uuid.New()),
		StoreName:    name,
		SubDomain:    slug + ".storeforge.com",
		BusinessType: business_type,
		Status:       "provisioning",
		CreatedAt:    time.Now(),
	}, nil
}
