package task

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"planelite-backend/internal/common"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, projectID, createdBy primitive.ObjectID, title, description string) (*Task, error) {
	if title == "" {
		return nil, common.ErrInvalidInput
	}
	t := &Task{
		Title:       title,
		Description: description,
		ProjectID:   projectID,
		Status:      StatusTodo,
		Priority:    PriorityMedium,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetByID(ctx context.Context, id primitive.ObjectID) (*Task, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) UpdateStatus(ctx context.Context, id primitive.ObjectID, status TaskStatus) error {
	if status != StatusTodo && status != StatusInProgress && status != StatusDone {
		return common.ErrInvalidInput
	}
	return s.repo.Update(ctx, id, bson.M{"status": status, "updated_at": time.Now()})
}

func (s *Service) UpdatePriority(ctx context.Context, id primitive.ObjectID, priority TaskPriority) error {
	if priority != PriorityLow && priority != PriorityMedium && priority != PriorityHigh {
		return common.ErrInvalidInput
	}
	return s.repo.Update(ctx, id, bson.M{"priority": priority, "updated_at": time.Now()})
}

func (s *Service) Update(ctx context.Context, id primitive.ObjectID, title, description string, status TaskStatus, priority TaskPriority) error {
	up := bson.M{"updated_at": time.Now()}
	if title != "" {
		up["title"] = title
	}
	if description != "" {
		up["description"] = description
	}
	if status != "" {
		up["status"] = status
	}
	if priority != "" {
		up["priority"] = priority
	}
	return s.repo.Update(ctx, id, up)
}

func (s *Service) ListByProject(ctx context.Context, projectID primitive.ObjectID, skip, limit int64) ([]*Task, int64, error) {
	return s.repo.ListByProject(ctx, projectID, skip, limit)
}
