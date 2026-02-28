package workspace

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MembershipStatus string

const (
	StatusPending  MembershipStatus = "PENDING"
	StatusApproved MembershipStatus = "APPROVED"
	StatusRejected MembershipStatus = "REJECTED"
)

type Membership struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	WorkspaceID primitive.ObjectID `bson:"workspace_id" json:"workspace_id"`

	Status MembershipStatus `bson:"status" json:"status"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
