package dtos

import (
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/entity"
)

// CreateTenantRequestDTO carries onboarding data from the transport layer.
type CreateTenantRequestDTO struct {
	StoreName    string
	BusinessType string
	ThemeID      string
	UserId       string
}

// New JWT token with Tenant Id and role
type TokenInfoDTO struct {
	AccessToken string
	ExpiresAt   time.Time
	ExpiresIn   int64
	IssuedAt    time.Time
	TokenType   string
}

// CreateTenantResponseDTO combines the created tenant and its chosen theme.
type CreateTenantResponseDTO struct {
	Tenant         *entity.Tenant
	Theme          *entity.Theme
	TenantSettings *entity.Settings
	Token *TokenInfoDTO
}

