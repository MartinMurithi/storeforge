package repository

import (
	"github.com/MartinMurithi/storeforge/usermanagement/internal/apperrors"
	"github.com/MartinMurithi/storeforge/usermanagement/internal/database"
)

// translateUserRepoError converts infra-level DB errors into domain-level errors.
func TranslateUserRepoError(err error) error {
	// Map infra error to a stable set of infra DB errors first
	switch database.MapPostgresError(err) {

	case database.ErrNotFound:
		return apperrors.ErrUserNotFound

	case database.ErrUniqueViolation:
		return apperrors.ErrUserAlreadyExists

	case database.ErrNotNull:
		return apperrors.ErrInvalidInput

	default:
		// Return original error if it doesn’t match known mappings
		return err
	}
}
