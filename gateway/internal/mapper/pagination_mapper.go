package mapper

import (
	userv1 "github.com/MartinMurithi/storeforge/api/protos/user/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto/shared"
)

func ProtoPaginationToDTO(meta *userv1.PaginationMeta) shared.PaginationMetaDTO {
	if meta == nil {
		return shared.PaginationMetaDTO{}
	}

	return shared.PaginationMetaDTO{
		Page:       int(meta.Page),
		Limit:      int(meta.Limit),
		Total:      int(meta.Total),
		TotalPages: int(meta.TotalPages),
		HasNext:    meta.HasNext,
		HasPrev:    meta.HasPrev,
	}
}