package errors

import "errors"

// Validation Errors
var (
	ErrMalformedJSON        = errors.New("malformed JSON")
	ErrInvalidInput         = errors.New("invalid input")
	ErrFullNameRequired     = errors.New("full name is required")
	ErrEmailRequired        = errors.New("email is required")
	ErrPhoneRequired        = errors.New("phone number is required")
	ErrPasswordRequired     = errors.New("password is required")
	ErrBusinessTypeRequired = errors.New("business type is required")
	ErrBusinessNameRequired = errors.New("business name is required")
	ErrIdIsRequired         = errors.New("id is required")
	ErrRoleIsRequired       = errors.New("role is required")
	ErrInvalidEmailFormat   = errors.New("invalid email format")
	ErrInvalidPhoneNumber   = errors.New("invalid phone number")
	ErrInvalidUUIDFormat    = errors.New("invalid UUID format")
	ErrInvalidID            = errors.New("invalid UUID")
	ErrInvalidPermissionID  = errors.New("invalid permission ID")
	ErrInvalidUserIdFormat  = errors.New("invalid user id format")
	ErrInvalidPageNumber    = errors.New("invalid page")
	ErrInvalidLimitNumber   = errors.New("invalid limit")
	ErrInvalidTenantName    = errors.New("tenant name cannot be empty")
)

// Database & Resource State Errors
var (
	ErrUserNotFound              = errors.New("user not found")
	ErrRoleNotFound              = errors.New("role not found")
	ErrTenantNotFound             = errors.New("tenant not found")
	ErrRoleAlreadyExists         = errors.New("role already exists")
	ErrTenantAlreadyExists       = errors.New("tenant with that name already exists")
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrUserEmailAlreadyExists    = errors.New("user with that email already exists")
	ErrUserMobileExists          = errors.New("user with that mobile already exists")
	ErrBusinessNameAlreadyExists = errors.New("business name already exists")
)

// Security & Credentials Errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountDeactivated = errors.New("your account has been deactivated. Please contact support.")
)

// Token-specific errors
var (
	ErrInvalidToken        = errors.New("invalid token")
	ErrExpiredToken        = errors.New("token expired")
	ErrNotValidYet         = errors.New("token not valid yet")
	ErrWrongIssuer         = errors.New("wrong issuer")
	ErrWrongAudience       = errors.New("wrong audience")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)
