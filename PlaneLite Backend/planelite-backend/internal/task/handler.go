package task

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
}

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
	if !CanCreateTask(u.Role) {
		common.Error(w, common.ErrForbidden)
		return
	}
	createdBy, err := primitive.ObjectIDFromHex(u.UserID)
	if err != nil {
		common.Error(w, common.ErrUnauthorized)
		return
	}
	projectID, err := primitive.ObjectIDFromHex(r.PathValue("pid"))
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	t, err := h.svc.Create(r.Context(), projectID, createdBy, req.Title, req.Description)
	if err != nil {
		common.Error(w, err)
		return
	}
	common.Created(w, t)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.Error(w, common.ErrBadRequest)
		return
	}
	tid := r.PathValue("tid")
	if tid == "" {
		common.Error(w, common.ErrBadRequest)
		return
	}
	id, err := primitive.ObjectIDFromHex(tid)
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	t, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		common.Error(w, common.ErrNotFound)
		return
	}
	common.OK(w, t)
}

func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		common.Error(w, common.ErrBadRequest)
		return
	}
	u := common.GetContextUser(r.Context())
	if u == nil || !CanUpdateTaskStatusOrPriority(u.Role) {
		common.Error(w, common.ErrForbidden)
		return
	}
	tid, err := primitive.ObjectIDFromHex(r.PathValue("tid"))
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	var req struct {
		Status TaskStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Status == "" {
		common.Error(w, common.ErrBadRequest)
		return
	}
	if err := h.svc.UpdateStatus(r.Context(), tid, req.Status); err != nil {
		common.Error(w, err)
		return
	}
	common.NoContent(w)
}

func (h *Handler) UpdatePriority(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		common.Error(w, common.ErrBadRequest)
		return
	}
	u := common.GetContextUser(r.Context())
	if u == nil || !CanUpdateTaskStatusOrPriority(u.Role) {
		common.Error(w, common.ErrForbidden)
		return
	}
	tid, err := primitive.ObjectIDFromHex(r.PathValue("tid"))
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	var req struct {
		Priority TaskPriority `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Priority == "" {
		common.Error(w, common.ErrBadRequest)
		return
	}
	if err := h.svc.UpdatePriority(r.Context(), tid, req.Priority); err != nil {
		common.Error(w, err)
		return
	}
	common.NoContent(w)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPut {
		common.Error(w, common.ErrBadRequest)
		return
	}
	u := common.GetContextUser(r.Context())
	if u == nil || !CanUpdateTaskFull(u.Role) {
		common.Error(w, common.ErrForbidden)
		return
	}
	tid, err := primitive.ObjectIDFromHex(r.PathValue("tid"))
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	if err := h.svc.Update(r.Context(), tid, req.Title, req.Description, req.Status, req.Priority); err != nil {
		common.Error(w, err)
		return
	}
	common.NoContent(w)
}

func (h *Handler) ListByProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		common.Error(w, common.ErrBadRequest)
		return
	}
	pid, err := primitive.ObjectIDFromHex(r.PathValue("pid"))
	if err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = common.DefaultPageSize
	}
	if pageSize > common.MaxPageSize {
		pageSize = common.MaxPageSize
	}
	skip := int64((page - 1) * pageSize)
	limit := int64(pageSize)
	list, total, err := h.svc.ListByProject(r.Context(), pid, skip, limit)
	if err != nil {
		common.Error(w, err)
		return
	}
	if list == nil {
		list = []*Task{}
	}
	common.OK(w, map[string]any{
		"items":        list,
		"page":         page,
		"page_size":    pageSize,
		"total_count":  total,
	})
}
