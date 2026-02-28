package api

import (
	"net/http"

	"planelite-backend/internal/project"
)

// RegisterProject registers project routes under workspaces. Uses Auth + WorkspaceAccess.
func RegisterProject(mux *http.ServeMux, h *project.Handler, mw Middleware) {
	mux.Handle("POST /workspaces/{id}/projects", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.Create))))
	mux.Handle("GET /workspaces/{id}/projects", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.ListByWorkspace))))
	mux.Handle("GET /workspaces/{id}/projects/{pid}", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.GetByID))))
}
