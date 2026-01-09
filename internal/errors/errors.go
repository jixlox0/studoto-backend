package errors

import "errors"

// User-related errors
var (
	ErrUserNotFound       = errors.New("user_not_found")
	ErrUserAlreadyExists  = errors.New("user_already_exists")
	ErrInvalidPassword    = errors.New("invalid_password")
	ErrInvalidProvider    = errors.New("invalid_provider")
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrInvalidToken       = errors.New("invalid_token")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrBadRequest         = errors.New("bad_request")
	ErrInternalError      = errors.New("internal_error")
	ErrValidationError    = errors.New("validation_error")
	ErrXAuthKeyRequired   = errors.New("x_auth_key_required")
	ErrXAuthKeyEmpty      = errors.New("x_auth_key_empty")
	ErrCodeRequired       = errors.New("code_required")
	ErrNotAuthenticated   = errors.New("not_authenticated")
	ErrInvalidUserIDType  = errors.New("invalid_user_id_type")
)
