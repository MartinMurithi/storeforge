package mapper

import (
	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
)

func MapCreateTenantResponseProtoToDTO(pbRes *tenantv1.CreateTenantResponse) dto.CreateTenantResponseDTO {
	return dto.CreateTenantResponseDTO{
		Message: "Store Created Successfully",
		Tenant: dto.TenantDTO{
			ID:           pbRes.Tenant.Id,
			StoreName:    pbRes.Tenant.StoreName,
			Slug:         pbRes.Tenant.Slug,
			SubDomain:    pbRes.Tenant.SubDomain,
			Domain:       pbRes.Tenant.Domain,
			BusinessType: pbRes.Tenant.BusinessType,
			Status:       pbRes.Tenant.Status,
			CreatedAt:    pbRes.Tenant.CreatedAt.AsTime(),
		},
		Theme: dto.ThemeDTO{
			ID:          pbRes.Theme.Id,
			Name:        pbRes.Theme.Name,
			Description: pbRes.Theme.Description,
			Config:      pbRes.Theme.ThemeConfig, // Assumes map[string]interface{}
			IsActive:    pbRes.Theme.IsActive,
		},
		Settings: dto.SettingsDTO{
			TenantID:    pbRes.Settings.TenantId,
			ThemeID:     pbRes.Settings.ThemeId,
			ThemeConfig: pbRes.Settings.ThemeConfig,
			Version:     pbRes.Settings.Version,
			UpdatedAt:   pbRes.Settings.Updated.AsTime(),
		},
		Token: dto.TokenDTO{
			AccessToken: pbRes.Token.AccessToken,
			TokenType:   pbRes.Token.TokenType,
			ExpiresIn:   pbRes.Token.ExpiresIn / 1e9, // Convert nanos to seconds if needed
			ExpiresAt:   pbRes.Token.ExpiresAt.AsTime(),
		},
	}
}

// ToProtoGetContextTenantResponse composes multiple entity mappers into a single response DTO.
func MapGetTenantTenantContextResponse(pbRes *tenantv1.GetTenantContextResponse) dto.TenantContextResponseDTO {
	return dto.TenantContextResponseDTO{
		Tenant: dto.TenantDTO{
			ID:           pbRes.Tenant.Id,
			StoreName:    pbRes.Tenant.StoreName,
			Slug:         pbRes.Tenant.Slug,
			SubDomain:    pbRes.Tenant.SubDomain,
			Domain:       pbRes.Tenant.Domain,
			BusinessType: pbRes.Tenant.BusinessType,
			Status:       pbRes.Tenant.Status,
			CreatedAt:    pbRes.Tenant.CreatedAt.AsTime(),
		},
		Settings: dto.SettingsDTO{
			TenantID:    pbRes.Settings.TenantId,
			ThemeID:     pbRes.Settings.ThemeId,
			ThemeConfig: pbRes.Settings.ThemeConfig,
			Version:     pbRes.Settings.Version,
			UpdatedAt:   pbRes.Settings.Updated.AsTime(),
		},
		Role: pbRes.Role,
	}
}
