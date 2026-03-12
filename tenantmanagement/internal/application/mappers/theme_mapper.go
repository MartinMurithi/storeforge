package mappers

import (
	"log"

	themev1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/theme/v1"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/entity"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoTheme maps the Theme domain entity to its Protobuf representation.
func ToProtoTheme(th *entity.Theme) *themev1.Theme {
	if th == nil {
		return nil
	}

	// Safely convert map[string]any to Protobuf Struct
	configStruct, err := structpb.NewStruct(th.DefaultConfig.Config)
	if err != nil {
		// If conversion fails, we return an empty struct to avoid a nil pointer
		log.Printf("[Mapper]: failed to convert ThemeConfig to structpb: %v", err)
		configStruct = &structpb.Struct{}
	}

	return &themev1.Theme{
		Id:          th.ID.String(),
		Name:        th.Name,
		Description: th.Description,
		ThemeConfig: configStruct,
		IsActive:    th.IsActive,
		CreatedAt:   timestamppb.New(th.CreatedAt),
	}
}
