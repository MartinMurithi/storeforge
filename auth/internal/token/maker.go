package token

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (maker *JWTMaker) CreateToken(id, tenantId uuid.UUID, email, role string, duration time.Duration) (string, *UserClaims, error) {
	const OP = "Token.CreateToken"

	if maker == nil || maker.PrivateKey == nil {
		return "", nil, fmt.Errorf("%s: jwt maker not initialized", OP)
	}

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
