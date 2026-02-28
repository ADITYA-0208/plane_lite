package api

import (
	"net/http"

	"planelite-backend/internal/user"
)

// RegisterUser registers user routes (me, get by id). Uses Auth middleware.
func RegisterUser(mux *http.ServeMux, h *user.Handler, mw Middleware) {
	mux.Handle("GET /me", mw.Auth(http.HandlerFunc(h.GetMe)))
	mux.Handle("GET /users/{id}", mw.Auth(http.HandlerFunc(h.GetByID)))
}
