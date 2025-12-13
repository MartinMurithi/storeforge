package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	Id    uuid.UUID `json:"id"`
	Role  string    `json:"role"`
	Email string    `json:"email"`
	jwt.RegisteredClaims
}

func NewUserClaims(id uuid.UUID, email, role string, duration time.Duration) (*UserClaims, error) {
	tokenId, err := uuid.NewRandom()

	if err != nil {
		return nil, fmt.Errorf("error generating token id %w", err)
	}

	return &UserClaims{
		Id:    id,
		Role:  role,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        tokenId.String(),
			Subject:   email,
			Issuer:    "storeforge",
			Audience:  []string{"storeforge-api"},
		},
	}, nil
}
