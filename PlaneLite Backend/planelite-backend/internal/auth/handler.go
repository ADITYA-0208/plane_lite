package auth

import (
	"encoding/json"
	"net/http"

	"planelite-backend/internal/common"
	"planelite-backend/internal/config"
)

// Handler handles signup and login; parses request and writes response only.
type Handler struct {
	svc *Service
	cfg *config.Config
}

// NewHandler returns an auth HTTP handler.
func NewHandler(svc *Service, cfg *config.Config) *Handler {
	return &Handler{svc: svc, cfg: cfg}
}
               
// SignupRequest is the JSON body for POST /auth/signup.
type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"` // optional; defaults to USER
}

// LoginRequest is the JSON body for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Signup handles POST /auth/signup.
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.Error(w, common.ErrBadRequest)
		return
	}
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	role := common.RoleUser
	if req.Role != "" {
		role = common.Role(req.Role)
	}
	u, token, err := h.svc.Signup(r.Context(), req.Email, req.Password, role)
	if err != nil {
		common.Error(w, err)
		return
	}
	common.Created(w, TokenResponse{
		AccessToken: token,
		UserID:      u.ID.Hex(),
		Email:       u.Email,
		Role:        string(u.Role),
	})
}

// Login handles POST /auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		common.Error(w, common.ErrBadRequest)
		return
	}
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.Error(w, common.ErrBadRequest)
		return
	}
	u, token, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		common.Error(w, err)
		return
	}
	common.OK(w, TokenResponse{
		AccessToken: token,
		UserID:      u.ID.Hex(),
		Email:       u.Email,
		Role:        string(u.Role),
	})
}
