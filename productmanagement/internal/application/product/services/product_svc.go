package services

import (
	"context"
	"fmt"
	"log"
	"strings"

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

// CreateProduct orchestrates creation of a new product
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

	if req.Status == "" {
		req.Status = entity.ProductStatusActive
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

// UpdateProduct handles business-level product updates.
//
// =========================
// DESIGN INTENT
// =========================
//
// This function acts as the application-layer gatekeeper for product updates.
//
// Responsibilities:
//
// 1. VALIDATION
//   - Ensures required identifiers are present
//   - Prevents invalid state (e.g., negative stock)
//
// 2. RBAC ENFORCEMENT
//   - Only tenant owner/admin can update products
//
// 3. VERSION-SAFE PROPERTIES
//   - Ensures ProductProperties always have valid version + data
//
// 5. DELEGATION
//   - Calls repository for transactional update
func (s *ProductService) UpdateProduct(
	ctx context.Context,
	req product.UpdateProductRequestDTO,
) (*entity.Product, error) {

	const op = "ProductService.UpdateProduct"

	log.Printf("[%s]: request received: %+v", op, req)

	// -------------------------
	// Basic validation
	// -------------------------
	if req.ProductID == "" || req.TenantID == "" || req.UserID == "" {
		return nil, fmt.Errorf("[%s]: product_id, tenant_id and userID are required", op)
	}

	if req.Stock != nil && *req.Stock < 0 {
		return nil, fmt.Errorf("[%s]: stock cannot be negative", op)
	}

	// -------------------------
	// Normalize properties (version-safe)
	// -------------------------
	if req.Properties != nil {
		if req.Properties.Version == 0 {
			req.Properties.Version = 1
		}
		if req.Properties.Data == nil {
			req.Properties.Data = map[string]any{}
		}
	}

	// -------------------------
	// Tenant context + RBAC
	// -------------------------
	tenantID, err := value_object.NewTenantID(req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: invalid tenant_id: %w", op, err)
	}

	productID, err := value_object.NewProductID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: invalid product_id: %w", op, err)
	}

	tenantCtxReq := &tenantv1.GetTenantContextRequest{
		TenantId: tenantID.String(),
		UserId:   req.UserID,
	}

	tenantCtx, err := s.TenantSvcClient.GetTenantContext(ctx, tenantCtxReq)
	if err != nil || tenantCtx == nil || tenantCtx.Tenant == nil {
		return nil, fmt.Errorf("[%s]: tenant context not found", op)
	}

	if tenantCtx.Role != rbac.RoleOwner && tenantCtx.Role != rbac.RoleAdmin {
		return nil, fmt.Errorf("[%s]: unauthorized", op)
	}

	// -------------------------
	// Map incoming update
	// -------------------------
	update := entity.Product{
		Name:        *req.Name,
		Description: *req.Description,
		Price:       *req.Price,
		Currency:    *req.Currency,
		SKU:         *req.SKU,
		Stock:       *req.Stock,
		Status:      *req.Status,
		Properties:  req.Properties,
	}

	// -------------------------
	// Call repository
	// -------------------------
	updatedProduct, err := s.ProductRepo.UpdateProduct(
		ctx,
		tenantID,
		productID,
		update,
	)
	if err != nil {
		log.Printf("[%s]: update failed: %v", op, err)
		return nil, err
	}

	log.Printf("[%s]: product updated successfully: id=%s", op, updatedProduct.ID.String())

	return updatedProduct, nil
}

// AddProductImages handles adding new images to a product.
//
// =========================
// DESIGN INTENT
// =========================
//
// - Allows incremental image addition
// - Does NOT affect existing images
// - Enforces RBAC (owner/admin)
//
// =========================
// NOTES
// =========================
// - Primary image enforcement is handled at repo level
func (s *ProductService) AddProductImages(
	ctx context.Context,
	req product.AddProductImagesRequestDTO,
) error {

	const op = "ProductService.AddProductImages"

	log.Printf("[%s]: request received: %+v", op, req)

	// -------------------------
	// Validate input
	// -------------------------
	if req.ProductID == "" || req.TenantID == "" || req.UserID == "" {
		return fmt.Errorf("[%s]: product_id, tenant_id and user_id are required", op)
	}

	if len(req.Images) == 0 {
		return fmt.Errorf("[%s]: no images provided", op)
	}

	for _, img := range req.Images {
		if strings.TrimSpace(img.URL) == "" {
			return fmt.Errorf("[%s]: image url cannot be empty", op)
		}
	}

	// -------------------------
	// Tenant + RBAC
	// -------------------------
	tenantID, err := value_object.NewTenantID(req.TenantID)
	if err != nil {
		return fmt.Errorf("[%s]: invalid tenant_id: %w", op, err)
	}

	productID, err := value_object.NewProductID(req.ProductID)
	if err != nil {
		return fmt.Errorf("[%s]: invalid product_id: %w", op, err)
	}

	tenantCtx, err := s.TenantSvcClient.GetTenantContext(ctx, &tenantv1.GetTenantContextRequest{
		TenantId: tenantID.String(),
		UserId:   req.UserID,
	})

	if err != nil || tenantCtx == nil || tenantCtx.Tenant == nil {
		return fmt.Errorf("[%s]: tenant context not found", op)
	}

	if tenantCtx.Role != rbac.RoleOwner && tenantCtx.Role != rbac.RoleAdmin {
		return fmt.Errorf("[%s]: unauthorized", op)
	}

	// -------------------------
	// Map + normalize images
	// -------------------------
	var images []entity.ProductImage

	for i, imgDTO := range req.Images {
		sortOrder := imgDTO.SortOrder
		if sortOrder == 0 {
			sortOrder = i
		}

		images = append(images, entity.ProductImage{
			ImageUrl:  imgDTO.URL,
			SortOrder: sortOrder,
			IsPrimary: imgDTO.IsPrimary,
		})
	}

	// -------------------------
	// Call repo
	// -------------------------
	if err := s.ProductRepo.AddProductImages(ctx, productID, images); err != nil {
		log.Printf("[%s]: failed: %v", op, err)
		return err
	}

	log.Printf("[%s]: images added successfully for product %s", op, productID.String())

	return nil
}

// DeleteProductImages handles soft deletion of specific product images.
//
// =========================
// DESIGN INTENT
// =========================
//
// - Allows selective image deletion
// - Does NOT affect other images
// - Maintains soft delete semantics
func (s *ProductService) DeleteProductImages(
	ctx context.Context,
	req product.DeleteProductImagesRequestDTO,
) error {

	const op = "ProductService.DeleteProductImages"

	log.Printf("[%s]: request received: %+v", op, req)

	// -------------------------
	// Validate input
	// -------------------------
	if req.ProductID == "" || req.TenantID == "" || req.UserID == "" {
		return fmt.Errorf("[%s]: product_id, tenant_id and user_id are required", op)
	}

	if len(req.ImageIDs) == 0 {
		return fmt.Errorf("[%s]: no image_ids provided", op)
	}

	// -------------------------
	// Tenant + RBAC
	// -------------------------
	tenantID, err := value_object.NewTenantID(req.TenantID)
	if err != nil {
		return fmt.Errorf("[%s]: invalid tenant_id: %w", op, err)
	}

	productID, err := value_object.NewProductID(req.ProductID)
	if err != nil {
		return fmt.Errorf("[%s]: invalid product_id: %w", op, err)
	}

	tenantCtx, err := s.TenantSvcClient.GetTenantContext(ctx, &tenantv1.GetTenantContextRequest{
		TenantId: tenantID.String(),
		UserId:   req.UserID,
	})

	if err != nil || tenantCtx == nil || tenantCtx.Tenant == nil {
		return fmt.Errorf("[%s]: tenant context not found", op)
	}

	if tenantCtx.Role != rbac.RoleOwner && tenantCtx.Role != rbac.RoleAdmin {
		return fmt.Errorf("[%s]: unauthorized", op)
	}

	// -------------------------
	// Map IDs
	// -------------------------
	var ids []value_object.ProductImageID

	for _, idStr := range req.ImageIDs {
		id, err := value_object.NewProductImageID(idStr)
		if err != nil {
			return fmt.Errorf("[%s]: invalid image_id: %w", op, err)
		}
		ids = append(ids, id)
	}

	// -------------------------
	// Call repo
	// -------------------------
	if err := s.ProductRepo.SoftDeleteProductImages(ctx, productID, ids); err != nil {
		log.Printf("[%s]: failed: %v", op, err)
		return err
	}

	log.Printf("[%s]: images deleted successfully for product %s", op, productID.String())

	return nil
}

// SoftDeleteProduct handles product deletion (soft delete).
//
// =========================
// DESIGN INTENT
// =========================
//
// - Marks product as deleted
// - Cascades to images (handled in repo)
// - Enforces RBAC
func (s *ProductService) SoftDeleteProduct(
	ctx context.Context,
	req product.DeleteProductRequestDTO,
) (*entity.Product, error) {

	const op = "ProductService.SoftDeleteProduct"

	log.Printf("[%s]: request received: %+v", op, req)

	// -------------------------
	// Validate input
	// -------------------------
	if req.ProductID == "" || req.TenantID == "" || req.UserID == "" {
		return nil, fmt.Errorf("[%s]: product_id, tenant_id and user_id are required", op)
	}

	// -------------------------
	// Tenant + RBAC
	// -------------------------
	tenantID, err := value_object.NewTenantID(req.TenantID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: invalid tenant_id: %w", op, err)
	}

	productID, err := value_object.NewProductID(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: invalid product_id: %w", op, err)
	}

	tenantCtx, err := s.TenantSvcClient.GetTenantContext(ctx, &tenantv1.GetTenantContextRequest{
		TenantId: tenantID.String(),
		UserId:   req.UserID,
	})

	if err != nil || tenantCtx == nil || tenantCtx.Tenant == nil {
		return nil, fmt.Errorf("[%s]: tenant context not found", op)
	}

	if tenantCtx.Role != rbac.RoleOwner && tenantCtx.Role != rbac.RoleAdmin {
		return nil, fmt.Errorf("[%s]: unauthorized", op)
	}

	// -------------------------
	// Call repo
	// -------------------------
	product, err := s.ProductRepo.SoftDeleteProduct(ctx, tenantID, productID)
	if err != nil {
		log.Printf("[%s]: failed: %v", op, err)
		return nil, err
	}

	log.Printf("[%s]: product deleted successfully: %s", op, productID.String())

	return product, nil
}
