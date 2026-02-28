package common

import "errors"

// Sentinel errors for handlers to map to HTTP status.
var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrBadRequest    = errors.New("bad request")
	ErrConflict      = errors.New("conflict (e.g. duplicate)")
	ErrInvalidInput  = errors.New("invalid input")
)
