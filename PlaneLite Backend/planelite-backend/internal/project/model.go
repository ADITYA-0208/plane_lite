package project

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	WorkspaceID primitive.ObjectID `bson:"workspace_id"`
	CreatedAt   time.Time          `bson:"created_at"`
}
