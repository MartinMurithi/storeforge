package domain

import (
	"context"
	"fmt"
	// "log"
	"time"

	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/entity"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/domain/value_object"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/database"
)

type ThemeRepository struct {
	DB database.DB
}

type IThemeRepository interface {
	// Fetches the pre-saved "Golden Template" for a theme
	GetThemeById(ctx context.Context, themeId value_object.ThemeID) (*entity.Theme, error)
}

func NewThemeRepository(db database.DB) IThemeRepository {
	return &ThemeRepository{DB: db}
}

// func (r *ThemeRepository) GetThemeById(ctx context.Context, id value_object.ThemeID) (*entity.Theme, error) {
// 	const op = "ThemeRepository.GetThemeById"

// 	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
// 	defer cancel()

// 	query := `SELECT id, name, description, is_active, default_config, created_at FROM themes
// 	WHERE id = $1`

// 	theme := &entity.Theme{}

// 	err := r.DB.QueryRow(ctx, query, id).Scan(&theme.ID, &theme.Name, &theme.Description, &theme.IsActive, &theme.DefaultConfig, &theme.CreatedAt)
// 	if err != nil {
// 		log.Printf("[%s]: failed to fetch requested theme: %v", op, err)
// 		return nil, fmt.Errorf("[%s]: failed to fetch requested theme:: %w", op, err)
// 	}

// 	return theme, nil
// }

func (r *ThemeRepository) GetThemeById(ctx context.Context, id value_object.ThemeID) (*entity.Theme, error) {
	const op = "ThemeRepository.GetThemeById"
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `SELECT id, name, description, is_active, default_config, created_at FROM themes WHERE id = $1`

	// Initialize the pointers so they aren't nil during Scan
	theme := &entity.Theme{
		DefaultConfig: &entity.Settings{
			Config: make(entity.ThemeConfig),
		},
	}

	// Scan directly into the map field inside the struct
	err := r.DB.QueryRow(ctx, query, id).Scan(
		&theme.ID,
		&theme.Name,
		&theme.Description,
		&theme.IsActive,
		&theme.DefaultConfig.Config,
		&theme.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("[%s]: %w", op, err)
	}

	return theme, nil
}
