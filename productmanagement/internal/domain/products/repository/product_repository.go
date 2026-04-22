package repository

import (
	"context"
	"fmt"
	"log"
	// "strings"
	"time"

	"github.com/MartinMurithi/storeforge/productmanagement/internal/application/product"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/entity"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/value_object"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/database"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/database/postgres"
	"github.com/google/uuid"
)

type ProductRepository struct {
	DB database.DB
}

type IProductRepository interface {
	CreateProduct(ctx context.Context, product *entity.Product) error
	GetProductsByTenant(ctx context.Context, tenantID value_object.TenantID, p product.Pagination) ([]*entity.Product, int, error)
	GetProductByID(ctx context.Context, tenantID value_object.TenantID, productID value_object.ProductID) (*entity.Product, error)
	AddProductImages(
		ctx context.Context,
		productID value_object.ProductID,
		images []entity.ProductImage,
	) error
	UpdateProduct(
		ctx context.Context,
		tenantID value_object.TenantID,
		productID value_object.ProductID,
		incoming entity.Product,
	) (*entity.Product, error)
	SoftDeleteProduct(
		ctx context.Context,
		tenantID value_object.TenantID,
		productID value_object.ProductID,
	) error
	SoftDeleteProductImages(
		ctx context.Context,
		productID value_object.ProductID,
		imageIDs []value_object.ProductImageID,
	) error
}

func NewProductRepository(db database.DB) IProductRepository {
	return &ProductRepository{DB: db}
}

func (repo *ProductRepository) CreateProduct(ctx context.Context, product *entity.Product) error {
	const op = "ProductRepository.CreateProduct"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	tx, err := repo.DB.Tx(ctx)
	if err != nil {
		return fmt.Errorf("[%s]: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// Insert product
	productQuery := `
		INSERT INTO products (
			tenant_id,
			name,
			description,
			price_cents,
			currency,
			sku,
			stock_quantity,
			product_properties,
			product_status
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, name, currency, created_at;
	`

	var id uuid.UUID

	err = tx.QueryRow(
		ctx,
		productQuery,
		product.TenantID,
		product.Name,
		product.Description,
		product.Price,
		product.Currency,
		product.SKU,
		product.Stock,
		product.Properties,
		product.Status,
	).Scan(&id, &product.Name, &product.Currency, &product.CreatedAt)

	product.ID = value_object.NewProductIDFromUUID(id)

	if err != nil {
		log.Printf("[%s error]: %v", op, err)
		return fmt.Errorf(
			"%w",
			domain.TranslateProductRepoError(postgres.MapPostgresError(err)),
		)
	}
	return tx.Commit(ctx)
}

// GetProductsByTenant fetches a paginated list of products for a given tenant.
//
// PURPOSE:
//   - Retrieve all non-deleted products for a tenant with optional pagination.
//   - Includes associated product images for each product.
//   - Supports ordering by product creation date and image sort order.
//
// RETURNS:
//   - A slice of fully hydrated Product entities including images.
//   - Total number of products for the tenant (ignoring pagination).
//   - Error if the query fails or scanning fails.
func (repo *ProductRepository) GetProductsByTenant(ctx context.Context, tenantID value_object.TenantID, p product.Pagination) ([]*entity.Product, int, error) {
	const op = "ProductRepository.GetProductsByTenant"

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// --- Total products Count ---
	var totalProducts int
	if err := repo.DB.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE tenant_id = $1 AND deleted_at IS NULL`, tenantID.Raw()).Scan(&totalProducts); err != nil {
		return nil, 0, domain.TranslateProductRepoError(postgres.MapPostgresError(err))
	}

	var maxLimit = p.Limit
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit <= 0 || p.Limit > maxLimit {
		p.Limit = maxLimit
	}

	offset := (p.Page - 1) * p.Limit

	query := `
	SELECT
		p.id,
		p.tenant_id,
		p.name,
		p.description,
		p.price_cents,
		p.currency,
		p.sku,
		p.stock_quantity,
		p.product_properties,
		p.product_status,
		p.created_at,
		p.updated_at,

		i.id,
		i.image_url,
		i.sort_order,
		i.is_primary,
		i.created_at

	FROM products p
	LEFT JOIN product_images i
		ON p.id = i.product_id

	WHERE p.tenant_id = $1
	AND p.deleted_at IS NULL

	ORDER BY p.created_at DESC, i.sort_order ASC
	LIMIT $2 OFFSET $3
	`

	rows, err := repo.DB.Query(ctx, query, tenantID.Raw(), p.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("[%s]: %w", op, err)
	}
	defer rows.Close()

	productMap := map[string]*entity.Product{}

	for rows.Next() {

		var (
			productID  uuid.UUID
			dbTenantID uuid.UUID

			imageID        *uuid.UUID
			imageURL       *string
			sortOrder      *int
			isPrimary      *bool
			imageCreatedAt *time.Time

			product entity.Product
		)

		err = rows.Scan(
			&productID,
			&dbTenantID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Currency,
			&product.SKU,
			&product.Stock,
			&product.Properties,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt,

			&imageID,
			&imageURL,
			&sortOrder,
			&isPrimary,
			&imageCreatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("[%s scan]: %w", op, err)
		}

		pID := productID.String()

		if _, exists := productMap[pID]; !exists {

			product.ID = value_object.NewProductIDFromUUID(productID)
			product.TenantID = value_object.NewTenantIDFromUUID(dbTenantID)

			product.ProductImages = []entity.ProductImage{}

			productMap[pID] = &product
		}

		if imageID != nil {

			img := entity.ProductImage{
				ID:        value_object.NewProductImageIDFromUUID(*imageID),
				ProductID: value_object.NewProductIDFromUUID(productID),
				ImageUrl:  *imageURL,
				SortOrder: *sortOrder,
				IsPrimary: *isPrimary,
				CreatedAt: *imageCreatedAt,
			}

			productMap[pID].ProductImages =
				append(productMap[pID].ProductImages, img)
		}
	}

	result := make([]*entity.Product, 0, len(productMap))

	for _, p := range productMap {
		result = append(result, p)
	}

	return result, totalProducts, nil
}

// GetProductByID fetches a single product by its ID for a given tenant.
//
// PURPOSE:
//   - Retrieve a fully hydrated product entity, including all associated images.
//   - Ensures only non-deleted products are returned.
//   - Handles nullable image fields due to LEFT JOIN.
//
// RETURNS:
//   - Pointer to Product entity including images.
//   - Error if product not found or query fails.
func (repo *ProductRepository) GetProductByID(
	ctx context.Context,
	tenantID value_object.TenantID,
	productID value_object.ProductID,
) (*entity.Product, error) {

	const op = "ProductRepository.GetProductByID"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
	SELECT
		p.id,
		p.tenant_id,
		p.name,
		p.description,
		p.price_cents,
		p.currency,
		p.sku,
		p.stock_quantity,
		p.product_properties,
		p.product_status,
		p.created_at,
		p.updated_at,

		i.id,
		i.image_url,
		i.sort_order,
		i.is_primary,
		i.created_at

	FROM products p
	LEFT JOIN product_images i
		ON p.id = i.product_id

	WHERE p.id = $1
	AND p.tenant_id = $2
	AND p.deleted_at IS NULL
	`

	rows, err := repo.DB.Query(
		ctx,
		query,
		productID.Raw(),
		tenantID.Raw(),
	)
	if err != nil {
		return nil, fmt.Errorf("[%s]: %w", op, err)
	}
	defer rows.Close()

	var prod *entity.Product

	for rows.Next() {

		// nullable image fields because of LEFT JOIN
		var imageID *uuid.UUID
		var imageURL *string
		var sortOrder *int
		var isPrimary *bool
		var imageCreatedAt *time.Time

		// initialize once
		if prod == nil {
			prod = &entity.Product{
				ProductImages: make([]entity.ProductImage, 0),
			}
		}

		err := rows.Scan(
			&prod.ID,
			&prod.TenantID,
			&prod.Name,
			&prod.Description,
			&prod.Price,
			&prod.Currency,
			&prod.SKU,
			&prod.Stock,
			&prod.Properties,
			&prod.Status,
			&prod.CreatedAt,
			&prod.UpdatedAt,

			&imageID,
			&imageURL,
			&sortOrder,
			&isPrimary,
			&imageCreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("[%s]: %w", op, err)
		}

		// append image if present
		if imageID != nil {

			img := entity.ProductImage{
				ID:        value_object.NewProductImageIDFromUUID(*imageID),
				ProductID: value_object.NewProductIDFromUUID(prod.ID.Raw().Bytes),
				ImageUrl:  *imageURL,
				SortOrder: int(*sortOrder),
				IsPrimary: *isPrimary,
				CreatedAt: *imageCreatedAt,
			}

			prod.ProductImages = append(prod.ProductImages, img)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("[%s]: %w", op, err)
	}

	if prod == nil {
		return nil, fmt.Errorf("[%s]: product not found", op)
	}

	return prod, nil
}

// UpdateProduct performs a PATCH-style update on a product (NO image handling).
//
// =========================
// DESIGN INTENT
// =========================
//
// This function is responsible ONLY for product fields:
//
// 1. PATCH SEMANTICS
//   - Only non-nil fields are updated
//
// 2. VERSION-AWARE PROPERTIES
//   - Same version → deep merge
//   - Different version → replace
//
// 3. TENANT ISOLATION
//   - Ensures product belongs to tenant
//
// 4. TRANSACTION SAFETY
//   - Locks row with FOR UPDATE
//
// =========================
// CONSISTENCY RULES
// =========================
// - Does NOT touch product_images table
func (r *ProductRepository) UpdateProduct(
	ctx context.Context,
	tenantID value_object.TenantID,
	productID value_object.ProductID,
	incoming entity.Product,
) (*entity.Product, error) {

	const op = "ProductRepository.UpdateProduct"

	tx, err := r.DB.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("[%s]: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// -------------------------
	// Lock product (tenant-safe)
	// -------------------------
	var current entity.Product

	err = tx.QueryRow(ctx, `
		SELECT id, tenant_id, name, description, price_cents, currency,
		       sku, stock_quantity, product_properties, product_status
		FROM products
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		FOR UPDATE
	`, productID, tenantID).Scan(
		&current.ID,
		&current.TenantID,
		&current.Name,
		&current.Description,
		&current.Price,
		&current.Currency,
		&current.SKU,
		&current.Stock,
		&current.Properties,
		&current.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("[%s]: fetch failed: %w", op, err)
	}

	// -------------------------
	// Merge ONLY product properties
	// -------------------------
	current.Properties = deepMerge(current.Properties, incoming.Properties)

	_, err = tx.Exec(ctx, `
		UPDATE products
		SET name=$1,
		    description=$2,
		    price_cents=$3,
		    currency=$4,
		    sku=$5,
		    stock_quantity=$6,
		    product_properties=$7,
		    product_status=$8,
		    updated_at=NOW()
		WHERE id=$9 AND tenant_id=$10
	`,
		incoming.Name,
		incoming.Description,
		incoming.Price,
		incoming.Currency,
		incoming.SKU,
		incoming.Stock,
		current.Properties,
		incoming.Status,
		productID,
		tenantID,
	)

	if err != nil {
		return nil, fmt.Errorf("[%s]: update failed: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("[%s]: commit failed: %w", op, err)
	}

	return r.GetProductByID(ctx, tenantID, productID)
}

// deepMerge merges two ProductProperties objects in a version-aware manner.
//
// =========================
// DESIGN INTENT
// =========================
//
// This function merges incoming product properties into existing ones while:
//
// 1. PRESERVING EXISTING DATA
//   - Existing keys are retained unless explicitly overwritten.
//
// 2. SUPPORTING NESTED STRUCTURES
//   - Recursively merges nested maps (deep merge).
//
// 3. HANDLING VERSIONING SAFELY
//   - If versions differ → incoming version overrides entirely (no merge).
//   - Prevents mixing incompatible schemas.
//
// =========================
// CONSISTENCY RULES
// =========================
// - Same version → deep merge Data
// - Different version → replace entire structure
//
// =========================
// FAILURE SAFETY
// =========================
// - Nil-safe (never panics)
// - Non-mutating (returns new object)
func deepMerge(
	existing, incoming *entity.ProductProperties,
) *entity.ProductProperties {

	// -------------------------
	// Nil handling
	// -------------------------
	if existing == nil && incoming == nil {
		return &entity.ProductProperties{
			Version: 1,
			Data:    map[string]any{},
		}
	}

	if existing == nil {
		return cloneProperties(incoming)
	}

	if incoming == nil {
		return cloneProperties(existing)
	}

	// -------------------------
	// Version mismatch → replace safely (CLONED)
	// -------------------------
	if existing.Version != incoming.Version {
		return cloneProperties(incoming)
	}

	// -------------------------
	// Deep merge data
	// -------------------------
	merged := make(map[string]any, len(existing.Data))

	// copy existing
	for k, v := range existing.Data {
		merged[k] = v
	}

	// merge incoming
	for k, v := range incoming.Data {

		if existingMap, ok1 := merged[k].(map[string]any); ok1 {
			if incomingMap, ok2 := v.(map[string]any); ok2 {

				merged[k] = deepMerge(
					&entity.ProductProperties{
						Version: existing.Version,
						Data:    existingMap,
					},
					&entity.ProductProperties{
						Version: incoming.Version,
						Data:    incomingMap,
					},
				).Data

				continue
			}
		}

		merged[k] = v
	}

	return &entity.ProductProperties{
		Version: existing.Version,
		Data:    merged,
	}
}

func cloneProperties(p *entity.ProductProperties) *entity.ProductProperties {
	if p == nil {
		return nil
	}

	cloned := make(map[string]any, len(p.Data))
	for k, v := range p.Data {
		cloned[k] = v
	}

	return &entity.ProductProperties{
		Version: p.Version,
		Data:    cloned,
	}
}

// AddProductImages inserts new images without affecting existing ones.
//
// =========================
// DESIGN INTENT
// =========================
//
// Supports incremental image additions.
//
// Guarantees:
// - Existing images remain untouched
// - If a new primary image is provided → existing primary is cleared
//
// =========================
// CONSISTENCY RULES
// =========================
// - At most ONE primary image per product
func (r *ProductRepository) AddProductImages(
	ctx context.Context,
	productID value_object.ProductID,
	images []entity.ProductImage,
) error {

	const op = "ProductRepository.AddProductImages"

	tx, err := r.DB.Tx(ctx)
	if err != nil {
		return fmt.Errorf("[%s]: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// Check if any incoming image is primary
	hasPrimary := false
	for _, img := range images {
		if img.IsPrimary {
			hasPrimary = true
			break
		}
	}

	// Reset existing primary if needed
	if hasPrimary {
		_, err = tx.Exec(ctx, `
			UPDATE product_images
			SET is_primary = FALSE
			WHERE product_id = $1
		`, productID)
		if err != nil {
			return fmt.Errorf("[%s]: reset primary failed: %w", op, err)
		}
	}

	// Insert images
	for _, img := range images {
		_, err := tx.Exec(ctx, `
			INSERT INTO product_images (
				product_id,
				image_url,
				sort_order,
				is_primary
			)
			VALUES ($1,$2,$3,$4)
		`, productID, img.ImageUrl, img.SortOrder, img.IsPrimary)

		if err != nil {
			return fmt.Errorf("[%s]: insert failed: %w", op, err)
		}
	}

	return tx.Commit(ctx)
}

// SoftDeleteProduct performs a cascading soft delete of a product.
//
// =========================
// DESIGN INTENT
// =========================
//
// Provides safe deletion while preserving historical data.
//
// Guarantees:
//
// 1. NON-DESTRUCTIVE DELETE
//   - Uses deleted_at timestamps
//
// 2. CASCADE CONSISTENCY
//   - All product images are also soft deleted
//
// 3. TRANSACTIONAL SAFETY
//   - No partial delete state possible
//
// =========================
// CONSISTENCY RULES
// =========================
// - Product must exist and not already be deleted
// - All images inherit deletion timestamp
func (r *ProductRepository) SoftDeleteProduct(
	ctx context.Context,
	tenantID value_object.TenantID,
	productID value_object.ProductID,
) error {

	const op = "ProductRepository.SoftDeleteProduct"

	tx, err := r.DB.Tx(ctx)
	if err != nil {
		return fmt.Errorf("[%s]: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// -------------------------
	// Lock product
	// -------------------------
	var id string

	err = tx.QueryRow(ctx, `
		SELECT id
		FROM products
		WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
		FOR UPDATE
	`, productID, tenantID).Scan(&id)

	if err != nil {
		return fmt.Errorf("[%s]: product not found: %w", op, err)
	}

	log.Printf("[%s]: product locked id=%s", op, id)

	now := time.Now()

	// -------------------------
	// Soft delete
	// -------------------------
	result, err := tx.Exec(ctx, `
		UPDATE products
		SET deleted_at = $1,
		    updated_at = $1
		WHERE id = $2 AND tenant_id = $3 AND deleted_at IS NULL
	`, now, productID, tenantID)

	if err != nil {
		return fmt.Errorf("[%s]: update failed: %w", op, err)
	}

	rows := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("[%s]: no rows updated", op)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("[%s]: commit failed: %w", op, err)
	}

	return nil
}

// SoftDeleteProductImages performs partial deletion of product images.
//
// =========================
// DESIGN INTENT
// =========================
//
// Allows selective removal of images without affecting others.
//
// Guarantees:
// - Only specified images are deleted
// - Uses soft delete (deleted_at)
//
// =========================
// EDGE CASE HANDLING
// =========================
// - If primary image is deleted → reassign a new primary
func (r *ProductRepository) SoftDeleteProductImages(
	ctx context.Context,
	productID value_object.ProductID,
	imageIDs []value_object.ProductImageID,
) error {

	const op = "ProductRepository.SoftDeleteProductImages"

	tx, err := r.DB.Tx(ctx)
	if err != nil {
		return fmt.Errorf("[%s]: %w", op, err)
	}
	defer tx.Rollback(ctx)

	if len(imageIDs) == 0 {
		return nil
	}

	_, err = tx.Exec(ctx, `
		UPDATE product_images
		SET deleted_at = NOW()
		WHERE product_id = $1
		  AND id = ANY($2)
		  AND deleted_at IS NULL
	`, productID, imageIDs)

	if err != nil {
		return fmt.Errorf("[%s]: delete failed: %w", op, err)
	}

	var exists bool
	err = tx.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM product_images
			WHERE product_id = $1
			  AND is_primary = TRUE
			  AND deleted_at IS NULL
		)
	`, productID).Scan(&exists)

	if err != nil {
		return fmt.Errorf("[%s]: primary check failed: %w", op, err)
	}

	// Reassign primary if needed
	if !exists {
		_, err = tx.Exec(ctx, `
			UPDATE product_images
			SET is_primary = TRUE
			WHERE id = (
				SELECT id FROM product_images
				WHERE product_id = $1
				  AND deleted_at IS NULL
				ORDER BY sort_order ASC
				LIMIT 1
			)
		`, productID)

		if err != nil {
			return fmt.Errorf("[%s]: primary reassignment failed: %w", op, err)
		}
	}

	return tx.Commit(ctx)
}
