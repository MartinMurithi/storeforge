package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey}
}

func (maker *JWTMaker) CreateToken(id uuid.UUID, email, role string, duration time.Duration) (string, *UserClaims, error) {
	const OP = "Token.CreateToken"
	claims, err := NewUserClaims(id, email, role, duration)

	if err != nil {
		return "", nil, fmt.Errorf("%s: error creating user claims %w", OP, err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenStr, err := token.SignedString([]byte(maker.secretKey))

	if err != nil {
		return "", nil, fmt.Errorf("%s error signing token %w", OP, err)
	}

	return tokenStr, claims, nil
}

// func VerifyToken(token string) (*UserClaims, error){
// 	jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		//verify signing method
// 		_, ok := token.Method.Verify(jwt.SigningMethodES256.Name)
// 	})
// }