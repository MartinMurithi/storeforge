package dto

import (
	"time"

	"github.com/MartinMurithi/storeforge/gateway/internal/dto/shared"
)


type UserResponseDTO struct {
	ID         string         `json:"id"`
	Email      string         `json:"email"`
	IsVerified bool           `json:"is_verified"`
	Profile    UserProfileDTO `json:"profile"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  *time.Time     `json:"deleted_at,omitempty"`
}

type UserProfileDTO struct {
	FullName     string `json:"full_name"`
	Phone        string `json:"phone"`
}

type GetCurrentUserResponseDTO struct {
	User UserResponseDTO `json:"user"`
}

type GetAllUsersResponseDTO struct {
	Users []UserResponseDTO   `json:"users"`
	Meta  shared.PaginationMetaDTO `json:"meta"`
}

type UpdateUserRequestDTO struct {
	Email *string `json:"email,omitempty" binding:"required"`
	Phone *string `json:"phone,omitempty" binding:"required"`
}

type UpdateUserResponseDTO struct {
	User UserResponseDTO `json:"user"`
    Message string `json:"message"`
}

type DeleteUserResponseDTO struct {
	Message string `json:"message"`
}