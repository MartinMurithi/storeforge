package grpc_errors

import (
	"errors"
	"log"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/apperrors"
	
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


// func MapGrpcError(err error) error {

// 	if err == nil {
// 		return nil
// 	}

// 	// Check if the error is already a gRPC status (prevents double-wrapping)
// 	if _, ok := status.FromError(err); ok {
// 		return err
// 	}

// 	switch {
// 	case errors.Is(err, apperrors.ErrFullNameRequired):
// 		return status.Error(codes.InvalidArgument, "full_name is required")
// 	case errors.Is(err, apperrors.ErrEmailRequired):
// 		return status.Error(codes.InvalidArgument, "email is required")
// 	case errors.Is(err, apperrors.ErrPhoneRequired):
// 		return status.Error(codes.InvalidArgument, "phone is required")
// 	case errors.Is(err, apperrors.ErrPasswordRequired):
// 		return status.Error(codes.InvalidArgument, "password is required")
// 	case errors.Is(err, apperrors.ErrBusinessTypeRequired):
// 		return status.Error(codes.InvalidArgument, "business_type is required")
// 	case errors.Is(err, apperrors.ErrBusinessNameRequired):
// 		return status.Error(codes.InvalidArgument, "business_name is required")
// 	case errors.Is(err, apperrors.ErrInvalidEmailFormat):
// 		return status.Error(codes.InvalidArgument, "invalid email format")
// 	case errors.Is(err, apperrors.ErrInvalidPhoneNumber):
// 		return status.Error(codes.InvalidArgument, "invalid phone number")
// 	case errors.Is(err, apperrors.ErrUserAlreadyExists):
// 		return status.Error(codes.AlreadyExists, "user already exists")
// 	case errors.Is(err, apperrors.ErrUserEmailAlreadyExists):
// 		return status.Error(codes.AlreadyExists, "user with that email already exists")
// 	case errors.Is(err, apperrors.ErrUserMobileExists):
// 		return status.Error(codes.AlreadyExists, "user with that mobile already exists")
// 	case errors.Is(err, apperrors.ErrInvalidUUIDFormat):
// 		return status.Error(codes.InvalidArgument, "invalid user Id")
// 	case errors.Is(err, apperrors.ErrInvalidPageNumber):
// 		return status.Error(codes.InvalidArgument, "invalid page number")
// 	case errors.Is(err, apperrors.ErrInvalidLimitNumber):
// 		return status.Error(codes.InvalidArgument, "invalid limit number")
// 	case errors.Is(err, apperrors.ErrAccountDeactivated):
// 		return status.Error(codes.PermissionDenied, "your account has been deactivated. please contact support")
// 	case errors.Is(err, apperrors.ErrUserNotFound):
// 		return status.Error(codes.NotFound, "user not found")
// 	case errors.Is(err, apperrors.ErrInvalidCredentials):
// 		return status.Error(codes.PermissionDenied, "invalid credentials")
// 	default:
// 		log.Printf("ERROR: Unhandled application error: %v", err)
// 		return status.Error(codes.Internal, "internal server error")
// 	}
// }


// MapAppErrorToGrpc converts internal domain errors into gRPC status errors.
// This is used by the gRPC service layer before returning errors to the caller.
func MapAppErrorToGrpc(err error) error {
    if err == nil {
        return nil
    }

    // If it's already a gRPC error (from another service), return as-is
    if _, ok := status.FromError(err); ok {
        return err
    }

    // Grouping by gRPC Code for readability
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
         errors.Is(err, apperrors.ErrInvalidLimitNumber):
        return status.Error(codes.InvalidArgument, err.Error())

    // 404 - Not Found
    case errors.Is(err, apperrors.ErrUserNotFound):
        return status.Error(codes.NotFound, err.Error())

    // 409 - Conflict
    case errors.Is(err, apperrors.ErrUserAlreadyExists),
         errors.Is(err, apperrors.ErrUserEmailAlreadyExists),
         errors.Is(err, apperrors.ErrUserMobileExists),
         errors.Is(err, apperrors.ErrBusinessNameAlreadyExists):
        return status.Error(codes.AlreadyExists, err.Error())

    // 403 - Permission Denied
    case errors.Is(err, apperrors.ErrAccountDeactivated),
         errors.Is(err, apperrors.ErrInvalidCredentials):
        return status.Error(codes.PermissionDenied, err.Error())

    // 500 - Internal
    default:
        // Log the actual error for debugging, but hide details from the client
        log.Printf("unhandled_error: %v", err)
        return status.Error(codes.Internal, "an internal server error occurred")
    }
}