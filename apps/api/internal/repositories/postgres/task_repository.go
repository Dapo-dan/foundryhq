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

// taskModel is the GORM row shape for the tasks table (see migrations
// 000001_init_schema and 000003_projects, which added project_id).
// DeletedAt uses GORM's built-in soft-delete type, same idiom as Phase 2's
// projectModel.
type taskModel struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;primaryKey;default:gen_random_uuid()"`
	WorkspaceID uuid.UUID      `gorm:"column:workspace_id"`
	ProjectID   uuid.UUID      `gorm:"column:project_id"`
	Title       string         `gorm:"column:title"`
	Status      string         `gorm:"column:status"`
	AssigneeID  *uuid.UUID     `gorm:"column:assignee_id"`
	SprintID    *uuid.UUID     `gorm:"column:sprint_id"`
	Priority    string         `gorm:"column:priority"`
	StoryPoints *int           `gorm:"column:story_points"`
	DueDate     *time.Time     `gorm:"column:due_date;type:date"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (taskModel) TableName() string { return "tasks" }

func (m taskModel) toDomain() *domain.Task {
	task := &domain.Task{
		ID:          m.ID,
		WorkspaceID: m.WorkspaceID,
		ProjectID:   m.ProjectID,
		Title:       m.Title,
		Status:      domain.TaskStatus(m.Status),
		AssigneeID:  m.AssigneeID,
		SprintID:    m.SprintID,
		Priority:    domain.TaskPriority(m.Priority),
		StoryPoints: m.StoryPoints,
		DueDate:     m.DueDate,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.DeletedAt.Valid {
		deletedAt := m.DeletedAt.Time
		task.DeletedAt = &deletedAt
	}
	return task
}

func taskModelFromDomain(t *domain.Task) *taskModel {
	return &taskModel{
		ID:          t.ID,
		WorkspaceID: t.WorkspaceID,
		ProjectID:   t.ProjectID,
		Title:       t.Title,
		Status:      string(t.Status),
		AssigneeID:  t.AssigneeID,
		SprintID:    t.SprintID,
		Priority:    string(t.Priority),
		StoryPoints: t.StoryPoints,
		DueDate:     t.DueDate,
	}
}

// TaskRepository implements domain.TaskRepository on top of GORM/Postgres.
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository constructs a TaskRepository.
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create inserts task, generating its ID and timestamps via the table's
// column defaults. Status/Priority default to domain.StatusTodo/PriorityMedium
// if unset, matching the columns' own DEFAULTs.
func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	model := taskModelFromDomain(task)
	if model.Status == "" {
		model.Status = string(domain.StatusTodo)
	}
	if model.Priority == "" {
		model.Priority = string(domain.PriorityMedium)
	}
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("creating task: %w", err)
	}
	task.ID = model.ID
	task.Status = domain.TaskStatus(model.Status)
	task.Priority = domain.TaskPriority(model.Priority)
	task.CreatedAt = model.CreatedAt
	task.UpdatedAt = model.UpdatedAt
	return nil
}

// GetByID returns the task with the given ID, or domain.ErrTaskNotFound —
// including for a soft-deleted task, which GORM's default scope treats as
// already gone.
func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Task, error) {
	var model taskModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("getting task %s: %w", id, err)
	}
	return model.toDomain(), nil
}

// ListByWorkspaceID returns every non-deleted task in workspaceID, narrowed
// by filter's non-nil fields.
func (r *TaskRepository) ListByWorkspaceID(ctx context.Context, workspaceID uuid.UUID, filter domain.TaskFilter) ([]*domain.Task, error) {
	query := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID)
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.AssigneeID != nil {
		query = query.Where("assignee_id = ?", *filter.AssigneeID)
	}
	if filter.SprintID != nil {
		query = query.Where("sprint_id = ?", *filter.SprintID)
	}

	var models []taskModel
	if err := query.Order("created_at").Find(&models).Error; err != nil {
		return nil, fmt.Errorf("listing tasks for workspace %s: %w", workspaceID, err)
	}

	tasks := make([]*domain.Task, len(models))
	for i, model := range models {
		tasks[i] = model.toDomain()
	}
	return tasks, nil
}

// Update replaces task's project/title/status/assignee/sprint/priority/
// story points/due date. Returns domain.ErrTaskNotFound if no such
// (non-deleted) task exists.
func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	model := taskModelFromDomain(task)
	result := r.db.WithContext(ctx).
		Model(&taskModel{}).
		Where("id = ?", task.ID).
		Updates(map[string]any{
			"project_id":   model.ProjectID,
			"title":        model.Title,
			"status":       model.Status,
			"assignee_id":  model.AssigneeID,
			"sprint_id":    model.SprintID,
			"priority":     model.Priority,
			"story_points": model.StoryPoints,
			"due_date":     model.DueDate,
		})
	if result.Error != nil {
		return fmt.Errorf("updating task %s: %w", task.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}

// Delete soft-deletes task. Returns domain.ErrTaskNotFound if no such
// (non-deleted) task exists.
func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&taskModel{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("deleting task %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}

// SumStoryPointsForSprint returns the sum of story_points for done tasks in
// sprintID whose updated_at falls within [startDate, endDate] — the end
// bound is exclusive of the *next* day so the entirety of endDate's
// calendar day counts, matching docs/database.md's closed date-range
// wording rather than a literal timestamp boundary at midnight.
func (r *TaskRepository) SumStoryPointsForSprint(ctx context.Context, sprintID uuid.UUID, startDate, endDate time.Time) (int, error) {
	var sum int
	err := r.db.WithContext(ctx).
		Model(&taskModel{}).
		Select("COALESCE(SUM(story_points), 0)").
		Where("sprint_id = ? AND status = ? AND updated_at >= ? AND updated_at < ?",
			sprintID, string(domain.StatusDone), startDate, endDate.AddDate(0, 0, 1)).
		Scan(&sum).Error
	if err != nil {
		return 0, fmt.Errorf("summing story points for sprint %s: %w", sprintID, err)
	}
	return sum, nil
}
