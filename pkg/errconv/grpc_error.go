package errconv

import (
	"errors"
	"log"

	apperrors "github.com/MartinMurithi/storeforge/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToGrpcError converts internal domain errors into gRPC status errors.
func ToGrpcError(err error) error {
	if err == nil {
		return nil
	}

	// MapAppErrorToGrpc converts internal domain errors into gRPC status errors.
	// This is used by the gRPC service layer before returning errors to the caller.
	if _, ok := status.FromError(err); ok {
		return err
	}

	switch {
	// 400 - Invalid Argument
	case errors.Is(err, apperrors.ErrFullNameRequired),
		errors.Is(err, apperrors.ErrEmailRequired),
		errors.Is(err, apperrors.ErrPhoneRequired),
		errors.Is(err, apperrors.ErrPasswordRequired),
		errors.Is(err, apperrors.ErrBusinessTypeRequired),
		errors.Is(err, apperrors.ErrBusinessNameRequired),
		errors.Is(err, apperrors.ErrInvalidEmailFormat),
		errors.Is(err, apperrors.ErrInvalidPhoneNumber),
		errors.Is(err, apperrors.ErrInvalidUUIDFormat),
		errors.Is(err, apperrors.ErrInvalidPageNumber),
		errors.Is(err, apperrors.ErrInvalidLimitNumber),
		errors.Is(err, apperrors.ErrInvalidPermissionID):
		return status.Error(codes.InvalidArgument, err.Error())

	// 401 - Unauthenticated
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, err.Error())

	// 403 - Permission Denied
	case errors.Is(err, apperrors.ErrAccountDeactivated):
		return status.Error(codes.PermissionDenied, err.Error())

	// 404 - Not Found
	case errors.Is(err, apperrors.ErrUserNotFound):
		return status.Error(codes.NotFound, err.Error())

	// 409 - Conflict
	case errors.Is(err, apperrors.ErrUserAlreadyExists),
		errors.Is(err, apperrors.ErrRoleAlreadyExists),
		errors.Is(err, apperrors.ErrUserEmailAlreadyExists),
		errors.Is(err, apperrors.ErrTenantAlreadyExists),
		errors.Is(err, apperrors.ErrUserMobileExists),
		errors.Is(err, apperrors.ErrBusinessNameAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())

	// 500 - Internal
	default:
		// Log the actual error for debugging, but hide details from the client
		log.Printf("unhandled_internal_error: %v", err)
		return status.Error(codes.Internal, "an internal server error occurred")
	}
}
