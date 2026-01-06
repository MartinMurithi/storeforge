package dto

import (
	"strings"
	"time"

	"github.com/MartinMurithi/storeforge/auth/internal/token"

	"github.com/jackc/pgx/v5/pgtype"

)

type UserResponseDTO struct {
	Id           pgtype.UUID  `json:"id"`
	FullName     string     `json:"fullName"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone"`
	BusinessType string     `json:"businessType"`
	BusinessName string     `json:"businessName"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
	IsVerified   bool       `json:"isVerified"`
}

type RegisterUserRequestDTO struct {
	FullName     string `json:"fullName" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Phone        string `json:"phone" binding:"required"`
	Password     string `json:"password" binding:"required,min=8"`
	BusinessType string `json:"businessType" binding:"required"`
	BusinessName string `json:"businessName" binding:"required"`
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
	regInput.BusinessType = strings.TrimSpace(regInput.BusinessType)
	regInput.BusinessName = strings.TrimSpace(regInput.BusinessName)
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
	Token *token.Token     `json:"token"`
}

type FetchAllUsersResponseDTO struct{
	Users []UserResponseDTO `json:"users"`
	Pagination PaginationMeta `json:"pagination"`
}

type FetchUserResponseDTO struct{
	User *UserResponseDTO `json:"user"`
}