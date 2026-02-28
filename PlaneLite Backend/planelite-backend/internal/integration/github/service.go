package github

import (
	"context"
	"net/http"
)

// Service integrates with GitHub (e.g. link repo, sync issues). Stub.
type Service struct {
	Token string // from env, never hardcoded
}

func NewService(token string) *Service {
	return &Service{Token: token}
}

// GetClient returns an HTTP client with GitHub auth. Stub.
func (s *Service) GetClient(ctx context.Context) *http.Client {
	_ = ctx
	return http.DefaultClient
}
