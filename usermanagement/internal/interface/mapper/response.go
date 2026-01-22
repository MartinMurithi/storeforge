package mapper

import (
	"github.com/MartinMurithi/storeforge/usermanagement/internal/interface/dto"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/token"
)

func ToUserResponse(user *entity.User) *dto.UserResponseDTO {
	if user == nil {
		return nil
	}
	return &dto.UserResponseDTO{
		Id:           user.ID,
		FullName:     user.FullName,
		Email:        user.Email,
		Phone:        user.Phone,
		BusinessType: user.BusinessType,
		BusinessName: user.BusinessName,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		IsVerified:   user.IsVerified,
	}
}

func ToRegisterUserResponse(user *entity.User) *dto.RegisterUserResponseDTO {
	if user == nil {
		return nil
	}
	return &dto.RegisterUserResponseDTO{
		User: &dto.UserResponseDTO{
			Id:           user.ID,
			FullName:     user.FullName,
			Email:        user.Email,
			Phone:        user.Phone,
			BusinessType: user.BusinessType,
			BusinessName: user.BusinessName,
			CreatedAt:    user.CreatedAt,
			IsVerified:   user.IsVerified,
		},
		Message: "Registration successful. Please verify your email.",
	}

}

func ToLoginUserResponse(token *token.Token, user *entity.User) *dto.LoginUserResponseDTO {
	if user == nil {
		return nil
	}
	return &dto.LoginUserResponseDTO{
		User: &dto.UserResponseDTO{
			Id:           user.ID,
			FullName:     user.FullName,
			Email:        user.Email,
			Phone:        user.Phone,
			BusinessType: user.BusinessType,
			BusinessName: user.BusinessName,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			IsVerified:   user.IsVerified,
		},
		Token: token,
	}
}

func ToFetchAllUsersResponse(users []*entity.User, meta dto.PaginationMeta) *dto.FetchAllUsersResponseDTO {

	usersDTO := make([]dto.UserResponseDTO, 0, len(users))

	for _, u := range users {
		usersDTO = append(usersDTO, dto.UserResponseDTO{
			Id:           u.ID,
			FullName:     u.FullName,
			Email:        u.Email,
			Phone:        u.Phone,
			BusinessType: u.BusinessType,
			BusinessName: u.BusinessName,
			CreatedAt:    u.CreatedAt,
			UpdatedAt:    u.UpdatedAt,
			IsVerified:   u.IsVerified,
		})
	}

	return &dto.FetchAllUsersResponseDTO{
		Users:      usersDTO,
		Pagination: meta,
	}
}

func ToFetchUserResponse(user *entity.User) *dto.FetchUserResponseDTO {
	if user == nil {
		return nil
	}
	return &dto.FetchUserResponseDTO{
		User: &dto.UserResponseDTO{
			Id:           user.ID,
			FullName:     user.FullName,
			Email:        user.Email,
			Phone:        user.Phone,
			BusinessType: user.BusinessType,
			BusinessName: user.BusinessName,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			IsVerified:   user.IsVerified,
		},
	}
}
