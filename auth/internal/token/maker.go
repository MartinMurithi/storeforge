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

// Token response
type Token struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	ExpiresIn   int64     `json:"expires_in"` // seconds
	IssuedAt    time.Time `json:"issued_at"`
	TokenType   string    `json:"token_type"` // "Bearer"
}

func NewJWTMaker(privateKey *rsa.PrivateKey) (*JWTMaker, error) {

	if privateKey == nil {
		return nil, errors.New("jwt private key is nil")
	}

	return &JWTMaker{PrivateKey: privateKey}, nil
}

func (maker *JWTMaker) CreateToken(id, tenantId uuid.UUID, email, role string, duration time.Duration) (*Token, *UserClaims, error) {
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

	return &Token{
		AccessToken: tokenStr,
		IssuedAt:    claims.IssuedAt.Local(),
		ExpiresAt:   claims.ExpiresAt.Local(),
		ExpiresIn:   claims.ExpiresAt.Unix() - claims.IssuedAt.Unix(),
		TokenType:   "Bearer",
	}, claims, nil
}
