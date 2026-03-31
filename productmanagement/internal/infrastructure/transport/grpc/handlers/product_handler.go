package handlers

import (
	"context"
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

	log.Printf("extracted tenant id %s", tenantID)

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

	return result, nil
}
