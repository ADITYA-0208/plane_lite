package workspace

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workspace struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"_id"`

	Name string `bson:"name" json:"name"`

	AdminID primitive.ObjectID `bson:"admin_id" json:"admin_id"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
