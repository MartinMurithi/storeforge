package auth

import (
	"errors"

	"github.com/MartinMurithi/storeforge/usermanagement/internal/apperrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapAuthError(err error) error {
	switch {
	// Registration field validation
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
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
