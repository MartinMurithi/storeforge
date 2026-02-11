package apperrors

import "errors"

// Transport / JSON errors
var (
	ErrMalformedJSON = errors.New("malformed JSON")
)

// Required field errors
var (
	ErrFullNameRequired     = errors.New("full name is required")
	ErrEmailRequired        = errors.New("email is required")
	ErrPhoneRequired        = errors.New("phone number is required")
	ErrPasswordRequired     = errors.New("password is required")
	ErrBusinessTypeRequired = errors.New("business type is required")
	ErrBusinessNameRequired = errors.New("business name is required")
	ErrIdIsRequired         = errors.New("id is required")
	ErrRoleIsRequired       = errors.New("role is required")
	ErrInvalidInput         = errors.New("invalid input")
)

// Format/validation errors
var (
	ErrInvalidEmailFormat  = errors.New("invalid email format")
	ErrInvalidPhoneNumber  = errors.New("invalid phone number")
	ErrInvalidUUIDFormat   = errors.New("invalid UUID format")
	ErrInvalidUserIdFormat = errors.New("invalid user id format")
)

// Business rule errors
var (
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrUserEmailAlreadyExists    = errors.New("user with that email already exists")
	ErrUserMobileExists          = errors.New("user with that mobile already exists")
	ErrBusinessNameAlreadyExists = errors.New("business name already exists")
	ErrAccountDeactivated        = errors.New("your account has been deactivated. Please contact support.")
)

// Credentials rule errors
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Database errors
var (
	ErrUserNotFound = errors.New("user not found")
)

// Pagination errors
var (
	ErrInvalidPageNumber  = errors.New("invalid page")
	ErrInvalidLimitNumber = errors.New("invalid limit")
)
