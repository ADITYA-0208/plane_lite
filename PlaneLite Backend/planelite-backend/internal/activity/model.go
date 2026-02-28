package activity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ActivityKind string

const (
	KindTaskCreated   ActivityKind = "task_created"
	KindTaskUpdated   ActivityKind = "task_updated"
	KindTaskCompleted ActivityKind = "task_completed"
	KindMemberAdded   ActivityKind = "member_added"
)

type Activity struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	WorkspaceID primitive.ObjectID `bson:"workspace_id"`
	ProjectID  primitive.ObjectID `bson:"project_id,omitempty"`
	TaskID     primitive.ObjectID `bson:"task_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	Kind       ActivityKind       `bson:"kind"`
	Payload    map[string]any     `bson:"payload,omitempty"`
	CreatedAt  time.Time          `bson:"created_at"`
}
