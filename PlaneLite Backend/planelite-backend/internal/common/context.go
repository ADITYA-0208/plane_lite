package common

import "context"

type contextKey string

const contextUserKey contextKey = "user"

// ContextUser holds authenticated user info set by auth middleware.
type ContextUser struct {
	UserID string
	Role   Role
}

// WithContextUser attaches user to context. Used by auth middleware only.
func WithContextUser(ctx context.Context, u ContextUser) context.Context {
	return context.WithValue(ctx, contextUserKey, &u)
}

// ContextUser returns the authenticated user from context, or nil.
func GetContextUser(ctx context.Context) *ContextUser {
	u, _ := ctx.Value(contextUserKey).(*ContextUser)
	return u
}
