package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserResponseDTO struct {
	Id           uuid.UUID  `json:"id"`
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
	Id uuid.UUID `json:"id"`
}

type LoginUserRequestDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserResponseDTO struct {
	Token string `json:"string"`
	User  *UserResponseDTO
}
