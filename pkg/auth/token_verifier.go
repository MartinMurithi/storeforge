package auth

import (
	"crypto/rsa"
	"fmt"

	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyToken(pubKey *rsa.PublicKey, tokenStr, expectedAudience, expectedIssuer string) (*UserClaims, error) {
	// Parse the token with custom claims
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		//Enforce that the signature method is RSA
		_, ok := t.Method.(*jwt.SigningMethodRSA)

		if !ok {
			return nil, fmt.Errorf("an error occurreed: %w", apperrors.ErrInvalidToken)
		}
		return pubKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error occurred when parsing token %w", err)
	}

	// Extract claims from token
	claims, ok := token.Claims.(*UserClaims)

	if !ok || !token.Valid {
		return nil, apperrors.ErrInvalidToken
	}

	return claims, nil
}
