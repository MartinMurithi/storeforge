package dto

import "time"

type RegisterRequestDTO struct {
	FullName     string `json:"full_name" binding:"required"`
	Email        string `json:"email" binding:"required"`
	Phone        string `json:"phone" binding:"required"`
	Password     string `json:"password" binding:"required"`
	BusinessName string `json:"business_name" binding:"required"`
	BusinessType string `json:"business_type" binding:"required"`
}

type RegisterResponseDTO struct {
	User    UserResponseDTO `json:"user"`
	Message string          `json:"message"`
}

type LoginRequestDTO struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseDTO struct {
	User  UserResponseDTO `json:"user"`
	Token TokenDTO        `json:"token"`
}

type RefreshTokenRequestDTO struct{
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponseDTO struct {
	 Token TokenDTO `json:"token"`
}

type TokenDTO struct {
	AccessToken  string    `json:"access_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	ExpiresIn    int64     `json:"expires_in"`
	IssuedAt     time.Time `json:"issued_at"`
	TokenType    string    `json:"token_type"`
}

type LogoutResponseDTO struct{
	success bool
}