package services

import (
	"context"
	"fmt"
	"log"

	tenantv1 "github.com/MartinMurithi/storeforge/api/protos/tenantmanagement/tenant/v1"
	"github.com/MartinMurithi/storeforge/pkg/rbac"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/application/product"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/entity"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/repository"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/value_object"
	grpcclient "github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/grpc_client"
)

type ProductService struct {
	ProductRepo     repository.IProductRepository
	TenantSvcClient grpcclient.TenantSvcClient
}

// NewProductService creates a new instance of the ProductService.
func NewProductService(pr repository.IProductRepository, tenantClientSvc *grpcclient.TenantSvcClient) *ProductService {
	return &ProductService{
		ProductRepo:     pr,
		TenantSvcClient: *tenantClientSvc,
	}
}

// CreateProduct orchestrates creation of a new product and its images.
// Handles stock/price validation, image normalization (first image = primary),
// and persists everything atomically.
func (s *ProductService) CreateProduct(ctx context.Context, req product.CreateProductRequestDTO) (*entity.Product, error) {
	const op = "ProductService.CreateProduct"

	log.Printf("[%s]: request received: %+v", op, req)

	// -------------------------
	// Validate required fields
	// -------------------------
	if req.Name == "" || req.TenantID == "" || req.UserID == "" || req.SKU == "" {
		return nil, fmt.Errorf("[%s]: tenant_id, user_id, sku and name are required", op)
	}

	if req.Stock < 0 {
		return nil, fmt.Errorf("[%s]: stock cannot be negative", op)
	}

	if req.Properties == nil {
		req.Properties = &entity.ProductProperties{
			Version: 1,
			Data:    map[string]any{},
		}
	}

	if req.Properties.Version == 0 {
		req.Properties.Version = 1
	}

	if req.Properties.Data == nil {
		req.Properties.Data = map[string]any{}
	}

	// -------------------------
	// TenantID value object
	// -------------------------
	log.Printf("received tenant id %s", req.TenantID)
	tenantID, err := value_object.NewTenantID(req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: invalid tenant_id: %w", op, err)
	}

	// Get tenant context, this wil return the current tenant and the currently loggedin entity
	// Call service from tenant service via grpc clients

	tenantCtxReq := &tenantv1.GetTenantContextRequest{
		TenantId: tenantID.String(),
		UserId:   req.UserID,
	}

	tenantCtx, err := s.TenantSvcClient.GetTenantContext(ctx, tenantCtxReq)

	if tenantCtx == nil || tenantCtx.Tenant == nil {
		return nil, fmt.Errorf("[%s]: tenant context not found", op)
	}

	log.Printf("user role and status : %s, %s", tenantCtx.Role, tenantCtx.Tenant.StoreName)

	if tenantCtx.Role != rbac.RoleOwner && tenantCtx.Role != rbac.RoleAdmin {
		return nil, fmt.Errorf("unauthorized, only admin and owner are allowed to create a product")
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

	return newProduct, nil
}

func (s *ProductService) GetProductsByTenant(ctx context.Context, tenantID string, p product.Pagination) ([]*entity.Product, product.PaginationMeta, error) {
	id, err := value_object.NewTenantID(tenantID)
	if err != nil {
		return nil, product.PaginationMeta{}, err
	}

	products, total, err := s.ProductRepo.GetProductsByTenant(ctx, id, p)
	if err != nil {
		return nil, product.PaginationMeta{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + p.Limit - 1) / p.Limit
	}

	meta := product.PaginationMeta{
		Page:       p.Page,
		Limit:      p.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    p.Page < totalPages,
		HasPrev:    p.Page > 1,
	}

	return products, meta, nil
}

func (s *ProductService) GetProductByID(ctx context.Context, tenantID string, productID string) (*entity.Product, error) {
	tID, err := value_object.NewTenantID(tenantID)
	if err != nil {
		return nil, err
	}

	pID, err := value_object.NewProductID(productID)
	if err != nil {
		return nil, err
	}

	product, err := s.ProductRepo.GetProductByID(ctx, tID, pID)

	if err != nil {
		return nil, err
	}

	return product, nil
}
