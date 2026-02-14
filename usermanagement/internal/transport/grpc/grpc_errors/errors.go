package grpc_errors

import (
	"errors"
	"log"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/apperrors"
	
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewMultiValidationError bundles multiple field errors into one gRPC status
func NewMultiValidationError(violations map[string]string) error {
	st := status.New(codes.InvalidArgument, "validation failed")

	br := &errdetails.BadRequest{}
	for field, desc := range violations {
		br.FieldViolations = append(br.FieldViolations, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: desc,
		})
	}

	st, _ = st.WithDetails(br)
	return st.Err()
}

func MapGrpcError(err error) error {

	if err == nil {
		return nil
	}

	// Check if the error is already a gRPC status (prevents double-wrapping)
	if _, ok := status.FromError(err); ok {
		return err
	}

	switch {
	case errors.Is(err, apperrors.ErrFullNameRequired):
		return status.Error(codes.InvalidArgument, "full_name is required")
	case errors.Is(err, apperrors.ErrEmailRequired):
		return status.Error(codes.InvalidArgument, "email is required")
	case errors.Is(err, apperrors.ErrPhoneRequired):
		return status.Error(codes.InvalidArgument, "phone is required")
	case errors.Is(err, apperrors.ErrPasswordRequired):
		return status.Error(codes.InvalidArgument, "password is required")
	case errors.Is(err, apperrors.ErrBusinessTypeRequired):
		return status.Error(codes.InvalidArgument, "business_type is required")
	case errors.Is(err, apperrors.ErrBusinessNameRequired):
		return status.Error(codes.InvalidArgument, "business_name is required")
	case errors.Is(err, apperrors.ErrInvalidEmailFormat):
		return status.Error(codes.InvalidArgument, "invalid email format")
	case errors.Is(err, apperrors.ErrInvalidPhoneNumber):
		return status.Error(codes.InvalidArgument, "invalid phone number")
	case errors.Is(err, apperrors.ErrUserAlreadyExists):
		return status.Error(codes.AlreadyExists, "user already exists")
	case errors.Is(err, apperrors.ErrUserEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, "user with that email already exists")
	case errors.Is(err, apperrors.ErrUserMobileExists):
		return status.Error(codes.AlreadyExists, "user with that mobile already exists")
	case errors.Is(err, apperrors.ErrInvalidUUIDFormat):
		return status.Error(codes.InvalidArgument, "invalid user Id")
	case errors.Is(err, apperrors.ErrInvalidPageNumber):
		return status.Error(codes.InvalidArgument, "invalid page number")
	case errors.Is(err, apperrors.ErrInvalidLimitNumber):
		return status.Error(codes.InvalidArgument, "invalid limit number")
	case errors.Is(err, apperrors.ErrAccountDeactivated):
		return status.Error(codes.PermissionDenied, "your account has been deactivated. please contact support")
	case errors.Is(err, apperrors.ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, apperrors.ErrInvalidCredentials):
		return status.Error(codes.PermissionDenied, "invalid credentials")
	default:
		log.Printf("ERROR: Unhandled application error: %v", err)
		return status.Error(codes.Internal, "internal server error")
	}
}
