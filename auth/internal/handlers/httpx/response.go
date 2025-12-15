package httpx

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Data  any `json:"data,omitempty"`
	Error any `json:"error,omitempty"`
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, APIResponse{Data: data})
}

func Error(c *gin.Context, status int, msg string) {
	c.JSON(status, APIResponse{Error: msg})
}
