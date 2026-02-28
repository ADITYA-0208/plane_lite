package api

import (
	"net/http"

	"planelite-backend/internal/task"
)

// RegisterTask registers task routes under workspaces/projects. Uses Auth + WorkspaceAccess.
func RegisterTask(mux *http.ServeMux, h *task.Handler, mw Middleware) {
	mux.Handle("POST /workspaces/{id}/projects/{pid}/tasks", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.Create))))
	mux.Handle("GET /workspaces/{id}/projects/{pid}/tasks", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.ListByProject))))
	mux.Handle("GET /workspaces/{id}/projects/{pid}/tasks/{tid}", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.GetByID))))
	mux.Handle("PATCH /workspaces/{id}/projects/{pid}/tasks/{tid}", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.Update))))
	mux.Handle("PUT /workspaces/{id}/projects/{pid}/tasks/{tid}", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.Update))))
}
