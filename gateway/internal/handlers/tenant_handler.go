package handlers

import (
	"log"
	"net/http"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/mapper"
	"github.com/MartinMurithi/storeforge/gateway/internal/response"
	"github.com/MartinMurithi/storeforge/gateway/internal/util"
	"github.com/MartinMurithi/storeforge/pkg/auth"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type TenantHandler struct {
	TenantClient tenantv1.TenantServiceClient
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {

	if h.TenantClient == nil {
		log.Println("Internal Error: TenantClient not initialized in TenantHandler")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Tenant service unavailable")
		return
	}

	userID, _ := c.Get(auth.CtxUserID)

	log.Printf("id of active user %s", userID)

	var req tenantv1.CreateTenantRequest

	if !util.BindAndValidateJSON(c, &req) {
		return
	}

	// Setting the Metadata
	md := metadata.Pairs("user-id", userID.(string))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	resp, err := h.TenantClient.CreateTenant(ctx, &req)
	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	response.JSON(c, http.StatusCreated, mapper.MapCreateTenantResponseProtoToDTO(resp))
}
