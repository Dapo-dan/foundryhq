package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
)

// workspaceMemberModel is the GORM row shape for the workspace_members table
// (see migration 000001_init_schema).
type workspaceMemberModel struct {
	ID          uuid.UUID  `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkspaceID uuid.UUID  `gorm:"column:workspace_id"`
	UserID      uuid.UUID  `gorm:"column:user_id"`
	Role        string     `gorm:"column:role"`
	InvitedAt   time.Time  `gorm:"column:invited_at"`
	JoinedAt    *time.Time `gorm:"column:joined_at"`
}

func (workspaceMemberModel) TableName() string { return "workspace_members" }

func (m workspaceMemberModel) toDomain() *domain.WorkspaceMember {
	return &domain.WorkspaceMember{
		ID:          m.ID,
		WorkspaceID: m.WorkspaceID,
		UserID:      m.UserID,
		Role:        domain.WorkspaceRole(m.Role),
		InvitedAt:   m.InvitedAt,
		JoinedAt:    m.JoinedAt,
	}
}

func workspaceMemberModelFromDomain(member *domain.WorkspaceMember) *workspaceMemberModel {
	return &workspaceMemberModel{
		ID:          member.ID,
		WorkspaceID: member.WorkspaceID,
		UserID:      member.UserID,
		Role:        string(member.Role),
		JoinedAt:    member.JoinedAt,
	}
}

// WorkspaceMemberRepository implements domain.WorkspaceMemberRepository on
// top of GORM/Postgres.
type WorkspaceMemberRepository struct {
	db *gorm.DB
}

// NewWorkspaceMemberRepository constructs a WorkspaceMemberRepository.
func NewWorkspaceMemberRepository(db *gorm.DB) *WorkspaceMemberRepository {
	return &WorkspaceMemberRepository{db: db}
}

// Create inserts member, generating its ID and InvitedAt via the table's
// column defaults.
func (r *WorkspaceMemberRepository) Create(ctx context.Context, member *domain.WorkspaceMember) error {
	model := workspaceMemberModelFromDomain(member)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("creating workspace member: %w", err)
	}
	member.ID = model.ID
	member.InvitedAt = model.InvitedAt
	return nil
}

// GetByID returns the membership row with the given id, or
// domain.ErrWorkspaceMemberNotFound.
func (r *WorkspaceMemberRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.WorkspaceMember, error) {
	var model workspaceMemberModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrWorkspaceMemberNotFound
		}
		return nil, fmt.Errorf("getting workspace member %s: %w", id, err)
	}
	return model.toDomain(), nil
}

// GetByWorkspaceAndUser returns the membership row for userID in
// workspaceID, or domain.ErrWorkspaceMemberNotFound.
func (r *WorkspaceMemberRepository) GetByWorkspaceAndUser(ctx context.Context, workspaceID, userID uuid.UUID) (*domain.WorkspaceMember, error) {
	var model workspaceMemberModel
	if err := r.db.WithContext(ctx).
		First(&model, "workspace_id = ? AND user_id = ?", workspaceID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrWorkspaceMemberNotFound
		}
		return nil, fmt.Errorf("getting workspace member: %w", err)
	}
	return model.toDomain(), nil
}

// memberWithEmailRow is the destination shape for ListByWorkspaceIDWithUser's
// workspace_members/users join.
type memberWithEmailRow struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
	UserID      uuid.UUID
	Role        string
	InvitedAt   time.Time
	JoinedAt    *time.Time
	Email       string
}

// ListByWorkspaceIDWithUser returns every member of workspaceID, joined with
// users for each member's email — the members-list response needs a
// human-readable identifier, not just a UserID.
func (r *WorkspaceMemberRepository) ListByWorkspaceIDWithUser(ctx context.Context, workspaceID uuid.UUID) ([]*domain.WorkspaceMemberWithUser, error) {
	var rows []memberWithEmailRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT workspace_members.id, workspace_members.workspace_id, workspace_members.user_id,
		       workspace_members.role, workspace_members.invited_at, workspace_members.joined_at,
		       users.email AS email
		FROM workspace_members
		JOIN users ON users.id = workspace_members.user_id
		WHERE workspace_members.workspace_id = ?
		ORDER BY workspace_members.invited_at
	`, workspaceID).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("listing workspace members: %w", err)
	}

	members := make([]*domain.WorkspaceMemberWithUser, len(rows))
	for i, row := range rows {
		members[i] = &domain.WorkspaceMemberWithUser{
			WorkspaceMember: domain.WorkspaceMember{
				ID:          row.ID,
				WorkspaceID: row.WorkspaceID,
				UserID:      row.UserID,
				Role:        domain.WorkspaceRole(row.Role),
				InvitedAt:   row.InvitedAt,
				JoinedAt:    row.JoinedAt,
			},
			Email: row.Email,
		}
	}
	return members, nil
}

// UpdateRole replaces the role for the membership row with the given id.
// Returns domain.ErrWorkspaceMemberNotFound if no such row exists.
func (r *WorkspaceMemberRepository) UpdateRole(ctx context.Context, id uuid.UUID, role domain.WorkspaceRole) error {
	result := r.db.WithContext(ctx).
		Model(&workspaceMemberModel{}).
		Where("id = ?", id).
		Update("role", string(role))
	if result.Error != nil {
		return fmt.Errorf("updating workspace member %s role: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrWorkspaceMemberNotFound
	}
	return nil
}

// MarkJoined sets joined_at to now for the membership row with the given id.
// Returns domain.ErrWorkspaceMemberNotFound if no such row exists.
func (r *WorkspaceMemberRepository) MarkJoined(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Model(&workspaceMemberModel{}).
		Where("id = ?", id).
		Update("joined_at", time.Now())
	if result.Error != nil {
		return fmt.Errorf("marking workspace member %s joined: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrWorkspaceMemberNotFound
	}
	return nil
}
