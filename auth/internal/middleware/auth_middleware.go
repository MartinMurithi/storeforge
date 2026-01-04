package middleware

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/MartinMurithi/storeforge/auth/internal/token"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtMaker *token.JWTMaker, pubKey *rsa.PublicKey, expectedAudience, expectedIssuer string) gin.HandlerFunc {

	return func(c *gin.Context) {

		serviceName := "storeforge-api"

		// Get Authorization Header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {

			c.Header("WWW-Authenticate", fmt.Sprintf(
				`Bearer realm="%s", error="invalid_request", error_description="Invalid Authorization Header Format"`,
				serviceName,
			))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		claims, err := token.VerifyToken(pubKey, tokenStr, expectedAudience, expectedIssuer)

		if err != nil {
			c.Header("WWW-Authenticate", fmt.Sprintf(
				`Bearer realm="%s", error="invalid_token", error_description="The access token is invalid or expired"`,
				serviceName,
			))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Token is valid
		c.Set("userID", claims.ID)
		c.Set("tenantID", claims.TenantId)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)
		c.Set("realm", fmt.Sprintf("%s/%s", serviceName, claims.TenantId.String()))

		c.Next()

	}
}
