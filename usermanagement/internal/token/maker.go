package token

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/MartinMurithi/storeforge/pkg/auth"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/domain/entity"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type JWTMaker struct {
	PrivateKey *rsa.PrivateKey
}

func NewJWTMaker(privateKey *rsa.PrivateKey) (*JWTMaker, error) {

	if privateKey == nil {
		return nil, errors.New("jwt private key is nil")
	}

	return &JWTMaker{PrivateKey: privateKey}, nil
}

func NewUserClaims(
    id pgtype.UUID,
    tenantId pgtype.UUID,
    email, role string,
    duration time.Duration,
) (*auth.UserClaims, error) {

    now := time.Now().UTC()

    return &auth.UserClaims{
        Id:       id.String(),
        Role:     role,
        Email:    email,
        TenantId: tenantId.String(),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
            IssuedAt:  jwt.NewNumericDate(now),
            ID:        uuid.NewString(),
            Subject:   id.String(),
            Issuer:    "usermanagement.auth.storeforge",
            Audience:  []string{"storeforge-api"},
        },
    }, nil
}

func (maker *JWTMaker) CreateToken(id, tenantId pgtype.UUID, email, role string, duration time.Duration) (*entity.Token, *auth.UserClaims, error) {
	const OP = "Token.CreateToken"

	if maker == nil || maker.PrivateKey == nil {
		return nil, nil, fmt.Errorf("%s: jwt maker not initialized", OP)
	}

	claims, err := NewUserClaims(id, tenantId, email, role, duration)

	if err != nil {
		return nil, nil, fmt.Errorf("%s: error creating user claims %w", OP, err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenStr, err := token.SignedString(maker.PrivateKey)

	if err != nil {
		return nil, nil, fmt.Errorf("%s error signing token %w", OP, err)
	}

	return &entity.Token{
		AccessToken: tokenStr,
		IssuedAt:    claims.IssuedAt.Local(),
		ExpiresAt:   claims.ExpiresAt.Local(),
		ExpiresIn:   claims.ExpiresAt.Unix() - claims.IssuedAt.Unix(),
		TokenType:   "Bearer",
	}, claims, nil
}
