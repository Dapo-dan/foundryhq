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

// projectModel is the GORM row shape for the projects table (see migration
// 000003_projects). DeletedAt uses GORM's built-in soft-delete type, so
// every normal query (First, Find, Updates, Delete) automatically excludes
// already-deleted rows without repeating "deleted_at IS NULL" by hand.
type projectModel struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkspaceID uuid.UUID      `gorm:"column:workspace_id"`
	Name        string         `gorm:"column:name"`
	Description *string        `gorm:"column:description"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (projectModel) TableName() string { return "projects" }

func (m projectModel) toDomain() *domain.Project {
	project := &domain.Project{
		ID:          m.ID,
		WorkspaceID: m.WorkspaceID,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.DeletedAt.Valid {
		deletedAt := m.DeletedAt.Time
		project.DeletedAt = &deletedAt
	}
	return project
}

func projectModelFromDomain(p *domain.Project) *projectModel {
	return &projectModel{
		ID:          p.ID,
		WorkspaceID: p.WorkspaceID,
		Name:        p.Name,
		Description: p.Description,
	}
}

// ProjectRepository implements domain.ProjectRepository on top of
// GORM/Postgres.
type ProjectRepository struct {
	db *gorm.DB
}

// NewProjectRepository constructs a ProjectRepository.
func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// Create inserts project, generating its ID and timestamps via the table's
// column defaults.
func (r *ProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	model := projectModelFromDomain(project)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("creating project: %w", err)
	}
	project.ID = model.ID
	project.CreatedAt = model.CreatedAt
	project.UpdatedAt = model.UpdatedAt
	return nil
}

// GetByID returns the project with the given ID, or domain.ErrProjectNotFound
// — including for a soft-deleted project, which GORM's default scope treats
// as already gone.
func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	var model projectModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProjectNotFound
		}
		return nil, fmt.Errorf("getting project %s: %w", id, err)
	}
	return model.toDomain(), nil
}

// ListByWorkspaceID returns every non-deleted project in workspaceID.
func (r *ProjectRepository) ListByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]*domain.Project, error) {
	var models []projectModel
	if err := r.db.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Order("created_at").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("listing projects for workspace %s: %w", workspaceID, err)
	}

	projects := make([]*domain.Project, len(models))
	for i, model := range models {
		projects[i] = model.toDomain()
	}
	return projects, nil
}

// Update replaces project's name/description. Returns
// domain.ErrProjectNotFound if no such (non-deleted) project exists.
func (r *ProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	model := projectModelFromDomain(project)
	result := r.db.WithContext(ctx).
		Model(&projectModel{}).
		Where("id = ?", project.ID).
		Updates(map[string]any{
			"name":        model.Name,
			"description": model.Description,
		})
	if result.Error != nil {
		return fmt.Errorf("updating project %s: %w", project.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrProjectNotFound
	}
	return nil
}

// Delete soft-deletes project and, in the same transaction, every task
// belonging to it (see domain.ProjectRepository.Delete's doc comment).
func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&projectModel{}, "id = ?", id)
		if result.Error != nil {
			return fmt.Errorf("deleting project %s: %w", id, result.Error)
		}
		if result.RowsAffected == 0 {
			return domain.ErrProjectNotFound
		}

		if err := tx.Exec(
			`UPDATE tasks SET deleted_at = now(), updated_at = now() WHERE project_id = ? AND deleted_at IS NULL`,
			id,
		).Error; err != nil {
			return fmt.Errorf("soft-deleting tasks for project %s: %w", id, err)
		}
		return nil
	})
}
