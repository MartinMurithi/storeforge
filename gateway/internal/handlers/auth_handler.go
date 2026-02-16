package handlers

import (
	"net/http"

	authv1 "github.com/MartinMurithi/storeforge/api/protos/auth/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/client"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
	"github.com/MartinMurithi/storeforge/gateway/internal/mapper"
	"github.com/MartinMurithi/storeforge/gateway/internal/response"
	"github.com/MartinMurithi/storeforge/pkg/errconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthClient *client.AuthClient
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
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
