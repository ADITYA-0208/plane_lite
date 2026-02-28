package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"planelite-backend/internal/common"
	"planelite-backend/internal/config"
	"planelite-backend/internal/user"
)

var _ error = (*validationError)(nil)

type validationError struct{ msg string }

func (e validationError) Error() string { return e.msg }


// Service handles signup, login, and JWT issuance/validation. No DB of its own.
type Service struct {
	user  *user.Service
	cfg   *config.Config
}

// NewService creates an auth service that uses user service and app config.
func NewService(userSvc *user.Service, cfg *config.Config) *Service {
	return &Service{user: userSvc, cfg: cfg}
}

// Signup creates a user and returns a JWT. Business rule: only first user or admin flow can create ADMIN.
func (s *Service) Signup(ctx context.Context, email, password string, role common.Role) (*user.User, string, error) {
	if email == "" || password == "" {
		return nil, "", fmt.Errorf("%w: email and password required", common.ErrInvalidInput)
	}
	if role != common.RoleAdmin && role != common.RoleUser && role != common.RoleProjectManager {
		role = common.RoleUser
	}
	if err := s.user.Create(ctx, email, password, role); err != nil {
		return nil, "", err
	}
	u, err := s.user.FindByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}
	token, err := s.issueToken(u.ID.Hex(), string(u.Role))
	if err != nil {
		return u, "", err
	}
	return u, token, nil
}

// Login authenticates and returns user + JWT.
func (s *Service) Login(ctx context.Context, email, password string) (*user.User, string, error) {
	if email == "" || password == "" {
		return nil, "", fmt.Errorf("%w: email and password required", common.ErrInvalidInput)
	}
	u, err := s.user.Authenticate(ctx, email, password)
	if err != nil {
		return nil, "", common.ErrUnauthorized
	}
	token, err := s.issueToken(u.ID.Hex(), string(u.Role))
	if err != nil {
		return u, "", err
	}
	return u, token, nil
}

// ValidateToken parses the JWT and returns claims. Returns nil, error if invalid.
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*Claims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *Service) issueToken(userID, role string) (string, error) {
	exp := time.Duration(s.cfg.JWTExpiryHours) * time.Hour
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
		Role:   common.Role(role),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(s.cfg.JWTSecret))
}

// GetUserByID is a convenience that delegates to user service; used by handlers that only have ID.
func (s *Service) GetUserByID(ctx context.Context, id primitive.ObjectID) (*user.User, error) {
	return s.user.GetByID(ctx, id)
}
