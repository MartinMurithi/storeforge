package repository

import (
	"context"
	"fmt"
	"log"
	"strings"
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
}

func NewProductRepository(db database.DB) IProductRepository {
	return &ProductRepository{DB: db}
}

/*
CreateProduct persists a product and its associated images atomically.

Image handling:
- 0..n images supported
- sort_order is assigned automatically based on array index
- first image in the array is automatically set as primary (is_primary = true)
- ensures transactional consistency: either product + all images are saved, or nothing is saved
*/
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
		product.TenantID,    // $1
		product.Name,        // $2
		product.Description, // $3
		product.Price,       // $4 -> price_cents
		product.Currency,    // $5
		product.SKU,         // $6
		product.Stock,       // $7 -> stock_quantity
		product.Properties,  // $8 -> jsonb
		product.Status,      // $9 -> enum
	).Scan(&id, &product.Name, &product.Currency, &product.CreatedAt)

	product.ID = value_object.NewProductIDFromUUID(id)

	if err != nil {
		log.Printf("[%s error]: %v", op, err)
		return fmt.Errorf(
			"%w",
			domain.TranslateProductRepoError(postgres.MapPostgresError(err)),
		)
	}

	// Normalize images: assign sort_order and set first image as primary
	if len(product.ProductImages) > 0 {
		for i := range product.ProductImages {
			product.ProductImages[i].SortOrder = i
			product.ProductImages[i].IsPrimary = false
		}
		product.ProductImages[0].IsPrimary = true
	}

	// Insert images
	if len(product.ProductImages) > 0 {
		imageQuery := `
			INSERT INTO product_images (
				product_id,
				image_url,
				sort_order,
				is_primary
			)
			VALUES ($1,$2,$3,$4)
			RETURNING id, product_id, created_at, sort_order, is_primary
		`

		for _, img := range product.ProductImages {

			var dbImgID uuid.UUID
			var dbProdID uuid.UUID

			img.ID = value_object.NewProductImageIDFromUUID(dbImgID)
			img.ProductID = value_object.NewProductIDFromUUID(dbProdID)

			err = tx.QueryRow(ctx, imageQuery, product.ID.Raw(), img.ImageUrl, img.SortOrder, img.IsPrimary).Scan(&img.ID, &img.ProductID, &img.CreatedAt, &img.SortOrder, &img.IsPrimary)
			if err != nil {
				log.Printf("[%s error inserting image]: %v", op, err)
				return fmt.Errorf(
					"%w",
					domain.TranslateProductRepoError(postgres.MapPostgresError(err)),
				)
			}
		}
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

// UpdateProductWithImages performs an atomic update of a product and its images.
//
// PURPOSE:
//   - Allows partial updates to product fields without overwriting unspecified fields.
//   - Performs deep merge for ProductProperties JSONB.
//   - Inserts, updates, and soft-deletes associated product images.
//   - Maintains a single primary image per product.
//   - Ensures transactional safety and consistency.
//
// PARAMETERS:
//   - incoming: fields to update (pointers indicate optional updates).
//   - incomingImages: list of images to upsert; missing images are soft-deleted.
//
// RETURNS:
//   - Fully hydrated Product entity with updated images.
//   - Error if any database operation fails.
func (r *ProductRepository) UpdateProductWithImages(
	ctx context.Context,
	tenantID value_object.TenantID,
	productID value_object.ProductID,
	incoming entity.ProductUpdate,
	incomingImages []entity.ProductImage,
) (*entity.Product, error) {

	const op = "ProductRepository.UpdateProductWithImages"

	tx, err := r.DB.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// Fetch current product
	var current entity.Product
	err = tx.QueryRow(ctx, `
        SELECT id, tenant_id, name, description, price_cents, currency, sku, stock_quantity, product_properties, product_status
        FROM products
        WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
        FOR UPDATE
    `, productID, tenantID).Scan(
		&current.ID, &current.TenantID, &current.Name, &current.Description,
		&current.Price, &current.Currency, &current.SKU, &current.Stock,
		&current.Properties, &current.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to fetch current product: %w", op, err)
	}

	if current.Properties == nil {
		current.Properties = &entity.ProductProperties{}
	}
	if incoming.Properties == nil {
		incoming.Properties = &entity.ProductProperties{}
	}

	// Deep merge properties
	mergedProps := deepMerge(*current.Properties, *incoming.Properties)

	// Apply other updates if provided
	if incoming.Name != nil {
		current.Name = *incoming.Name
	}
	if incoming.Description != nil {
		current.Description = *incoming.Description
	}
	if incoming.PriceCents != nil {
		current.Price = *incoming.PriceCents
	}
	if incoming.Currency != nil {
		current.Currency = *incoming.Currency
	}
	if incoming.SKU != nil {
		current.SKU = *incoming.SKU
	}
	if incoming.StockQuantity != nil {
		current.Stock = *incoming.StockQuantity
	}
	if incoming.Status != nil {
		current.Status = *incoming.Status
	}

	// Update product row
	_, err = tx.Exec(ctx, `
        UPDATE products
        SET name=$1, description=$2, price_cents=$3, currency=$4,
            sku=$5, stock_quantity=$6, product_properties=$7,
            product_status=$8, updated_at=NOW()
        WHERE id=$9
    `, current.Name, current.Description, current.Price, current.Currency, current.SKU,
		current.Stock, mergedProps, current.Status, current.ID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to update product: %w", op, err)
	}

	// If any incoming image is primary, unset all others
	primaryExists := false
	for _, img := range incomingImages {
		if img.IsPrimary {
			primaryExists = true
			break
		}
	}
	if primaryExists {
		_, err := tx.Exec(ctx, `
            UPDATE product_images
            SET is_primary = FALSE
            WHERE product_id = $1
        `, current.ID)
		if err != nil {
			return nil, fmt.Errorf("[%s]: failed to unset existing primary images: %w", op, err)
		}
	}

	// Insert/update incoming images
	for _, img := range incomingImages {
		if img.ID.String() == "" {
			_, err := tx.Exec(ctx, `
                INSERT INTO product_images (product_id, image_url, sort_order, is_primary)
                VALUES ($1, $2, $3, $4)
            `, current.ID, img.ImageUrl, img.SortOrder, img.IsPrimary)
			if err != nil {
				return nil, fmt.Errorf("[%s]: failed to insert product image: %w", op, err)
			}
		} else {
			_, err := tx.Exec(ctx, `
                UPDATE product_images
                SET image_url=$1, sort_order=$2, is_primary=$3
                WHERE id=$4 AND product_id=$5
            `, img.ImageUrl, img.SortOrder, img.IsPrimary, img.ID, current.ID)
			if err != nil {
				return nil, fmt.Errorf("[%s]: failed to update product image: %w", op, err)
			}
		}
	}

	// Soft-delete images not in incomingImages
	imageIDs := make([]interface{}, 0, len(incomingImages))
	for _, img := range incomingImages {
		if img.ID != "" {
			imageIDs = append(imageIDs, img.ID)
		}
	}
	if len(imageIDs) > 0 {
		query := fmt.Sprintf(`
            UPDATE product_images
            SET deleted_at = NOW()
            WHERE product_id = $1 AND id NOT IN (%s)
        `, placeholders(len(imageIDs), 2))
		args := append([]interface{}{current.ID}, imageIDs...)
		if _, err := tx.Exec(ctx, query, args...); err != nil {
			return nil, fmt.Errorf("[%s]: failed to soft-delete removed images: %w", op, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("[%s]: failed to commit tx: %w", op, err)
	}

	// Return full product with images
	return r.GetProductByID(ctx, current.TenantID, current.ID)
}

// deepMerge recursively merges two ProductProperties maps.
//
// PURPOSE:
//   - Used to merge incoming ProductProperties into existing product properties.
//   - Preserves existing nested keys not provided in incoming map.
//
// RETURNS:
//   - A new ProductProperties map containing the merged result.
func deepMerge(existing, incoming entity.ProductProperties) entity.ProductProperties {
	merged := make(entity.ProductProperties)
	for k, v := range existing {
		merged[k] = v
	}
	for k, v := range incoming {
		if existingMap, ok1 := merged[k].(map[string]any); ok1 {
			if incomingMap, ok2 := v.(map[string]any); ok2 {
				merged[k] = deepMerge(existingMap, incomingMap)
				continue
			}
		}
		merged[k] = v
	}
	return merged
}

// placeholders generates SQL placeholders for an IN clause.
//
// PURPOSE:
//   - Generates "$2,$3,$4..." style placeholders for variable-length SQL queries.
//   - Useful when building dynamic IN clauses in queries like soft-deleting images.
//
// PARAMETERS:
//   - n: number of placeholders to generate.
//   - startIdx: starting index for the first placeholder.
//
// RETURNS:
//   - A comma-separated string of placeholders, e.g., "$2,$3,$4".
func placeholders(n int, startIdx int) string {
	parts := make([]string, n)
	for i := 0; i < n; i++ {
		parts[i] = fmt.Sprintf("$%d", startIdx+i)
	}
	return strings.Join(parts, ",")
}

// SoftDeleteProduct performs an atomic soft-delete of a product and its images.
//
// PURPOSE:
//   - Marks a product and all associated images as deleted without removing data.
//   - Ensures transactional safety and consistency.
//
// HOW IT WORKS:
//  1. Begins a transaction.
//  2. Locks the product row using FOR UPDATE.
//  3. Sets deleted_at for the product.
//  4. Sets deleted_at for all associated images.
//  5. Commits the transaction.
//  6. Returns the full product context (with deleted_at timestamps populated).
func (r *ProductRepository) SoftDeleteProduct(
	ctx context.Context,
	tenantID value_object.TenantID,
	productID value_object.ProductID,
) (*entity.Product, error) {

	const op = "ProductRepository.SoftDeleteProduct"

	tx, err := r.DB.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to begin tx: %w", op, err)
	}
	defer tx.Rollback(ctx)

	// 1️⃣ Lock the product
	var p entity.Product
	err = tx.QueryRow(ctx, `
        SELECT id, tenant_id, name, description, price_cents, currency, sku,
               stock_quantity, product_properties, product_status, created_at, updated_at, deleted_at
        FROM products
        WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL
        FOR UPDATE
    `, productID, tenantID).Scan(
		&p.ID, &p.TenantID, &p.Name, &p.Description, &p.Price, &p.Currency, &p.SKU,
		&p.Stock, &p.Properties, &p.Status, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to fetch product for deletion: %w", op, err)
	}

	now := time.Now()

	// 2️⃣ Soft-delete product
	_, err = tx.Exec(ctx, `
        UPDATE products
        SET deleted_at = $1, updated_at = $1
        WHERE id = $2
    `, now, p.ID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to soft-delete product: %w", op, err)
	}

	// 3️⃣ Soft-delete all associated images
	_, err = tx.Exec(ctx, `
        UPDATE product_images
        SET deleted_at = $1
        WHERE product_id = $2 AND deleted_at IS NULL
    `, now, p.ID)
	if err != nil {
		return nil, fmt.Errorf("[%s]: failed to soft-delete product images: %w", op, err)
	}

	// 4️⃣ Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("[%s]: failed to commit tx: %w", op, err)
	}

	// 5️⃣ Return the product context (deleted_at now populated)
	return r.GetProductByID(ctx, p.TenantID, p.ID)
}
