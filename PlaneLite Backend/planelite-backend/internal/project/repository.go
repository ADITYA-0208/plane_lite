package project

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	col *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{col: db.Collection("projects")}
}

func (r *Repository) Create(ctx context.Context, p *Project) error {
	_, err := r.col.InsertOne(ctx, p)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID) (*Project, error) {
	var p Project
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) ListByWorkspace(ctx context.Context, workspaceID primitive.ObjectID) ([]*Project, error) {
	cur, err := r.col.Find(ctx, bson.M{"workspace_id": workspaceID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []*Project
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
