package api

import (
	"net/http"

	"planelite-backend/internal/auth"
)

// RegisterAuth registers auth routes (signup, login). No middleware.
func RegisterAuth(mux *http.ServeMux, h *auth.Handler) {
	mux.HandleFunc("POST /auth/signup", h.Signup)
	mux.HandleFunc("POST /auth/login", h.Login)
}
