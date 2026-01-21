package mapper

import (
	"github.com/MartinMurithi/storeforge/auth/internal/dto"
	"github.com/MartinMurithi/storeforge/auth/internal/services"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPatchUserRequest(id pgtype.UUID, req *dto.PatchUserRequestDTO) *services.PatchUserInput {
	return &services.PatchUserInput{
		Id:           id,
		BusinessName: req.BusinessName,
		BusinessType: req.BusinessType,
	}

}
