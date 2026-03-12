package dto

import (
	"strings"
	"time"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserResponseDTO struct {
	Id         pgtype.UUID `json:"id"`
	FullName   string      `json:"fullName"`
	Email      string      `json:"email"`
	Phone      string      `json:"phone"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  *time.Time  `json:"updatedAt"`
	IsVerified bool        `json:"isVerified"`
}

type RegisterUserRequestDTO struct {
	FullName string
	Email    string
	Phone    string
	Password string
}

type RegisterUserResponseDTO struct {
	User    *UserResponseDTO `json:"user"`
	Message string           `json:"message"`
}

// Normalize Registration user input
// Email and phone are also normalized in the validators
func (regInput *RegisterUserRequestDTO) Normalize() {
	regInput.FullName = strings.TrimSpace(regInput.FullName)
	regInput.Email = strings.TrimSpace(regInput.Email)
	regInput.Phone = strings.TrimSpace(regInput.Phone)
	regInput.Password = strings.TrimSpace(regInput.Password)
}

// Normalize Login user input
// Email and phone are also normalized in the validators
func (regInput *LoginUserRequestDTO) Normalize() {
	regInput.Email = strings.TrimSpace(regInput.Email)
	regInput.Password = strings.TrimSpace(regInput.Password)
}

type LoginUserRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserResponseDTO struct {
	User  *UserResponseDTO `json:"user"`
	Token *entity.Token    `json:"token"`
}

type FetchAllUsersResponseDTO struct {
	Users      []UserResponseDTO `json:"users"`
	Pagination PaginationMeta    `json:"pagination"`
}

type FetchUserResponseDTO struct {
	User *UserResponseDTO `json:"user"`
}

type PatchUserRequestDTO struct {
	Email *string `json:"email"`
	Phone *string `json:"phone"`
}

type UpdateActiveSessionContextRequestDTO struct {
	UserId   pgtype.UUID `json:"userId"`
	TenantId pgtype.UUID `json:"tenantId"`
	Role     string      `json:"role"`
}
