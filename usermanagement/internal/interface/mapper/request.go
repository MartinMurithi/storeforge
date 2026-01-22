package mapper

import (
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPatchUserRequest(id pgtype.UUID, req *dto.PatchUserRequestDTO) *application.PatchUserInput {
	return &application.PatchUserInput{
		Id:           id,
		BusinessName: req.BusinessName,
		BusinessType: req.BusinessType,
	}

}
