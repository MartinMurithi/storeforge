package mapper

import (
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/application/user"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPatchUserRequest(id pgtype.UUID, req *dto.PatchUserRequestDTO) *user.PatchUserInput {
	return &user.PatchUserInput{
		Id:           id,
		BusinessName: req.BusinessName,
		BusinessType: req.BusinessType,
	}

}
