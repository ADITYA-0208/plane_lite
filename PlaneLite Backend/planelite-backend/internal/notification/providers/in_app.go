package providers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// InAppProvider stores in-app notifications (e.g. in a collection for polling or SSE).
type InAppProvider struct{}

func NewInAppProvider() *InAppProvider {
	return &InAppProvider{}
}

func (p *InAppProvider) Send(ctx context.Context, userID primitive.ObjectID, title, body string) error {
	// Stub: persist to notifications collection; real impl would use DB.
	_ = userID
	_, _ = title, body
	return nil
}
