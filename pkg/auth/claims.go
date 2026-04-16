package auth

import	"github.com/golang-jwt/jwt/v5"

// UserClaims defines the structure of the JWT payload.
// It supports both an "Identity-Only" state (post-registration) and 
// a "Tenant-Context" state (after store selection/creation).
type UserClaims struct {
    Id       string  `json:"id"`
    Email    string  `json:"email"`
    // Role and TenantId are pointers to allow null/omitted values in the JSON 
    // payload when the user has not yet entered a specific tenant context.
    Role     *string `json:"role,omitempty"`
    TenantId *string `json:"tenantId,omitempty"`
    jwt.RegisteredClaims
}
