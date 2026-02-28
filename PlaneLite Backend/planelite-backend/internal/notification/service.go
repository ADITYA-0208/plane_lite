package notification

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"planelite-backend/internal/notification/providers"
)

// Service dispatches notifications via providers (in-app, whatsapp). Stub.
type Service struct {
	InApp    *providers.InAppProvider
	WhatsApp *providers.WhatsAppProvider
}

func NewService(inApp *providers.InAppProvider, whatsApp *providers.WhatsAppProvider) *Service {
	return &Service{InApp: inApp, WhatsApp: whatsApp}
}

// Notify sends to user via configured providers. Stub implementation.
func (s *Service) Notify(ctx context.Context, userID primitive.ObjectID, title, body string) error {
	if s.InApp != nil {
		_ = s.InApp.Send(ctx, userID, title, body)
	}
	return nil
}
