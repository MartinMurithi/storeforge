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
		Email: reqDTO.Email,
		Password: reqDTO.Password,
	}

	res, err := h.AuthClient.Login(c.Request.Context(), grpcRequest)

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	resp := mapper.MapLoginResponseProtoToDTO(res)
	response.JSON(c, http.StatusOK, resp)
}
