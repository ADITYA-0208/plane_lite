package api

import "net/http"

// Middleware holds shared middleware used by route packages.
// Filled in main and passed to each Register* function.
type Middleware struct {
	Auth            func(http.Handler) http.Handler
	AdminOnly       func(http.Handler) http.Handler
	WorkspaceAccess func(http.Handler) http.Handler
}
