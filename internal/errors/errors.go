package errors

import "errors"

// User-related errors
var (
	ErrUserNotFound       = errors.New("User not found")
	ErrUserAlreadyExists  = errors.New("User already exists")
	ErrInvalidPassword    = errors.New("Invalid password")
	ErrInvalidProvider    = errors.New("Invalid provider")
	ErrInvalidCredentials = errors.New("Invalid credentials")
	ErrInvalidToken       = errors.New("Invalid token")
	ErrUnauthorized       = errors.New("Unauthorized")
	ErrBadRequest         = errors.New("Bad request")
	ErrInternalError      = errors.New("Internal error")
	ErrValidationError    = errors.New("Validation error")
	ErrXAuthKeyRequired   = errors.New("X-Auth-Key required")
	ErrXAuthKeyEmpty      = errors.New("X-Auth-Key empty")
	ErrCodeRequired       = errors.New("Code required")
	ErrNotAuthenticated   = errors.New("Not authenticated")
	ErrInvalidUserIDType  = errors.New("Invalid user ID type")
)
