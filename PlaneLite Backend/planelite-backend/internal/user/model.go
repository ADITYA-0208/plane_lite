package user

import (
	"time"

	"planelite-backend/internal/common"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Role      common.Role        `bson:"role"`
	CreatedAt time.Time          `bson:"created_at"`
}

