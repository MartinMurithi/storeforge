package dto

import (
	"errors"
	"fmt"
	// "log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// func GetUserId(c *gin.Context) (pgtype.UUID, error) {

// 	var id pgtype.UUID

// 	val, exists := c.Get("userId")

// 	log.Printf("user id from parser %v", val)

// 	if !exists {
// 		return pgtype.UUID{}, errors.New("id not found")
// 	}

// 	// Use pgtype's internal Scan method to parse the string
// 	// This handles the conversion from "550e8400-e2..." to the byte array

// 	err := id.Scan(val)

// 	if err != nil {
// 		return pgtype.UUID{}, errors.New("invalid user id type")
// 	}

// 	return id, nil
// }

func GetUserId(c *gin.Context) (pgtype.UUID, error) {
    val, exists := c.Get("userId")

    if !exists {
        return pgtype.UUID{}, errors.New("id not found in context")
    }

    // 1. Cast the context value to a string
    strID, ok := val.(string)
	
    if !ok {
        // If it's already a pgtype.UUID (unlikely but possible), return it
        if uuid, ok := val.(pgtype.UUID); ok {
            return uuid, nil
        }
        return pgtype.UUID{}, errors.New("user id in context is not a string")
    }

    var id pgtype.UUID

    err := id.Scan(strID) 
	
    if err != nil {
        return pgtype.UUID{}, fmt.Errorf("failed to parse uuid string: %w", err)
    }

    if !id.Valid {
        return pgtype.UUID{}, errors.New("parsed uuid is marked as invalid")
    }

    return id, nil
}