package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidationError(c *gin.Context) {
	Error(c, http.StatusBadRequest, "invalid request")
}
