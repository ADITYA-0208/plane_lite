package workspace

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"planelite-backend/internal/common"
)

type Service struct {
	repo     *Repository
	memRepo  *MembershipRepository
}

func NewService(repo *Repository, memRepo *MembershipRepository) *Service {
	return &Service{repo: repo, memRepo: memRepo}
}

// Create creates a workspace. Caller must be ADMIN; ADMIN can have only one workspace.
func (s *Service) Create(ctx context.Context, adminID primitive.ObjectID, name string) (*Workspace, error) {
	if name == "" {
		return nil, common.ErrInvalidInput
	}
	list, err := s.repo.ListByAdminID(ctx, adminID)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return nil, common.ErrConflict // ADMIN can belong to only one workspace
	}
	w := &Workspace{
		Name:      name,
		AdminID:   adminID,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

// GetByID returns workspace by ID.
func (s *Service) GetByID(ctx context.Context, id primitive.ObjectID) (*Workspace, error) {
	return s.repo.FindByID(ctx, id)
}

// AddMember adds a user to workspace as PENDING. Admin must approve later.
func (s *Service) AddMember(ctx context.Context, workspaceID, userID primitive.ObjectID) (*Membership, error) {
	existing, _ := s.memRepo.FindByUserAndWorkspace(ctx, userID, workspaceID)
	if existing != nil {
		return nil, common.ErrConflict // already a membership
	}
	m := &Membership{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
	}
	if err := s.memRepo.Create(ctx, m); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, common.ErrConflict
		}
		return nil, err
	}
	return m, nil
}

// ApproveMember sets membership to APPROVED. Caller must be workspace admin.
func (s *Service) ApproveMember(ctx context.Context, adminID, membershipID primitive.ObjectID) error {
	mem, err := s.memRepo.FindByID(ctx, membershipID)
	if err != nil {
		return common.ErrNotFound
	}
	ws, err := s.repo.FindByID(ctx, mem.WorkspaceID)
	if err != nil || ws.AdminID != adminID {
		return common.ErrForbidden
	}
	return s.memRepo.UpdateStatus(ctx, membershipID, StatusApproved)
}

// HasApprovedAccess returns true if user is workspace admin or has approved membership.
func (s *Service) HasApprovedAccess(ctx context.Context, userID, workspaceID primitive.ObjectID) (bool, error) {
	ws, err := s.repo.FindByID(ctx, workspaceID)
	if err != nil {
		return false, err
	}
	if ws.AdminID == userID {
		return true, nil
	}
	return s.memRepo.HasApproved(ctx, userID, workspaceID)
}

// ListByAdminID returns workspaces owned by admin.
func (s *Service) ListByAdminID(ctx context.Context, adminID primitive.ObjectID) ([]*Workspace, error) {
	return s.repo.ListByAdminID(ctx, adminID)
}

// ListMembers returns memberships for a workspace.
func (s *Service) ListMembers(ctx context.Context, workspaceID primitive.ObjectID) ([]*Membership, error) {
	return s.memRepo.ListByWorkspace(ctx, workspaceID)
}
