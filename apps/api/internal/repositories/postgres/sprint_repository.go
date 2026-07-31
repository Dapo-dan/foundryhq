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

// sprintModel is the GORM row shape for the sprints table (see migration
// 000005_sprints). StartDate/EndDate use gorm's "date" type — Postgres
// returns them as midnight-UTC time.Time values, no time-of-day component.
type sprintModel struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkspaceID uuid.UUID `gorm:"column:workspace_id"`
	Name        string    `gorm:"column:name"`
	StartDate   time.Time `gorm:"column:start_date;type:date"`
	EndDate     time.Time `gorm:"column:end_date;type:date"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (sprintModel) TableName() string { return "sprints" }

func (m sprintModel) toDomain() *domain.Sprint {
	return &domain.Sprint{
		ID:          m.ID,
		WorkspaceID: m.WorkspaceID,
		Name:        m.Name,
		StartDate:   m.StartDate,
		EndDate:     m.EndDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func sprintModelFromDomain(s *domain.Sprint) *sprintModel {
	return &sprintModel{
		ID:          s.ID,
		WorkspaceID: s.WorkspaceID,
		Name:        s.Name,
		StartDate:   s.StartDate,
		EndDate:     s.EndDate,
	}
}

// SprintRepository implements domain.SprintRepository on top of
// GORM/Postgres.
type SprintRepository struct {
	db *gorm.DB
}

// NewSprintRepository constructs a SprintRepository.
func NewSprintRepository(db *gorm.DB) *SprintRepository {
	return &SprintRepository{db: db}
}

// Create inserts sprint, generating its ID and timestamps via the table's
// column defaults.
func (r *SprintRepository) Create(ctx context.Context, sprint *domain.Sprint) error {
	model := sprintModelFromDomain(sprint)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("creating sprint: %w", err)
	}
	sprint.ID = model.ID
	sprint.CreatedAt = model.CreatedAt
	sprint.UpdatedAt = model.UpdatedAt
	return nil
}

// GetByID returns the sprint with the given ID, or domain.ErrSprintNotFound.
func (r *SprintRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Sprint, error) {
	var model sprintModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSprintNotFound
		}
		return nil, fmt.Errorf("getting sprint %s: %w", id, err)
	}
	return model.toDomain(), nil
}

// ListByWorkspaceID returns every sprint in workspaceID, most recently
// started first.
func (r *SprintRepository) ListByWorkspaceID(ctx context.Context, workspaceID uuid.UUID) ([]*domain.Sprint, error) {
	var models []sprintModel
	if err := r.db.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Order("start_date desc").
		Find(&models).Error; err != nil {
		return nil, fmt.Errorf("listing sprints for workspace %s: %w", workspaceID, err)
	}

	sprints := make([]*domain.Sprint, len(models))
	for i, model := range models {
		sprints[i] = model.toDomain()
	}
	return sprints, nil
}
