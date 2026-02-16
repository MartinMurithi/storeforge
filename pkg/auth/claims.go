package auth

import	"github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	Id       string `json:"id"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	TenantId string `json:"tenantId"`
	jwt.RegisteredClaims
}
