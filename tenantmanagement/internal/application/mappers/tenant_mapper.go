package mappers

import (
	"log"

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
		BusinessType: t.BusinessType,
		CreatedAt:    timestamppb.New(t.CreatedAt),
	}
}


// ToProtoCreateTenantResponse assembles the final gRPC response using the DTO.
// It handles the conversion of map[string]any to google.protobuf.Struct.
func ToProtoCreateTenantResponse(data *dtos.CreateTenantResponseDTO) *tenantv1.CreateTenantResponse {
	if data == nil || data.Tenant == nil {
		return nil
	}

	// Safely convert map[string]any to Protobuf Struct
	configStruct, err := structpb.NewStruct(data.Tenant.Settings.Config)
	if err != nil {
		// If conversion fails, we return an empty struct to avoid a nil pointer 
		log.Printf("[Mapper]: failed to convert ThemeConfig to structpb: %v", err)
		configStruct = &structpb.Struct{}
	}

	return &tenantv1.CreateTenantResponse{
		Tenant: ToProtoTenant(data.Tenant),
		Settings: &tenantv1.TenantSettings{
			ThemeId:     data.Tenant.Settings.ThemeID.String(),
			ThemeConfig: configStruct,
			Version:     int32(data.Tenant.Settings.Version),
		},
	}
}