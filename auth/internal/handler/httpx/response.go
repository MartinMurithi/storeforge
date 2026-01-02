package httpx

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Data  any `json:"data,omitempty"`
	Error any `json:"error,omitempty"`
}

func JSON(c *gin.Context, status int, data any) {
	c.Header("Content-Type", "application/json")
	c.JSON(status, data)
}

func Error(c *gin.Context, status int, msg string) {
	c.Header("Content-Type", "application/json")
	c.JSON(status, APIResponse{Error: msg})
}
