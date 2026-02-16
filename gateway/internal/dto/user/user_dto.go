package user

import (
	"time"

	"github.com/MartinMurithi/storeforge/gateway/internal/dto/shared"
)


type UserResponse struct {
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
	BusinessName string `json:"business_name"`
	BusinessType string `json:"business_type"`
}

type GetCurrentUserResponse struct {
	User UserResponse `json:"user"`
}

type GetAllUsersResponse struct {
	Users []UserResponse   `json:"users"`
	Meta  shared.PaginationMetaDTO `json:"meta"`
}

type UpdateUserRequest struct {
	BusinessName *string `json:"business_name,omitempty"`
	BusinessType *string `json:"business_type,omitempty"`
}

type UpdateUserResponse struct {
	User UserResponse `json:"user"`
    Message string `json:"message"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}