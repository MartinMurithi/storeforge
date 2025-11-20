package models

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID uuid.UUID `json:"userId"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
