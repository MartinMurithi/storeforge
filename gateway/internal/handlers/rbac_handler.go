package handlers

import (
	"log"
	"net/http"

	rbacv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/rbac/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
	"github.com/MartinMurithi/storeforge/gateway/internal/mapper"
	"github.com/MartinMurithi/storeforge/gateway/internal/response"
	"github.com/MartinMurithi/storeforge/gateway/internal/util"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
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

	var reqDTO dto.CreateRoleRequestDTO

	if !util.BindAndValidateJSON(c, &reqDTO) {
		return
	}

	grpcRequest := &rbacv1.CreateRoleRequest{
		Name:          reqDTO.Name,
		Slug:          reqDTO.Slug,
		Description:   reqDTO.Description,
		PermissionIds: reqDTO.PermissionIDs,
	}

	res, err := h.RbacClient.CreateRole(c.Request.Context(), grpcRequest)

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	resp := mapper.MapCreateRoleResponseProtoToDTO(res)
	response.JSON(c, http.StatusCreated, resp)
}
