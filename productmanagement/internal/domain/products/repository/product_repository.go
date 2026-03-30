package repository

import (
	"context"
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

func (repo *ProductRepository) CreateProduct(ctx context.Context, product *entity.Product) error {
	const op = "product_repository.CreateProduct"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, name
	`

	err := repo.DB.QueryRow(
		ctx,
		query,
		product.TenantID,
		product.Name,
		product.Description,
		product.Price,
		product.SKU,
		product.Stock,
		product.Properties,
		product.Status,
	).Scan(&product.ID, &product.Name)

	if err != nil {
		log.Printf("[%s]: error creating product: %v", op, err)
		infraErr := postgres.MapPostgresError(err)
		return domain.TranslateProductRepoError(infraErr)
	}

	return nil
}
