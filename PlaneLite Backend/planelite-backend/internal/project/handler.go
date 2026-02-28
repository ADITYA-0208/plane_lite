package project

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

type CreateRequest struct {
	Name string `json:"name"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.Error(w, common.ErrBadRequest)
		return
	}
	wsID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	p, err := h.svc.Create(r.Context(), wsID, req.Name)
	if err != nil {
		common.Error(w, err)
		return
	}
	common.Created(w, p)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.Error(w, common.ErrBadRequest)
		return
	}
	pid := r.PathValue("pid")
	if pid == "" {
		common.Error(w, common.ErrBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(pid)
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		common.Error(w, common.ErrNotFound)
		return
	}
	common.OK(w, p)
}

func (h *Handler) ListByWorkspace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.Error(w, common.ErrBadRequest)
		return
	}
	wsID, err := primitive.ObjectIDFromHex(r.PathValue("id"))
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	list, err := h.svc.ListByWorkspace(r.Context(), wsID)
	if err != nil {
		common.Error(w, err)
		return
	}
	if list == nil {
		list = []*Project{}
	}
	common.OK(w, list)
}
