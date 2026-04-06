package handlers

import (
	"log"
	"net/http"

	productv1 "github.com/MartinMurithi/storeforge/api/protos/productmanagement/product/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/mapper"
	"github.com/MartinMurithi/storeforge/gateway/internal/request"
	"github.com/MartinMurithi/storeforge/gateway/internal/response"
	"github.com/MartinMurithi/storeforge/gateway/internal/util"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type ProductHandler struct {
	ProductClient productv1.ProductServiceClient
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	const op = "ProductHandler.CreateProduct"

	if h.ProductClient == nil {
		log.Println("Internal Error: ProductClient not initialized in ProductHandler")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Product service unavailable")
		return
	}

	// Get logged in owner
	userID, err := request.GetUserId(c)
	if err != nil {
		log.Printf("[%s]: error getting user ID: %v", op, err)
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "User session not found")
		return
	}
	log.Printf("[%s]: active user ID : %s", op, userID)

	// gET TENANT ID FROM PARAMS
	tenantID, err := request.GetParamId(c)
	if err != nil {
		log.Printf("[%s]: error getting tenant ID: %v", op, err)
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Tenant identification not found")
		return
	}
	log.Printf("[%s]: current tenant ID : %s", op, tenantID)

	// Setting the Metadata for the product grpc service
	md := metadata.Pairs(
		"user-id", userID,
		"tenant-id", tenantID)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	var req productv1.CreateProductRequest

	if !util.BindAndValidateJSON(c, &req) {
		return
	}

	resp, err := h.ProductClient.CreateProduct(ctx, &req)
	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	response.JSON(c, http.StatusAccepted, mapper.MapCreateProductResponse(resp))
}
