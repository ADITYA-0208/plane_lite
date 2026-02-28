package activity

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Service records activity for audit/feed. Stub: can be backed by a collection later.
type Service struct {
	col *mongo.Collection
}

func NewService(db *mongo.Database) *Service {
	return &Service{col: db.Collection("activities")}
}

func (s *Service) Record(ctx context.Context, workspaceID, projectID, taskID, userID primitive.ObjectID, kind ActivityKind, payload map[string]any) error {
	a := &Activity{
		WorkspaceID: workspaceID,
		ProjectID:   projectID,
		TaskID:      taskID,
		UserID:      userID,
		Kind:        kind,
		Payload:     payload,
		CreatedAt:   time.Now(),
	}
	_, err := s.col.InsertOne(ctx, a)
	return err
}
