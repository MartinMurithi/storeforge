package mapper

import (
	userv1 "github.com/MartinMurithi/storeforge/api/protos/user/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto/shared"
)

// MapPaginationMetaProtoToDTO converts gRPC metadata into a shared PaginationMetaDTO.
func MapPaginationMetaProtoToDTO(pb *userv1.PaginationMeta) *shared.PaginationMetaDTO {
	if pb == nil {
		return &shared.PaginationMetaDTO{}
	}

	return &shared.PaginationMetaDTO{
		Page:       int(pb.Page),
		Limit:      int(pb.Limit),
		Total:      int(pb.Total),
		TotalPages: int(pb.TotalPages),
		HasNext:    pb.HasNext,
		HasPrev:    pb.HasPrev,
	}
}
