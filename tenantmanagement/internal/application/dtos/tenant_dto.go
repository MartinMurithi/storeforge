package dtos

import "github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/entity"

// CreateTenantRequestDTO carries onboarding data from the transport layer.
type CreateTenantRequestDTO struct {
	StoreName    string
	Slug         string
	BusinessType string
	SubDomain    string
	ThemeID      string
}

// CreateTenantResponseDTO combines the created tenant and its chosen theme.
type CreateTenantResponseDTO struct {
	Tenant *entity.Tenant
	Theme *entity.Theme
	TenantSettings *entity.Settings
}

