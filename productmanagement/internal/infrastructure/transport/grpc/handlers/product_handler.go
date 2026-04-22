package handlers

import (
	"context"
	"fmt"
	"log"

	productv1 "github.com/MartinMurithi/storeforge/api/protos/productmanagement/product/v1"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
	"github.com/MartinMurithi/storeforge/pkg/grpcx"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/application/product"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/application/product/services"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/entity"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductGrpcHandler struct {
	productv1.UnimplementedProductServiceServer
	ProductService *services.ProductService
}

// NewProductGrpcHandler initializes the handler with the required application service.
func NewProductGrpcHandler(s *services.ProductService) *ProductGrpcHandler {
	if s == nil {
		panic("NewProductGrpcHandler: service is nil")
	}
	return &ProductGrpcHandler{
		ProductService: s,
	}
}

// CreateProduct converts the protobuf request into an internal DTO,
// executes the creation logic, and returns a mapped protobuf response.
func (h *ProductGrpcHandler) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.CreateProductResponse, error) {
	const op = "ProductHandler.CreateProduct"

	tenantID, err := grpcx.GetTenantIDFromMetadata(ctx)

	log.Printf("[%s]: extracted tenant id %s", op, tenantID)

	if err != nil {
		log.Printf("[%s]: failed to extract tenant id from metadata %s", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing tenant identity: %v", err)
	}

	userID, err := grpcx.GetUserIDFromMetadata(ctx)
	log.Printf("[%s]: extracted user id %s", op, userID)

	if err != nil {
		log.Printf("[%s]: failed to extract user id from metadata %s", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing user identity: %v", err)
	}

	log.Printf("[%s]: extracted user id %s", op, userID)

	// Set the tenant and user ID in the grpc metadata and send them to tenant service
	ctx = grpcx.ForwardMetadata(ctx)

	var props *entity.ProductProperties

	if req.Properties != nil {
		m := req.Properties.AsMap()

		version := 1
		if v, ok := m["version"].(float64); ok {
			version = int(v)
			delete(m, "version")
		}

		props = &entity.ProductProperties{
			Version: version,
			Data:    m,
		}
	}

	var prodStatus entity.ProductStatus

	if req.Status != "" {
		s := entity.ProductStatus(req.Status)
		prodStatus = s
	}

	dtoReq := product.CreateProductRequestDTO{
		TenantID:    tenantID,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SKU:         req.Sku,
		Stock:       req.Stock,
		Status:      prodStatus,
		Properties:  props,
	}

	// Call the service
	result, err := h.ProductService.CreateProduct(ctx, dtoReq)
	if err != nil {
		return nil, errconv.ToGrpcError(err)
	}

	res := &productv1.CreateProductResponse{
		Product: product.ToProtoProduct(result),
		Message: "Product Created Successfully",
	}

	return res, nil
}

func (h *ProductGrpcHandler) GetTenantProducts(ctx context.Context, req *productv1.GetTenantProductsRequest) (*productv1.GetTenantProductsResponse, error) {
	tenantID, err := grpcx.GetTenantIDFromMetadata(ctx)

	log.Printf("extracted tenant id from metadata %s", tenantID)

	if err != nil {
		log.Printf("failed to extract tenant id from metadata %s", err)
		return nil, status.Errorf(codes.Unauthenticated, "missing tenant identity: %v", err)
	}

	pagination := &product.Pagination{
		Page:  int(req.Page),
		Limit: int(req.Limit),
	}

	products, meta, err := h.ProductService.GetProductsByTenant(ctx, tenantID, *pagination)

	log.Printf("metadata %v", meta.Total)

	if err != nil {
		fmt.Printf("[PRODUCTGRPCHANDLER]: failed to fetch products %s\n", err)
		return nil, errconv.ToGrpcError(err)
	}

	return product.ToProtoFetchTenantProductsResponse(products, &meta), nil
}

func (h *ProductGrpcHandler) GetProductByID(ctx context.Context, req *productv1.GetProductByIDRequest) (*productv1.GetProductByIDResponse, error) {
	tenantID, err := grpcx.GetTenantIDFromMetadata(ctx)

	log.Printf("extracted tenant id from metadata %s", tenantID)

	if err != nil {
		log.Printf("failed to extract tenant id from metadata %s", err)
		return nil, status.Errorf(codes.Unauthenticated, "missing tenant identity: %v", err)
	}

	prod, err := h.ProductService.GetProductByID(ctx, tenantID, req.ProductId)

	if err != nil {
		fmt.Printf("[PRODUCTGRPCHANDLER]: failed to fetch product %s\n", err)
		return nil, errconv.ToGrpcError(err)
	}

	return &productv1.GetProductByIDResponse{Product: product.ToProtoProduct(prod)}, nil
}

func (h *ProductGrpcHandler) AddProductImages(
	ctx context.Context,
	req *productv1.AddProductImagesRequest,
) (*productv1.AddProductImagesResponse, error) {

	const op = "ProductGrpcHandler.AddProductImages"

	tenantID, err := grpcx.GetTenantIDFromMetadata(ctx)

	log.Printf("[%s]: extracted tenant id %s", op, tenantID)

	if err != nil {
		log.Printf("[%s]: failed to extract tenant id from metadata %s", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing tenant identity: %v", err)
	}

	userID, err := grpcx.GetUserIDFromMetadata(ctx)
	log.Printf("[%s]: extracted user id %s", op, userID)

	if err != nil {
		log.Printf("[%s]: failed to extract user id from metadata %s", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing user identity: %v", err)
	}

	log.Printf("[%s]: extracted user id %s", op, userID)

	// Set the tenant and user ID in the grpc metadata and send them to tenant service
	ctx = grpcx.ForwardMetadata(ctx)

	// -------------------------
	// Map request DTO
	// -------------------------
	dto := product.AddProductImagesRequestDTO{
		ProductID: req.ProductId,
		TenantID:  tenantID,
		UserID:    userID,
	}

	for _, img := range req.Images {
		dto.Images = append(dto.Images, product.AddProductImageInputDTO{
			URL:       img.ImageUrl,
			SortOrder: int(img.SortOrder),
			IsPrimary: img.IsPrimary,
		})
	}

	// -------------------------
	// Service call
	// -------------------------
	err = h.ProductService.AddProductImages(ctx, dto)
	if err != nil {
		log.Printf("[%s]: failed to add images: %v", op, err)
		return nil, errconv.ToGrpcError(err)
	}

	return &productv1.AddProductImagesResponse{
		Message: "images added successfully",
	}, nil
}

func (h *ProductGrpcHandler) DeleteProductImages(
	ctx context.Context,
	req *productv1.DeleteProductImagesRequest,
) (*productv1.DeleteProductImagesResponse, error) {

	const op = "ProductGrpcHandler.DeleteProductImages"

	tenantID, err := grpcx.GetTenantIDFromMetadata(ctx)
	log.Printf("[%s]: extracted tenant id %s", op, tenantID)

	if err != nil {
		log.Printf("[%s]: failed to extract tenant id: %v", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing tenant identity: %v", err)
	}

	userID, err := grpcx.GetUserIDFromMetadata(ctx)
	log.Printf("[%s]: extracted user id %s", op, userID)

	if err != nil {
		log.Printf("[%s]: failed to extract user id: %v", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing user identity: %v", err)
	}

	// Set the tenant and user ID in the grpc metadata and send them to tenant service
	ctx = grpcx.ForwardMetadata(ctx)

	// -------------------------
	// Map DTO
	// -------------------------
	dto := product.DeleteProductImagesRequestDTO{
		ProductID: req.ProductId,
		TenantID:  tenantID,
		UserID:    userID,
		ImageIDs:  req.ImageIds,
	}

	// -------------------------
	// Service call
	// -------------------------
	err = h.ProductService.DeleteProductImages(ctx, dto)
	if err != nil {
		log.Printf("[%s]: failed to delete images: %v", op, err)
		return nil, errconv.ToGrpcError(err)
	}

	return &productv1.DeleteProductImagesResponse{
		Message: "images deleted successfully",
	}, nil
}

func (h *ProductGrpcHandler) UpdateProduct(
	ctx context.Context,
	req *productv1.UpdateProductRequest,
) (*productv1.UpdateProductResponse, error) {

	const op = "ProductGrpcHandler.UpdateProduct"

	tenantID, err := grpcx.GetTenantIDFromMetadata(ctx)
	log.Printf("[%s]: extracted tenant id %s", op, tenantID)

	if err != nil {
		log.Printf("[%s]: failed to extract tenant id: %v", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing tenant identity: %v", err)
	}

	userID, err := grpcx.GetUserIDFromMetadata(ctx)
	log.Printf("[%s]: extracted user id %s", op, userID)

	if err != nil {
		log.Printf("[%s]: failed to extract user id: %v", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing user identity: %v", err)
	}

	// Set the tenant and user ID in the grpc metadata and send them to tenant service
	ctx = grpcx.ForwardMetadata(ctx)

	// -------------------------
	// Map properties (protobuf Struct → domain)
	// -------------------------
	var props *entity.ProductProperties

	if req.Input != nil && req.Input.ProductProperties != nil {
		var dataMap map[string]any
		if req.Input.ProductProperties.Properties != nil {
			dataMap = req.Input.ProductProperties.Properties.AsMap()
		}

		props = &entity.ProductProperties{
			Version: int(req.Input.ProductProperties.Version),
			Data:    dataMap,
		}
	}

	// -------------------------
	// Map status
	// -------------------------
	var statusVal *entity.ProductStatus
	if req.Input != nil && req.Input.Status != nil {
		s := entity.ProductStatus(*req.Input.Status)
		statusVal = &s
	}

	dto := product.UpdateProductRequestDTO{
		ProductID:  req.ProductId,
		TenantID:   tenantID,
		UserID:     userID,
		Status:     statusVal,
		Properties: props,
	}

	if req.Input != nil {
		dto.Name = req.Input.Name
		dto.Description = req.Input.Description
		dto.Price = req.Input.Price
		dto.Currency = req.Input.Currency
		dto.SKU = req.Input.Sku
		dto.Stock = req.Input.Stock
	}

	updatedProduct, err := h.ProductService.UpdateProduct(ctx, dto)
	if err != nil {
		log.Printf("[%s]: failed to update product: %v", op, err)
		return nil, errconv.ToGrpcError(err)
	}

	return &productv1.UpdateProductResponse{
		Message: "Product Updated Successfully",
		Product: product.ToProtoProduct(updatedProduct),
	}, nil
}

func (h *ProductGrpcHandler) SoftDeleteProduct(
	ctx context.Context,
	req *productv1.SoftDeleteProductRequest,
) (*productv1.SoftDeleteProductResponse, error) {

	const op = "ProductGrpcHandler.SoftDeleteProduct"

	tenantID, err := grpcx.GetTenantIDFromMetadata(ctx)
	log.Printf("[%s]: extracted tenant id %s", op, tenantID)

	if err != nil {
		log.Printf("[%s]: failed to extract tenant id: %v", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing tenant identity: %v", err)
	}

	userID, err := grpcx.GetUserIDFromMetadata(ctx)
	log.Printf("[%s]: extracted user id %s", op, userID)

	if err != nil {
		log.Printf("[%s]: failed to extract user id: %v", op, err)
		return nil, status.Errorf(codes.Unauthenticated, "missing user identity: %v", err)
	}

	// Set the tenant and user ID in the grpc metadata and send them to tenant service
	ctx = grpcx.ForwardMetadata(ctx)

	dto := product.DeleteProductRequestDTO{
		ProductID: req.ProductId,
		TenantID:  tenantID,
		UserID:    userID,
	}

	if err := h.ProductService.SoftDeleteProduct(ctx, dto); err != nil {
		log.Printf("[%s]: failed to delete product: %v", op, err)
		return nil, errconv.ToGrpcError(err)
	}

	return &productv1.SoftDeleteProductResponse{
		Message: "product deleted successfully",
	}, nil
}
