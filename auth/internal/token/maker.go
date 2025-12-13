package token

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTMaker struct {
	PrivateKey *rsa.PrivateKey
}

func NewJWTMaker(privateKey *rsa.PrivateKey) *JWTMaker {
	return &JWTMaker{PrivateKey: privateKey}
}

func (maker *JWTMaker) CreateToken(id, tenantId uuid.UUID, email, role string, duration time.Duration) (string, *UserClaims, error) {
	const OP = "Token.CreateToken"
	claims, err := NewUserClaims(id, tenantId, email, role, duration)

	if err != nil {
		return "", nil, fmt.Errorf("%s: error creating user claims %w", OP, err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenStr, err := token.SignedString(maker.PrivateKey)

	if err != nil {
		return "", nil, fmt.Errorf("%s error signing token %w", OP, err)
	}

	return tokenStr, claims, nil
}
