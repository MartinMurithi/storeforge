package handlers

import (
	"fmt"
	"log"
	"net/http"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/mapper"
	"github.com/MartinMurithi/storeforge/gateway/internal/request"
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

func (h *TenantHandler) GetTenantContext(c *gin.Context) {
	if h.TenantClient == nil {
		log.Println("Internal Error: TenantClient not initialized in TenantHandler")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Tenant service unavailable")
		return
	}

	// Get logged in owner
	userID, err := request.GetUserId(c)

	fmt.Println("error getting user: %w", err)

	if err != nil {
		fmt.Println("error getting user ID: %w", err)
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "User session not found")
		return
	}

	// Setting the Metadata
	md := metadata.Pairs("user-id", userID)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	// gET TENANT ID FROM PARAMS
	tenantID, err := request.GetParamId(c)

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	resp, err := h.TenantClient.GetTenantContext(ctx, &tenantv1.GetTenantContextRequest{
		TenantId: tenantID,
		UserId:   userID,
	})
	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	response.JSON(c, http.StatusAccepted, mapper.MapGetTenantTenantContextResponse(resp))
}

func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	if h.TenantClient == nil {
		log.Println("Internal Error: TenantClient not initialized in TenantHandler")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Tenant service unavailable")
		return
	}

	var req tenantv1.UpdateTenantRequest

	if !util.BindAndValidateJSON(c, &req) {
		return
	}

	// Get logged in owner
	userID, err := request.GetUserId(c)

	fmt.Println("error getting user: %w", err)

	if err != nil {
		fmt.Println("error getting user ID: %w", err)
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "User session not found")
		return
	}

	// Setting the Metadata
	md := metadata.Pairs("user-id", userID)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	// gET TENANT ID FROM PARAMS
	tenantID, err := request.GetParamId(c)

	if err != nil {
		// code, slug, msg := errconv.FromGrpcToHttp(err)
		fmt.Println("error getting tenant ID: %w", err)
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Tenant ID not found")
		return
	}

	resp, err := h.TenantClient.UpdateTenant(ctx, &tenantv1.UpdateTenantRequest{
		TenantId: tenantID,
		UserId:   userID,
		Settings: &tenantv1.TenantSettingsUpdate{
			ThemeConfig: req.Settings.ThemeConfig,
		},
	})
	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	response.JSON(c, http.StatusCreated, mapper.MapGetTenantTenantContextResponse(resp))
}
