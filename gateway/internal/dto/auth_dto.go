package dto

import "time"

type RegisterRequestDTO struct {
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	BusinessName string `json:"business_name"`
	BusinessType string `json:"business_type"`
}

type RegisterResponseDTO struct {
	User    UserResponseDTO `json:"user"`
	Message string       `json:"message"`
}

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	User  UserResponseDTO `json:"user"`
	Token TokenDTO     `json:"token"`
}

type TokenDTO struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	ExpiresIn   int64     `json:"expires_in"`
	IssuedAt    time.Time `json:"issued_at"`
	TokenType   string    `json:"token_type"`
}