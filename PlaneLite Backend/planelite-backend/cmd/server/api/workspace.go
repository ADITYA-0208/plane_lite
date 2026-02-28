package api

import (
	"net/http"

	"planelite-backend/internal/workspace"
)

// RegisterWorkspace registers workspace routes. Uses Auth, AdminOnly, WorkspaceAccess.
func RegisterWorkspace(mux *http.ServeMux, h *workspace.Handler, mw Middleware) {
	mux.Handle("POST /workspaces", mw.Auth(mw.AdminOnly(http.HandlerFunc(h.Create))))
	mux.Handle("GET /workspaces", mw.Auth(mw.AdminOnly(http.HandlerFunc(h.ListMine))))
	mux.Handle("GET /workspaces/{id}", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.GetByID))))
	mux.Handle("POST /workspaces/{id}/members", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.AddMember))))
	mux.Handle("GET /workspaces/{id}/members", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.ListMembers))))
	mux.Handle("POST /workspaces/{id}/members/{mid}/approve", mw.Auth(mw.WorkspaceAccess(http.HandlerFunc(h.ApproveMember))))
}
