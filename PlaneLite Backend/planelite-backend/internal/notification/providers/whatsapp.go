package providers

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WhatsAppProvider sends notifications via WhatsApp. Stub; real impl would use WhatsApp Business API.
type WhatsAppProvider struct{}

func NewWhatsAppProvider() *WhatsAppProvider {
	return &WhatsAppProvider{}
}

func (p *WhatsAppProvider) Send(ctx context.Context, userID primitive.ObjectID, title, body string) error {
	// Stub: would resolve user phone and call WhatsApp API.
	_ = userID
	_, _ = title, body
	return nil
}
