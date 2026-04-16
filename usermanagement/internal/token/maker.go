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

// NewUserClaims constructs a new UserClaims instance.
// If tenantId.Valid is false or role is empty, the resulting claims represent 
// a partial identity, restricting the user to onboarding or global actions.
func NewUserClaims(
    id pgtype.UUID,
    tenantId pgtype.UUID,
    email, role string,
    duration time.Duration,
) (*auth.UserClaims, error) {
    now := time.Now().UTC()

    claims := &auth.UserClaims{
        Id:    id.String(),
        Email: email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
            IssuedAt:  jwt.NewNumericDate(now),
            ID:        uuid.NewString(),
            Subject:   id.String(),
            Issuer:    "usermanagement.auth.storeforge",
            Audience:  []string{"storeforge-api"},
        },
    }

    // Populate tenant context only if valid. 
    // This distinction is used by the gateway to enforce PBAC boundaries.
    if tenantId.Valid {
        tID := tenantId.String()
        claims.TenantId = &tID
    }

    if role != "" {
        claims.Role = &role
    }

    return claims, nil
}


// CreateToken generates a signed JWT. It handles the transition from 
// basic authentication to tenant-specific authorization by checking the 
// validity of the provided tenantId.
func (maker *JWTMaker) CreateToken(
    id, tenantId pgtype.UUID, 
    email, role string, 
    duration time.Duration,
) (*entity.Token, *auth.UserClaims, error) {
    const OP = "Token.CreateToken"

    if maker == nil || maker.PrivateKey == nil {
        return nil, nil, fmt.Errorf("%s: jwt maker not initialized", OP)
    }

    claims, err := NewUserClaims(id, tenantId, email, role, duration)
   
	if err != nil {
        return nil, nil, fmt.Errorf("%s: error creating user claims %w", OP, err)
    }

    // Use RS256 for asymmetric signing; requires a valid RSA Private Key.
    token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
    tokenStr, err := token.SignedString(maker.PrivateKey)

    if err != nil {
        return nil, nil, fmt.Errorf("%s: error signing token %w", OP, err)
    }

    return &entity.Token{
        AccessToken: tokenStr,
        IssuedAt:    claims.IssuedAt.Time,
        ExpiresAt:   claims.ExpiresAt.Time,
        ExpiresIn:   int64(duration.Seconds()),
        TokenType:   "Bearer",
    }, claims, nil
}