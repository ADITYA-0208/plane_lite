package user

import (
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

// GetMe returns the current user from JWT. Use after Auth middleware.
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.Error(w, common.ErrBadRequest)
		return
	}
	u := common.GetContextUser(r.Context())
	if u == nil || u.UserID == "" {
		common.Error(w, common.ErrUnauthorized)
		return
	}
	userID, err := primitive.ObjectIDFromHex(u.UserID)
	if err != nil {
		common.Error(w, common.ErrUnauthorized)
		return
	}
	usr, err := h.svc.GetByID(r.Context(), userID)
	if err != nil {
		common.Error(w, common.ErrNotFound)
		return
	}
	// Do not expose password
	out := map[string]any{
		"id":         usr.ID.Hex(),
		"email":      usr.Email,
		"role":       string(usr.Role),
		"created_at": usr.CreatedAt,
	}
	common.OK(w, out)
}

// GetByID returns a user by ID (e.g. for workspace member list). Admin/PM only in real impl.
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
	u, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		common.Error(w, common.ErrNotFound)
		return
	}
	out := map[string]any{
		"id":         u.ID.Hex(),
		"email":      u.Email,
		"role":       string(u.Role),
		"created_at": u.CreatedAt,
	}
	common.OK(w, out)
}
