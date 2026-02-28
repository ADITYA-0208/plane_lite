package workspace

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
	return &Repository{col: db.Collection("workspaces")}
}

func (r *Repository) Create(ctx context.Context, w *Workspace) error {
	// Insert without _id so MongoDB generates one (omitempty can still serialize zero ObjectID).
	doc := bson.M{
		"name":       w.Name,
		"admin_id":   w.AdminID,
		"created_at": w.CreatedAt,
	}
	result, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		w.ID = oid
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID) (*Workspace, error) {
	var w Workspace
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&w)
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *Repository) ListByAdminID(ctx context.Context, adminID primitive.ObjectID) ([]*Workspace, error) {
	cur, err := r.col.Find(ctx, bson.M{"admin_id": adminID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []*Workspace
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
