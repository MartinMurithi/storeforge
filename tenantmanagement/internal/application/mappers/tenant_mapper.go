package mappers

import (
	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/application/dtos"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/entity"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoTenant maps the Tenant domain entity to its Protobuf representation.
func ToProtoTenant(t *entity.Tenant) *tenantv1.Tenant {
	if t == nil {
		return nil
	}
	return &tenantv1.Tenant{
		Id:           t.ID.String(),
		StoreName:    t.StoreName,
		Slug:         t.Slug,
		SubDomain:    t.SubDomain,
		Status: t.Status,
		BusinessType: t.BusinessType,
		CreatedAt:    timestamppb.New(t.CreatedAt),
	}
}


// ToProtoCreateTenantResponse composes multiple entity mappers into a single response DTO.
func ToProtoCreateTenantResponse(data *dtos.CreateTenantResponseDTO) *tenantv1.CreateTenantResponse {
	if data == nil || data.Tenant == nil {
		return nil
	}

	return &tenantv1.CreateTenantResponse{
		Tenant:   ToProtoTenant(data.Tenant),
		Theme:    ToProtoTheme(data.Theme),
		Settings: ToProtoSettings(data.Tenant.Settings),
	}
}

// ToProtoSettings maps the Tenant's active configuration instance.
func ToProtoSettings(s *entity.Settings) *tenantv1.TenantSettings {
	if s == nil {
		return nil
	}
	config, _ := structpb.NewStruct(s.Config)

	return &tenantv1.TenantSettings{
		TenantId:    s.TenantID.String(),
		ThemeId:     s.ThemeID.String(),
		ThemeConfig: config,
		Version:     int32(s.Version),
		Updated:     timestamppb.New(s.UpdatedAt),
	}
}