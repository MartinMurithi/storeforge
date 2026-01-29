package entity

import (
	"time"

)

type Token struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	ExpiresIn   int64     `json:"expires_in"` // seconds
	IssuedAt    time.Time `json:"issued_at"`
	TokenType   string    `json:"token_type"` // "Bearer"
}
