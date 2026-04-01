package handlers

import (
	"fmt"
	"log"
	"net/http"

	productv1 "github.com/MartinMurithi/storeforge/api/protos/productmanagement/product/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/mapper"
	"github.com/MartinMurithi/storeforge/gateway/internal/request"
	"github.com/MartinMurithi/storeforge/gateway/internal/response"
	"github.com/MartinMurithi/storeforge/gateway/internal/util"
	"github.com/MartinMurithi/storeforge/pkg/auth"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

type ProductHandler struct {
	ProductClient productv1.ProductServiceClient
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	if h.ProductClient == nil {
		log.Println("Internal Error: ProductClient not initialized in ProductHandler")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Product service unavailable")
		return
	}

	// // Get logged in owner
	// userID, err := request.GetUserId(c)

	// fmt.Println("error getting user: %w", err)

	// if err != nil {
	// 	fmt.Println("error getting user ID: %w", err)
	// 	response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "User session not found")
	// 	return
	// }

	userID, _ := c.Get(auth.CtxUserID)

	log.Printf("id of active user %s", userID)

	var req productv1.CreateProductRequest

	if !util.BindAndValidateJSON(c, &req) {
		return
	}

// Setting the Metadata
	md := metadata.Pairs("user-id",userID)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	// gET TENANT ID FROM PARAMS
	tenantID, err := request.GetParamId(c)

	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
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
