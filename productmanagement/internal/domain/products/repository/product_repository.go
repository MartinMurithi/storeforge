package repository

import (
	"context"
	"fmt"
	"log"
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
