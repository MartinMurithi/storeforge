package handlers

import (
	"log"
	"net/http"

	productv1 "github.com/MartinMurithi/storeforge/api/protos/productmanagement/product/v1"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto"
	"github.com/MartinMurithi/storeforge/gateway/internal/dto/shared"
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

func (h *ProductHandler) GetTenantProducts(c *gin.Context) {
	const op = "ProductHandler.GetTenantProducts"

	if h.ProductClient == nil {
		log.Println("Internal Error: ProductClient not initialized in ProductHandler")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Product service unavailable")
		return
	}

	// gET TENANT ID FROM PARAMS
	tenantID, err := request.GetParamId(c)
	if err != nil {
		log.Printf("[%s]: error getting tenant ID: %v", op, err)
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Tenant identification not found")
		return
	}
	log.Printf("[%s]: current tenant ID : %s", op, tenantID)

	// Setting the Metadata for the product grpc service
	md := metadata.Pairs("tenant-id", tenantID)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	// Get pagination params
	pagination, err := shared.ParsePagination(c)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "INVALID_PAGINATION", err.Error())
		return
	}

	resp, err := h.ProductClient.GetTenantProducts(ctx, &productv1.GetTenantProductsRequest{
		TenantId: tenantID,
		Page:     pagination.Page,
		Limit:    pagination.Limit,
	})
	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	res := mapper.MapGetTenantProductsResponse(resp)

	response.JSON(c, http.StatusAccepted, dto.GetTenantProductsResponseDTO{
		Products: res,
		Meta: shared.PaginationMetaDTO{
			Page:       resp.Meta.Page,
			Limit:      resp.Meta.Limit,
			Total:      int64(resp.Meta.Total),
			TotalPages: resp.Meta.TotalPages,
			HasNext:    resp.Meta.HasNext,
			HasPrev:    resp.Meta.HasPrev,
		},
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	const op = "ProductHandler.GetProductByID"

	if h.ProductClient == nil {
		log.Println("Internal Error: ProductClient not initialized in ProductHandler")
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Product service unavailable")
		return
	}

	// gET Tenant ID FROM PARAMS
	tenantID, err := request.GetParamId(c)
	if err != nil {
		log.Printf("[%s]: error getting tenant ID: %v", op, err)
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Tenant ID not found")
		return
	}
	log.Printf("[%s]: tenant ID : %s", op, tenantID)

	// gET Product ID FROM PARAMS
	productID, err := request.GetNamedParamID(c, "productID")
	if err != nil {
		log.Printf("[%s]: error getting product ID: %v", op, err)
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Product ID not found")
		return
	}
	log.Printf("[%s]: product ID : %s", op, productID)

	// Setting the Metadata for the product grpc service
	md := metadata.Pairs(
		"tenant-id", tenantID,
		"product-id", productID,
	)
	ctx := metadata.NewOutgoingContext(c.Request.Context(), md)

	resp, err := h.ProductClient.GetProductByID(ctx, &productv1.GetProductByIDRequest{
		TenantId:  tenantID,
		ProductId: productID,
	})
	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	res := mapper.MapGetProductResponse(resp)

	response.JSON(c, http.StatusAccepted, res)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	const op = "ProductHandler.UpdateProduct"

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

	var req productv1.UpdateProductRequest

	if !util.BindAndValidateJSON(c, &req) {
		return
	}

	resp, err := h.ProductClient.UpdateProduct(ctx, &req)
	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	response.JSON(c, http.StatusAccepted, &productv1.UpdateProductResponse{
		Message: "Product Updated Successfully",
		Product: resp.Product,
	})
}

func (h *ProductHandler) SoftDeleteProduct(c *gin.Context) {
	const op = "ProductHandler.SoftDeleteProduct"

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

	var req productv1.SoftDeleteProductRequest

	if !util.BindAndValidateJSON(c, &req) {
		return
	}

	resp, err := h.ProductClient.SoftDeleteProduct(ctx, &req)
	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	response.JSON(c, http.StatusAccepted, &productv1.SoftDeleteProductResponse{
		Message: resp.Message,
	})
}

func (h *ProductHandler) AddProductImages(c *gin.Context) {
	const op = "ProductHandler.AddProductImages"

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

	var req productv1.AddProductImagesRequest

	if !util.BindAndValidateJSON(c, &req) {
		return
	}

	resp, err := h.ProductClient.AddProductImages(ctx, &req)
	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	response.JSON(c, http.StatusAccepted, &productv1.AddProductImagesResponse{
		Message: resp.Message,
	})
}

func (h *ProductHandler) DeleteProductImages(c *gin.Context) {
	const op = "ProductHandler.DeleteProductImages"

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

	var req productv1.DeleteProductImagesRequest

	if !util.BindAndValidateJSON(c, &req) {
		return
	}

	resp, err := h.ProductClient.DeleteProductImages(ctx, &req)
	if err != nil {
		code, slug, msg := errconv.FromGrpcToHttp(err)
		response.Error(c, code, slug, msg)
		return
	}

	response.JSON(c, http.StatusAccepted, &productv1.DeleteProductImagesResponse{
		Message: resp.Message,
	})
}
