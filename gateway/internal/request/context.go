package request

import (
	"errors"
	"fmt"

	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

// GetUserId retrieves the logged-in user's ID as a string.
// The Gateway treats IDs as strings; only the internal services
// will handle database-specific types like pgtype.
func GetUserId(c *gin.Context) (string, error) {
	val, exists := c.Get("userId")

	if !exists {
		return "", errors.New("user ID not found in context")
	}

	id, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("user ID in context is not a string (got: %T)", val)
	}

	if id == "" {
		return "", errors.New("user ID is empty")
	}

	return id, nil
}

// GetUserParamId parses the "id" URL parameter and returns it as a string.
// It performs a basic structural check to ensure the ID is a valid UUID format.
func GetUserParamId(c *gin.Context) (string, error) {
	strID := c.Param("id")

	if strID == "" {
		return "", errors.New("id parameter is required")
	}

	if _, err := uuid.Parse(strID); err != nil {
		return "", fmt.Errorf("invalid UUID format: %w", err)
	}

	return strID, nil
}

func GetUserRole(c *gin.Context) (string, error) {
	val, exists := c.Get("role")
	if !exists {
		return "", apperrors.ErrRoleIsRequired
	}

	role, ok := val.(string)
	if !ok {
		return "", errors.New("user role in context is not a string")
	}

	if role == "" {
		return "", errors.New("user role is empty")
	}

	return role, nil
}
