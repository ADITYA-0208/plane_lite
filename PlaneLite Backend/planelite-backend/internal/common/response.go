package common

import (
	"encoding/json"
	"errors"
	"net/http"
)

// JSON writes a JSON response with status code.
func JSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

// Success is the standard success envelope.
type Success struct {
	Data any `json:"data,omitempty"`
}

// ErrorResp is the standard error envelope.
type ErrorResp struct {
	Error string `json:"error"`
}

// OK writes 200 with data.
func OK(w http.ResponseWriter, data any) {
	JSON(w, http.StatusOK, Success{Data: data})
}

// Created writes 201 with data.
func Created(w http.ResponseWriter, data any) {
	JSON(w, http.StatusCreated, Success{Data: data})
}

// NoContent writes 204.
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Error writes JSON error and status from common errors.
func Error(w http.ResponseWriter, err error) {
	status := StatusFromError(err)
	msg := err.Error()
	if msg == "" {
		msg = http.StatusText(status)
	}
	JSON(w, status, ErrorResp{Error: msg})
}

// StatusFromError maps sentinel/common errors to HTTP status.
func StatusFromError(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrBadRequest), errors.Is(err, ErrInvalidInput), errors.Is(err, ErrConflict):
		return http.StatusBadRequest
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrConflict), errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
