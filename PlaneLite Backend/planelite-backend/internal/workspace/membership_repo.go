package workspace

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MembershipRepository handles membership collection.
type MembershipRepository struct {
	col *mongo.Collection
}

func NewMembershipRepository(db *mongo.Database) *MembershipRepository {
	return &MembershipRepository{col: db.Collection("memberships")}
}

func (r *MembershipRepository) Create(ctx context.Context, m *Membership) error {
	doc := bson.M{
		"user_id":      m.UserID,
		"workspace_id": m.WorkspaceID,
		"status":       m.Status,
		"created_at":   m.CreatedAt,
	}
	result, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		m.ID = oid
	}
	return nil
}

func (r *MembershipRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*Membership, error) {
	var m Membership
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MembershipRepository) FindByUserAndWorkspace(ctx context.Context, userID, workspaceID primitive.ObjectID) (*Membership, error) {
	var m Membership
	err := r.col.FindOne(ctx, bson.M{"user_id": userID, "workspace_id": workspaceID}).Decode(&m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MembershipRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status MembershipStatus) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	return err
}

func (r *MembershipRepository) HasApproved(ctx context.Context, userID, workspaceID primitive.ObjectID) (bool, error) {
	n, err := r.col.CountDocuments(ctx, bson.M{
		"user_id": userID, "workspace_id": workspaceID, "status": StatusApproved,
	})
	return n > 0, err
}

func (r *MembershipRepository) ListByWorkspace(ctx context.Context, workspaceID primitive.ObjectID) ([]*Membership, error) {
	cur, err := r.col.Find(ctx, bson.M{"workspace_id": workspaceID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []*Membership
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
