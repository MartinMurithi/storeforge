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
)

// Format/validation errors
var (
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
)

// Business rule errors
var (
	ErrUserAlreadyExists         = errors.New("user with that email already exists")
	ErrBusinessNameAlreadyExists = errors.New("business name already exists")
)

// Credentials rule errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
)

// Database errors
var (
	ErrUserNotFound = errors.New("user not found")
)
