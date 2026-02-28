package config

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureIndexes creates required indexes: users.email unique, memberships (user_id+workspace_id) unique.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	users := db.Collection("users")
	_, err := users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]int{"email": 1},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	memberships := db.Collection("memberships")
	_, err = memberships.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    map[string]int{"user_id": 1, "workspace_id": 1},
		Options: options.Index().SetUnique(true),
	})
	return err
}
