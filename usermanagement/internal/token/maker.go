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

func NewUserClaims(id pgtype.UUID, tenantId pgtype.UUID, email, role string, duration time.Duration) (*auth.UserClaims, error) {

	googleUUID, err := uuid.NewRandom()

	if err != nil {
		return nil, fmt.Errorf("error generating google id %w", err)
	}

	// Convert the Google UUID to pgtypeUUID for compatibility
	pgUUID := pgtype.UUID{
		Bytes: googleUUID,
		Valid: true,
	}

	return &auth.UserClaims{
		Id:       id.String(), //Id of the store owner
		Role:     role,
		Email:    email,
		TenantId: tenantId.String(), //Id of the tenant(actual store)
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        pgUUID.String(),
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
