package task

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	col *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{col: db.Collection("tasks")}
}

func (r *Repository) Create(ctx context.Context, t *Task) error {
	_, err := r.col.InsertOne(ctx, t)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id primitive.ObjectID) (*Task, error) {
	var t Task
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	if update["updated_at"] == nil {
		update["updated_at"] = time.Now()
	}
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (r *Repository) ListByProject(ctx context.Context, projectID primitive.ObjectID, skip, limit int64) ([]*Task, int64, error) {
	filter := bson.M{"project_id": projectID}
	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().SetSkip(skip).SetLimit(limit)
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)
	var out []*Task
	if err := cur.All(ctx, &out); err != nil {
		return nil, 0, err
	}
	if out == nil {
		out = []*Task{}
	}
	return out, total, nil
}
