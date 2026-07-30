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

// workspaceModel is the GORM row shape for the workspaces table (see
// migration 000001_init_schema). LogoURL is a pointer because the column is
// nullable — domain.Workspace represents "no logo" as "".
type workspaceModel struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"column:name"`
	Slug      string    `gorm:"column:slug"`
	LogoURL   *string   `gorm:"column:logo_url"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (workspaceModel) TableName() string { return "workspaces" }

func (m workspaceModel) toDomain() *domain.Workspace {
	workspace := &domain.Workspace{
		ID:        m.ID,
		Name:      m.Name,
		Slug:      m.Slug,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.LogoURL != nil {
		workspace.LogoURL = *m.LogoURL
	}
	return workspace
}

func workspaceModelFromDomain(w *domain.Workspace) *workspaceModel {
	model := &workspaceModel{ID: w.ID, Name: w.Name, Slug: w.Slug}
	if w.LogoURL != "" {
		model.LogoURL = &w.LogoURL
	}
	return model
}

// WorkspaceRepository implements domain.WorkspaceRepository on top of
// GORM/Postgres.
type WorkspaceRepository struct {
	db *gorm.DB
}

// NewWorkspaceRepository constructs a WorkspaceRepository.
func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

// Create inserts workspace and owner's membership row in the same DB
// transaction — see domain.WorkspaceRepository.Create's doc comment for why
// these two inserts must succeed or fail together. owner.WorkspaceID is set
// from the newly-generated workspace ID before its row is inserted.
func (r *WorkspaceRepository) Create(ctx context.Context, workspace *domain.Workspace, owner *domain.WorkspaceMember) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		workspaceRow := workspaceModelFromDomain(workspace)
		if err := tx.Create(workspaceRow).Error; err != nil {
			return fmt.Errorf("creating workspace: %w", err)
		}
		workspace.ID = workspaceRow.ID
		workspace.CreatedAt = workspaceRow.CreatedAt
		workspace.UpdatedAt = workspaceRow.UpdatedAt

		owner.WorkspaceID = workspace.ID
		memberRow := workspaceMemberModelFromDomain(owner)
		if err := tx.Create(memberRow).Error; err != nil {
			return fmt.Errorf("creating owner membership: %w", err)
		}
		owner.ID = memberRow.ID
		owner.InvitedAt = memberRow.InvitedAt
		return nil
	})
}

// GetByID returns the workspace with the given ID, or domain.ErrWorkspaceNotFound.
func (r *WorkspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Workspace, error) {
	var model workspaceModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrWorkspaceNotFound
		}
		return nil, fmt.Errorf("getting workspace %s: %w", id, err)
	}
	return model.toDomain(), nil
}

// Update replaces workspace's name/slug/logo_url. Returns
// domain.ErrWorkspaceNotFound if no such workspace exists.
func (r *WorkspaceRepository) Update(ctx context.Context, workspace *domain.Workspace) error {
	model := workspaceModelFromDomain(workspace)
	result := r.db.WithContext(ctx).
		Model(&workspaceModel{}).
		Where("id = ?", workspace.ID).
		Updates(map[string]any{
			"name":     model.Name,
			"slug":     model.Slug,
			"logo_url": model.LogoURL,
		})
	if result.Error != nil {
		return fmt.Errorf("updating workspace %s: %w", workspace.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrWorkspaceNotFound
	}
	return nil
}

// SlugExists reports whether slug is already taken by some workspace.
func (r *WorkspaceRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&workspaceModel{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return false, fmt.Errorf("checking slug existence: %w", err)
	}
	return count > 0, nil
}

// ListForUser returns every workspace userID is a member of.
func (r *WorkspaceRepository) ListForUser(ctx context.Context, userID uuid.UUID) ([]*domain.Workspace, error) {
	var models []workspaceModel
	if err := r.db.WithContext(ctx).
		Joins("JOIN workspace_members ON workspace_members.workspace_id = workspaces.id").
		Where("workspace_members.user_id = ?", userID).
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("listing workspaces for user %s: %w", userID, err)
	}

	workspaces := make([]*domain.Workspace, len(models))
	for i, model := range models {
		workspaces[i] = model.toDomain()
	}
	return workspaces, nil
}
