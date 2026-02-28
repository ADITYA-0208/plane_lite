package project

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"planelite-backend/internal/common"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, workspaceID primitive.ObjectID, name string) (*Project, error) {
	if name == "" {
		return nil, common.ErrInvalidInput
	}
	p := &Project{
		Name:        name,
		WorkspaceID: workspaceID,
		CreatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) GetByID(ctx context.Context, id primitive.ObjectID) (*Project, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) ListByWorkspace(ctx context.Context, workspaceID primitive.ObjectID) ([]*Project, error) {
	return s.repo.ListByWorkspace(ctx, workspaceID)
}
