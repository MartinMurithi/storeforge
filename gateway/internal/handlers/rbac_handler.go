package handlers

import (
	"log"
	"net/http"

	rbacv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/rbac/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
	"github.com/MartinMurithi/storeforge/gateway/internal/response"
	"github.com/gin-gonic/gin"
)

type RbacHandler struct {
	RbacClient rbacv1.RbacServiceClient
}

func (h *RbacHandler) CreateRole(c *gin.Context) {
	if h.RbacClient == nil {
		log.Println("Internal Error: RBAC Client not initialized in RBAC Handler")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "RBAC service unavailable")
		return
	}

	var reqDTO dto.RegisterRequestDTO

	// 1. Validate JSON Input
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		response.Error(c, http.StatusBadRequest, "MALFORMED_JSON", "kindly check your request body")
		return
	}

	grpcRequest := &authv1.RegisterRequest{
		FullName: reqDTO.FullName,
		Email:    reqDTO.Email,
		Phone:    reqDTO.Phone,
		Password: reqDTO.Password,
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
