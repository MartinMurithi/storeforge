package handlers

import (
	"log"
	"net/http"

	rbacv1 "github.com/MartinMurithi/storeforge/api/protos/usermanagement/rbac/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
	"github.com/MartinMurithi/storeforge/gateway/internal/mapper"
	"github.com/MartinMurithi/storeforge/gateway/internal/request"
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

func (h *RbacHandler) GetRoleById(c *gin.Context) {
	if h.RbacClient == nil {
		log.Println("Internal Error: RBAC Client not initialized in RBAC Handler")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "RBAC service unavailable")
		return
	}

	roleID, err := request.GetParamId(c)

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	grpcRequest := &rbacv1.GetRoleByIDRequest{
		RoleId: roleID,
	}

	res, err := h.RbacClient.GetRoleByID(c.Request.Context(), grpcRequest)

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	resp := mapper.MapGetRoleByIDResponseProtoToDTO(res)
	response.JSON(c, http.StatusAccepted, resp)
}

func (h *RbacHandler) UpdateRole(c *gin.Context) {
	if h.RbacClient == nil {
		log.Println("Internal Error: RBAC Client not initialized in RBAC Handler")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "RBAC service unavailable")
		return
	}

	roleID, err := request.GetParamId(c)

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	var reqDTO dto.UpdateRoleRequestDTO

	if !util.BindAndValidateJSON(c, &reqDTO) {
		return
	}

	grpcRequest := &rbacv1.UpdateRoleRequest{
		RoleId:        roleID,
		Name:          reqDTO.Name,
		Description:   reqDTO.Description,
		PermissionIds: reqDTO.PermissionIDs,
	}

	updateRes, err := h.RbacClient.UpdateRole(c.Request.Context(), grpcRequest)
	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	// Fetch the "Fresh" Role using the ID returned from the update
	// We use GetRoleByID to ensure we have all fields (slug, permissions, etc.)
	fullRoleRes, err := h.RbacClient.GetRoleByID(c.Request.Context(), &rbacv1.GetRoleByIDRequest{
		RoleId: updateRes.Role.Id,
	})
	if err != nil {
		// If the fetch fails after a successful update, we still report an error
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	resp := dto.UpdateRoleResponseDTO{
		Role:    mapper.MapRoleProtoToDTO(fullRoleRes.Role),
		Message: "Role updated successfully",
	}

	response.JSON(c, http.StatusOK, resp)
}
