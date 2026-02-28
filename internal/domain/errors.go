package domain

import "errors"

var (
	ErrTitleRequired       = errors.New("ERR_INVALID_ARGUMENT: title is required")
	ErrInvalidMaterialType = errors.New("ERR_INVALID_ARGUMENT: type must be 'videos' or 'documents'")
	ErrEmailRequired       = errors.New("ERR_INVALID_ARGUMENT: email is required")
	ErrPasswordRequired    = errors.New("ERR_INVALID_ARGUMENT: password is required")
	ErrInvalidCredentials  = errors.New("ERR_UNAUTHORIZED: invalid credentials")
	ErrUserInactive        = errors.New("ERR_FORBIDDEN: user account is inactive")
	ErrUserAlreadyExists   = errors.New("ERR_CONFLICT: user with this email already exists")
	ErrNotFound            = errors.New("ERR_NOT_FOUND: not found")
	ErrInternalServerError = errors.New("ERR_INTERNAL_SERVER_ERROR: internal server error")
)
