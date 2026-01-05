package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserClaims struct {
	Id       pgtype.UUID `json:"id"`
	Role     string      `json:"role"`
	Email    string      `json:"email"`
	TenantId pgtype.UUID `json:"tenantId"`
	jwt.RegisteredClaims
}

func NewUserClaims(id pgtype.UUID, tenantId pgtype.UUID, email, role string, duration time.Duration) (*UserClaims, error) {

	googleUUID, err := uuid.NewRandom()

	if err != nil {
		return nil, fmt.Errorf("error generating google id %w", err)
	}

	// Convert the Google UUID to pgtypeUUID for compatibility
	pgUUID := pgtype.UUID{
		Bytes: googleUUID,
		Valid: true,
	}

	return &UserClaims{
		Id:       id,
		Role:     role,
		Email:    email,
		TenantId: tenantId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        pgUUID.String(),
			Subject:   id.String(),
			Issuer:    "auth.storeforge",
			Audience:  []string{"storeforge-api"},
		},
	}, nil
}
