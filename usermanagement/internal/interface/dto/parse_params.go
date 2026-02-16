package dto

import (
	"errors"
	"fmt"

	"github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// GetUserID retrieves the logged-in user's ID from the Gin context
func GetUserId(c *gin.Context) (pgtype.UUID, error) {
    val, exists := c.Get("userId")
    if !exists {
        return pgtype.UUID{}, errors.New("user ID not found in context")
    }

    switch v := val.(type) {
    case string:
        var id pgtype.UUID
        if err := id.Scan(v); err != nil {
            return pgtype.UUID{}, fmt.Errorf("failed to parse UUID string: %w", err)
        }
        if !id.Valid {
            return pgtype.UUID{}, errors.New("parsed UUID is invalid")
        }
        return id, nil

    case pgtype.UUID:
        if !v.Valid {
            return pgtype.UUID{}, errors.New("UUID in context is invalid")
        }
        return v, nil

    default:
        return pgtype.UUID{}, errors.New("user ID in context has invalid type")
    }
}


// GetUserParamID parses the "id" URL parameter into a pgtype.UUID
func GetUserParamId(c *gin.Context) (pgtype.UUID, error) {
    strID := c.Param("id")

    if strID == "" {
        return pgtype.UUID{}, errors.New("id parameter is required")
    }

    var id pgtype.UUID

    if err := id.Scan(strID); err != nil {
        return pgtype.UUID{}, fmt.Errorf("failed to parse uuid string: %w", err)
    }

    if !id.Valid {
        return pgtype.UUID{}, errors.New("parsed uuid is invalid")
    }

    return id, nil
}


func GetUserRole(c *gin.Context) (string, error) {
	val, exists := c.Get("role")

	if !exists {
		return "", apperrors.ErrRoleIsRequired
	}

	role, ok := val.(string)

	if !ok {
		return "", errors.New("role has invalid type")
	}

	return role, nil
}
