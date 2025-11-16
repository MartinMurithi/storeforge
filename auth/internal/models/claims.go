package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserID string `json:"userId"`
	Roles []Role `json:"roles"`
	jwt.StandardClaims
}
