package middleware

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/MartinMurithi/storeforge/gateway/internal/response"
	"github.com/MartinMurithi/storeforge/pkg/auth"
	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(pubKey *rsa.PublicKey, expectedAudience, expectedIssuer string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const serviceName = "storeforge-api"

		// Extract Authorization Header
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "AUTH_REQUIRED", "Authorization header is missing")
			c.Abort()
			return
		}

		// Check for Bearer prefix
		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			c.Header("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s", error="invalid_request"`, serviceName))
			response.Error(c, http.StatusUnauthorized, "INVALID_FORMAT", "Must use Bearer token format")
			c.Abort()
			return
		}

		if len(authHeader) <= 7 {
			response.Error(c, http.StatusUnauthorized, "INVALID_FORMAT", "Token missing")
			c.Abort()
			return
		}

		tokenStr := authHeader[7:] // Extract the token part

		// Verify using shared pkg/auth function
		claims, err := auth.VerifyToken(pubKey, tokenStr, expectedAudience, expectedIssuer)
		if err != nil {
			c.Header("WWW-Authenticate", fmt.Sprintf(`Bearer realm="%s", error="invalid_token"`, serviceName))

			errCode := "INVALID_TOKEN"
			if errors.Is(err, apperrors.ErrExpiredToken) {
				errCode = "TOKEN_EXPIRED"
			}

			response.Error(c, http.StatusUnauthorized, errCode, err.Error())
			c.Abort()
			return
		}

		c.Set(auth.CtxUserID, claims.Id)         //UserId
		c.Set(auth.CtxTenantID, claims.TenantId) //TenantId(alias StoreId)
		c.Set(auth.CtxRole, claims.Role)
		c.Set(auth.CtxEmail, claims.Email)

		// Contextual realm for logs/tracing
		c.Set("realm", fmt.Sprintf("%s/%s", serviceName, claims.TenantId))

		c.Next()
	}
}
