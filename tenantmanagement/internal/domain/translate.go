package domain

import (
	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"
	"github.com/MartinMurithi/storeforge/tenantmanagement/internal/infrastructure/database/postgres"
)

// translateUserRepoError converts infra-level DB errors into domain-level errors.
func TranslateUserRepoError(err error) error {
	// Map infra error to a stable set of infra DB errors first
	switch postgres.MapPostgresError(err) {

	case postgres.ErrNotFound:
		return apperrors.ErrUserNotFound

	case postgres.ErrUniqueViolation:
		return apperrors.ErrUserAlreadyExists

	case postgres.ErrNotNull:
		return apperrors.ErrInvalidInput

	default:
		// Return original error if it doesn’t match known mappings
		return err
	}
}
