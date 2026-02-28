package middleware

import (
	"net/http"
	"strings"

	"planelite-backend/internal/auth"
	"planelite-backend/internal/common"
)

// Auth validates JWT and sets claims in context. Use ContextClaims to read.
func Auth(authSvc *auth.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				common.Error(w, common.ErrUnauthorized)
				return
			}
			const prefix = "Bearer "
			if !strings.HasPrefix(header, prefix) {
				common.Error(w, common.ErrUnauthorized)
				return
			}
			token := strings.TrimPrefix(header, prefix)
			claims, err := authSvc.ValidateToken(token)
			if err != nil {
				common.Error(w, common.ErrUnauthorized)
				return
			}
			ctx := common.WithContextUser(r.Context(), common.ContextUser{UserID: claims.UserID, Role: claims.Role})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
