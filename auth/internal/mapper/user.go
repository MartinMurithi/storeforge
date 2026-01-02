package mapper

import (
	"github.com/MartinMurithi/storeforge/auth/internal/dto"
	"github.com/MartinMurithi/storeforge/auth/internal/models"
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
		Status: "success",
		Message: "User registered successfully",
		Data: &dto.UserResponseDTO{
			Id:           user.ID,
			FullName:     user.FullName,
			Email:        user.Email,
			Phone:        user.Phone,
			BusinessType: user.BusinessType,
			BusinessName: user.BusinessName,
			CreatedAt:    user.CreatedAt,
			IsVerified:   user.IsVerified,
		},
	}
	
}

func ToLoginUserResponse(token string, user *models.User) *dto.LoginUserResponseDTO {
	if user == nil {
		return nil
	}
	return &dto.LoginUserResponseDTO{
		Token: token,
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
