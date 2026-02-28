package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"planelite-backend/internal/common"
)

type contextKey string

const claimsKey contextKey = "claims"

// WithClaims attaches claims to context.
func WithClaims(ctx context.Context, c *Claims) context.Context {
	return context.WithValue(ctx, claimsKey, c)
}

// ContextClaims returns claims from context, or nil.
func ContextClaims(ctx context.Context) *Claims {
	c, _ := ctx.Value(claimsKey).(*Claims)
	return c
}

// Claims holds JWT payload: user_id and role for auth and role middleware.
type Claims struct {
	jwt.RegisteredClaims
	UserID string      `json:"user_id"`
	Role   common.Role `json:"role"`
}

// TokenResponse is returned on login/signup.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}

// ObjectIDFromClaims returns UserID as ObjectID if valid.
func ObjectIDFromClaims(c *Claims) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.UserID)
	return id, err == nil
}
