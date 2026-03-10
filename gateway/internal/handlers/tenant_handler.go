package handlers

import (
	"log"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	"github.com/MartinMurithi/storeforge/pkg/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type TenantHandler struct {
	TenantClient tenantv1.TenantServiceClient
}

func (h *TenantHandler) CreateTenant(c *gin.Context) {
	userID, _ := c.Get(auth.CtxUserID)

	log.Printf("id of active user %s", userID)

	var req tenantv1.CreateTenantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	// 3. The Relay (Setting the Metadata)
	md := metadata.Pairs("user-id", userID.(string))
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	// 4. The Execution
	resp, err := h.TenantClient.CreateTenant(ctx, &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, resp)
}
