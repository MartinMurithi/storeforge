package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/domain/products/entity"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/database"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/database/postgres"
)

type ProductRepository struct {
	DB database.DB
}

type IProductRepository interface {
	CreateProduct(ctx context.Context, product *entity.Product) error
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
			price,
			sku,
			stock_quantity,
			product_properties,
			product_status
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, name, created_at
	`

	err = tx.QueryRow(
		ctx,
		productQuery,
		product.TenantID,
		product.Name,
		product.Description,
		product.Price,
		product.SKU,
		product.Stock,
		product.Properties,
		product.Status,
	).Scan(&product.ID, &product.Name, &product.CreatedAt)

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
		`

		for _, img := range product.ProductImages {
			_, err = tx.Exec(ctx, imageQuery, product.ID, img.ImageUrl, img.SortOrder, img.IsPrimary)
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
