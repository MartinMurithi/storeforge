package handlers

import (
	"log"
	"net/http"

	authv1 "github.com/MartinMurithi/storeforge/api/protos/auth/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
	"github.com/MartinMurithi/storeforge/gateway/internal/mapper"
	"github.com/MartinMurithi/storeforge/gateway/internal/response"
	"github.com/MartinMurithi/storeforge/pkg/errconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthClient authv1.AuthServiceClient
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	if h.AuthClient == nil {
		log.Println("Internal Error: AuthClient not initialized in AuthHandler")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Auth service unavailable")
		return
	}

	var reqDTO dto.RegisterRequestDTO

	// 1. Validate JSON Input
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		response.Error(c, http.StatusBadRequest, "MALFORMED_JSON", "kindly check your request body")
		return
	}

	grpcRequest := &authv1.RegisterRequest{
		FullName:     reqDTO.FullName,
		Email:        reqDTO.Email,
		Phone:        reqDTO.Phone,
		Password:     reqDTO.Password,
		BusinessName: reqDTO.BusinessName,
		BusinessType: reqDTO.BusinessType,
	}

	res, err := h.AuthClient.Register(c.Request.Context(), grpcRequest)

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	resp := mapper.MapRegisterResponseProtoToDTO(res)
	response.JSON(c, http.StatusCreated, resp)
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	var reqDTO dto.LoginRequestDTO

	// 1. Validate JSON Input
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		response.Error(c, http.StatusBadRequest, "MALFORMED_JSON", "kindly check your request body")
		return
	}

	grpcRequest := &authv1.LoginRequest{
		Email:    reqDTO.Email,
		Password: reqDTO.Password,
	}

	res, err := h.AuthClient.Login(c.Request.Context(), grpcRequest)

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	// Set the NEW Refresh Token into the Cookie
	c.SetCookie(
		"refresh_token",
		res.Token.RefreshToken,
		3600*24*7, // 7 days
		"/api/v1/auth",
		"",    // domain (empty for current)
		false, // secure (set to true for HTTPS)
		true,  // httpOnly (Critical: JS can't see this)
	)

	safeRes := mapper.MapLoginResponseProtoToDTO(res)

	response.JSON(c, http.StatusOK, safeRes)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Extract refresh token from Cookie
	refreshToken, err := c.Cookie("refresh_token")

	if err != nil {
		// If cookie is missing, the user is not logged in/session expired
		response.Error(c, http.StatusUnauthorized, "MISSING_REFRESH_TOKEN", "session expired or invalid")
		return
	}

	grpcRequest := &authv1.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	res, err := h.AuthClient.RefreshToken(c.Request.Context(), grpcRequest)

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	// Set the NEW rotated Refresh Token back into the Cookie
	c.SetCookie(
		"refresh_token",
		res.Token.RefreshToken,
		3600*24*7, // 7 days
		"/api/v1/auth",
		"",    // domain (empty for current)
		false, // secure (set to true for HTTPS)
		true,  // httpOnly (Critical: JS can't see this)
	)

	safeRes := mapper.MapRefreshTokenResponseProtoToSafeDTO(res)

	response.JSON(c, http.StatusOK, safeRes)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if err == nil && refreshToken != "" {
		grpcRequest := &authv1.LogoutRequest{
			RefreshToken: refreshToken,
		}

		_, _ = h.AuthClient.Logout(c.Request.Context(), grpcRequest)
	}

	// Clear the Cookie from the browser
	c.SetCookie(
		"refresh_token",
		"",
		-1, // MaxAge -1 tells the browser to delete the cookie immediately
		"/api/v1/auth",
		"",    // Domain
		false, // Secure (Match your dev/prod setting)
		true,  // HttpOnly
	)

	response.JSON(c, http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
