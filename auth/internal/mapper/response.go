package mapper

import (
	"github.com/MartinMurithi/storeforge/auth/internal/dto"
	"github.com/MartinMurithi/storeforge/auth/internal/models"
	"github.com/MartinMurithi/storeforge/auth/internal/token"
)

func ToUserResponse(user *models.User) *dto.UserResponseDTO {
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

func ToRegisterUserResponse(user *models.User) *dto.RegisterUserResponseDTO {
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

func ToLoginUserResponse(token *token.Token, user *models.User) *dto.LoginUserResponseDTO {
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

func ToFetchAllUsersResponse(users []*models.User, meta dto.PaginationMeta) *dto.FetchAllUsersResponseDTO {

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

func ToFetchUserResponse(user *models.User) *dto.FetchUserResponseDTO {
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
