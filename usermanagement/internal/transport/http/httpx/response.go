package httpx

import "github.com/gin-gonic/gin"

type APIError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type APIResponse struct {
	Data  any       `json:"data,omitempty"`
	Error *APIError `json:"error,omitempty"`
}

func JSON(c *gin.Context, status int, data any) {
	c.Header("Content-Type", "application/json")
	c.JSON(status, APIResponse{Data: data})
}

func Error(c *gin.Context, status int, code, msg string) {
	c.JSON(status, APIResponse{
		Error: &APIError{
			Code:    code,
			Message: msg,
		},
	})
}
