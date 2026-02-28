package middleware

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"planelite-backend/internal/common"
	"planelite-backend/internal/workspace"
)

// WorkspaceAccess checks that the current user has approved access to the workspace
// identified by URL path or query (e.g. workspace_id). Use after Auth middleware.
// getWorkspaceID extracts workspace ID from request (path param or query).
type WorkspaceAccess struct {
	Membership *workspace.Service
	// GetWorkspaceID returns workspace ID from request; e.g. from path "GET /workspaces/:id"
	GetWorkspaceID func(*http.Request) (primitive.ObjectID, bool)
}

// Middleware returns a middleware that calls Check and returns 403 if no access.
func (w *WorkspaceAccess) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			u := common.GetContextUser(r.Context())
			if u == nil || u.UserID == "" {
				common.Error(rw, common.ErrUnauthorized)
				return
			}
			wsID, ok := w.GetWorkspaceID(r)
			if !ok {
				common.Error(rw, common.ErrBadRequest)
				return
			}
			userID, err := primitive.ObjectIDFromHex(u.UserID)
			if err != nil {
				common.Error(rw, common.ErrUnauthorized)
				return
			}
			allowed, err := w.Membership.HasApprovedAccess(r.Context(), userID, wsID)
			if err != nil || !allowed {
				common.Error(rw, common.ErrForbidden)
				return
			}
			next.ServeHTTP(rw, r)
		})
	}
}
