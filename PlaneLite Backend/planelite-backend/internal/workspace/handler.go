package workspace

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"planelite-backend/internal/common"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// CreateRequest is the JSON body for POST /workspaces.
type CreateRequest struct {
	Name string `json:"name"`
}

// Create handles POST /workspaces. Caller must be ADMIN (enforced by route).
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.Error(w, common.ErrBadRequest)
		return
	}
	u := common.GetContextUser(r.Context())
	if u == nil || u.UserID == "" {
		common.Error(w, common.ErrUnauthorized)
		return
	}
	adminID, err := primitive.ObjectIDFromHex(u.UserID)
	if err != nil {
		common.Error(w, common.ErrUnauthorized)
		return
	}
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	ws, err := h.svc.Create(r.Context(), adminID, req.Name)
	if err != nil {
		common.Error(w, err)
		return
	}
	common.Created(w, ws)
}

// GetByID handles GET /workspaces/:id.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.Error(w, common.ErrBadRequest)
		return
	}
	idHex := r.PathValue("id")
	if idHex == "" {
		common.Error(w, common.ErrBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	ws, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		common.Error(w, common.ErrNotFound)
		return
	}
	common.OK(w, ws)
}

// ListMine handles GET /workspaces (workspaces I admin).
func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.Error(w, common.ErrBadRequest)
		return
	}
	u := common.GetContextUser(r.Context())
	if u == nil || u.UserID == "" {
		common.Error(w, common.ErrUnauthorized)
		return
	}
	adminID, err := primitive.ObjectIDFromHex(u.UserID)
	if err != nil {
		common.Error(w, common.ErrUnauthorized)
		return
	}
	list, err := h.svc.ListByAdminID(r.Context(), adminID)
	if err != nil {
		common.Error(w, err)
		return
	}
	if list == nil {
		list = []*Workspace{}
	}
	common.OK(w, list)
}

// AddMemberRequest is the JSON body for POST /workspaces/:id/members.
type AddMemberRequest struct {
	UserID string `json:"user_id"`
}

// AddMember handles POST /workspaces/:id/members. Caller must be workspace admin.
func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.Error(w, common.ErrBadRequest)
		return
	}
	wsID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	var req AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" {
		common.Error(w, common.ErrBadRequest)
		return
	}
	userID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	m, err := h.svc.AddMember(r.Context(), wsID, userID)
	if err != nil {
		common.Error(w, err)
		return
	}
	common.Created(w, m)
}

// ApproveMember handles POST /workspaces/:id/members/:mid/approve.
func (h *Handler) ApproveMember(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.Error(w, common.ErrBadRequest)
		return
	}
	u := common.GetContextUser(r.Context())
	if u == nil || u.UserID == "" {
		common.Error(w, common.ErrUnauthorized)
		return
	}
	adminID, err := primitive.ObjectIDFromHex(u.UserID)
	if err != nil {
		common.Error(w, common.ErrUnauthorized)
		return
	}
	midHex := r.PathValue("mid")
	if midHex == "" {
		common.Error(w, common.ErrBadRequest)
		return
	}
	membershipID, err := primitive.ObjectIDFromHex(midHex)
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	if err := h.svc.ApproveMember(r.Context(), adminID, membershipID); err != nil {
		common.Error(w, err)
		return
	}
	common.NoContent(w)
}

// ListMembers handles GET /workspaces/:id/members.
func (h *Handler) ListMembers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.Error(w, common.ErrBadRequest)
		return
	}
	wsID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	list, err := h.svc.ListMembers(r.Context(), wsID)
	if err != nil {
		common.Error(w, err)
		return
	}
	if list == nil {
		list = []*Membership{}
	}
	common.OK(w, list)
}
