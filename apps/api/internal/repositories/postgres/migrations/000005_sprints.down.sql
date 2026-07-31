DROP INDEX IF EXISTS idx_tasks_sprint_id;

ALTER TABLE tasks
    DROP COLUMN IF EXISTS due_date,
    DROP COLUMN IF EXISTS story_points,
    DROP COLUMN IF EXISTS priority,
    DROP COLUMN IF EXISTS sprint_id;

DROP INDEX IF EXISTS idx_sprints_workspace_id;
DROP TABLE IF EXISTS sprints;
