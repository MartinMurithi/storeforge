package services

import (
	"context"
	"fmt"
	"log"

	productv1 "github.com/MartinMurithi/storeforge/api/protos/productmanagement/product/v1"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/application/product"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/entity"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/repository"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/value_object"
)

type ProductService struct {
	ProductRepo repository.IProductRepository
}

// NewProductService creates a new instance of the ProductService.
func NewProductService(pr repository.IProductRepository) *ProductService {
	return &ProductService{
		ProductRepo: pr,
	}
}

// CreateProduct orchestrates creation of a new product and its images.
// Handles stock/price validation, image normalization (first image = primary),
// and persists everything atomically.
func (s *ProductService) CreateProduct(ctx context.Context, req product.CreateProductRequestDTO) (*productv1.CreateProductResponse, error) {
	const op = "ProductService.CreateProduct"

	log.Printf("[%s]: request received: %+v", op, req)

	// -------------------------
	// Validate required fields
	// -------------------------
	if req.Name == "" || req.TenantID == "" || req.SKU == "" {
		return nil, fmt.Errorf("[%s]: tenant_id, sku and name are required", op)
	}

	if req.Stock < 0 {
		return nil, fmt.Errorf("[%s]: stock cannot be negative", op)
	}

	if req.Properties == nil {
		req.Properties = &entity.ProductProperties{} // empty properties map
	}

	// -------------------------
	// TenantID value object
	// -------------------------
	tenantID, err := value_object.NewTenantID(req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: invalid tenant_id: %w", op, err)
	}

	newProduct := &entity.Product{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		SKU:         req.SKU,
		Stock:       req.Stock,
		Status:      req.Status,
		Properties:  req.Properties,
		Price:       req.Price,
	}

	// -------------------------
	// Map and normalize images
	// -------------------------
	for _, imgDTO := range req.Images {
		newProduct.ProductImages = append(newProduct.ProductImages, entity.ProductImage{
			ImageUrl:  imgDTO.URL,
			SortOrder: imgDTO.SortOrder,
			IsPrimary: imgDTO.IsPrimary,
		})
	}

	if len(newProduct.ProductImages) > 0 {
		for i := range newProduct.ProductImages {
			if newProduct.ProductImages[i].SortOrder == 0 {
				newProduct.ProductImages[i].SortOrder = i
			}
			newProduct.ProductImages[i].IsPrimary = false
		}
		newProduct.ProductImages[0].IsPrimary = true
	}

	if err := s.ProductRepo.CreateProduct(ctx, newProduct); err != nil {
		log.Printf("[%s]: error creating product: %v", op, err)
		return nil, err
	}

	log.Printf("[%s]: product created successfully: id=%s, name=%s", op, newProduct.ID.String(), newProduct.Name)

	productResp := product.ToProtoProduct(newProduct)

	res := &productv1.CreateProductResponse{
		Product: productResp,
		Message: "Product Created Successfully",
	}

	return res, nil
}
