package entity

import (
	"time"

	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"

	"github.com/google/uuid"
)

// NewTenant is a Domain Factory that ensures business rules are met upon creation.
func NewTenant(name, slug, business_type string, defaultThemeConfig map[string]any) (*Tenant, error) {
	if name == "" {
		return nil, apperrors.ErrBusinessNameRequired
	}

	tenantId := value_object.NewTenantIDFromUUID(uuid.New())
	themeId := value_object.NewThemeIDFromUUID(uuid.New())

	return &Tenant{
		ID:           tenantId,
		StoreName:    name,
		SubDomain:    slug + ".storeforge.com",
		BusinessType: business_type,
		Status:       "provisioning",
		CreatedAt:    time.Now(),
		Settings: &Settings{
			TenantID:  tenantId,
			ThemeID:   themeId,
			Config:    defaultThemeConfig,
			Version:   1,
			UpdatedAt: time.Now(),
		},
	}, nil
}
