package domain

import (
	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/MartinMurithi/storeforge/productmanagement/internal/infrastructure/database/postgres"
)

// translateTenantRepoError converts infra-level DB errors into domain-level errors.
func TranslateProductRepoError(err error) error {
	// Map infra error to a stable set of infra DB errors first
	switch postgres.MapPostgresError(err) {

	case postgres.ErrNotFound:
		return apperrors.ErrProductNotFound

	case postgres.ErrUniqueViolation:
		return apperrors.ErrProductAlreadyExists

	case postgres.ErrNotNull:
		return apperrors.ErrInvalidInput

	default:
		// Return original error if it doesn’t match known mappings
		return err
	}
}
