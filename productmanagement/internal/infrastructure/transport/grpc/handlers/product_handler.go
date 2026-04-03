package handlers

import (
	"context"
	"fmt"
	"log"

	productv1 "github.com/MartinMurithi/storeforge/api/protos/productmanagement/product/v1"
	"github.com/MartinMurithi/storeforge/pkg/auth"
	"github.com/MartinMurithi/storeforge/pkg/errconv"
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

	tenantID, err := auth.GetTenantIDFromMetadata(ctx)

	if err != nil {
		log.Printf("failed to extract tenant id from metadata %s", err)
		return nil, status.Errorf(codes.Unauthenticated, "missing tenant identity: %v", err)
	}

	userID, err := auth.GetUserIDFromMetadata(ctx)

	if err != nil {
		log.Printf("failed to extract user id from metadata %s", err)
		return nil, status.Errorf(codes.Unauthenticated, "missing user identity: %v", err)
	}

	log.Printf("extracted user id %s", userID)

	var props *entity.ProductProperties

	if req.Properties != nil {
		m := entity.ProductProperties(req.Properties.AsMap())
		props = &m
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
		Images:      make([]product.ProductImageInputDTO, 0, len(req.Images)),
	}
	// Map repeated images
	for _, img := range req.Images {
		dtoReq.Images = append(dtoReq.Images, product.ProductImageInputDTO{
			URL: img.ImageUrl,
		})
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
	tenantID, err := auth.GetTenantIDFromMetadata(ctx)

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
	tenantID, err := auth.GetTenantIDFromMetadata(ctx)

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
