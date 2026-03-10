package mapper

import (
	userv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/user/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto/shared"
)

// MapPaginationMetaProtoToDTO converts gRPC metadata into a shared PaginationMetaDTO.
func MapPaginationMetaProtoToDTO(pb *userv1.PaginationMeta) *shared.PaginationMetaDTO {
	if pb == nil {
		return &shared.PaginationMetaDTO{}
	}

	return &shared.PaginationMetaDTO{
		Page:       int32(pb.Page),
		Limit:      int32(pb.Limit),
		Total:      int64(pb.Total),
		TotalPages: int32(pb.TotalPages),
		HasNext:    pb.HasNext,
		HasPrev:    pb.HasPrev,
	}
}
