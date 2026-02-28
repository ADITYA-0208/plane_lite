package middleware

import (
	"net/http"

	"planelite-backend/internal/common"
)

// RequireRole returns middleware that allows only the given roles.
func RequireRole(roles ...common.Role) func(http.Handler) http.Handler {
	allowed := make(map[common.Role]struct{})
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := common.GetContextUser(r.Context())
			if u == nil {
				common.Error(w, common.ErrUnauthorized)
				return
			}
			if _, ok := allowed[u.Role]; !ok {
				common.Error(w, common.ErrForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
