package util

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValidateJSON(c *gin.Context, dst any) bool {
	if err := c.ShouldBindJSON(dst); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fields := make(map[string]string)

			for _, fe := range ve {
				field := toSnakeCase(fe.Field())

				switch fe.Tag() {
				case "required":
					fields[field] = "is required"
				case "uuid4":
					fields[field] = "must be a valid UUID"
				case "min":
					fields[field] = fmt.Sprintf("must be at least %s characters/items", fe.Param())
				case "max":
					fields[field] = fmt.Sprintf("must not exceed %s characters/items", fe.Param())
				case "oneof":
					fields[field] = fmt.Sprintf("must be one of: %s", fe.Param())
				default:
					fields[field] = "is invalid"
				}
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error": gin.H{
					"code":    "VALIDATION_ERROR",
					"message": "one or more fields are invalid",
					"fields":  fields,
				},
			})
			return false
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "BAD_REQUEST",
				"message": "invalid request body",
			},
		})
		return false
	}

	return true
}

func toSnakeCase(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}