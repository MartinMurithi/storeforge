package user

import "github.com/MartinMurithi/storeforge/gateway/internal/dto/shared"

type UserResponseDTO struct {
    Id           string  `json:"id"`
    FullName     string  `json:"fullName"`
    Email        string  `json:"email"`
    Phone        string  `json:"phone"`
    BusinessType string  `json:"businessType"`
    BusinessName string  `json:"businessName"`
    CreatedAt    string  `json:"createdAt"`
    UpdatedAt    *string `json:"updatedAt,omitempty"`
    IsVerified   bool    `json:"isVerified"`
}

type FetchUserResponseDTO struct {
    User *UserResponseDTO `json:"user"`
}

type FetchAllUsersResponseDTO struct {
    Users      []UserResponseDTO `json:"users"`
    Pagination shared.PaginationMeta `json:"pagination"`
}