CREATE TABLE sprints (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id uuid NOT NULL REFERENCES workspaces (id) ON DELETE CASCADE,
    name         text NOT NULL,
    start_date   date NOT NULL,
    end_date     date NOT NULL CHECK (end_date >= start_date),
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_sprints_workspace_id ON sprints (workspace_id);

-- sprint_id is nullable (null = backlog, not any sprint) and ON DELETE SET
-- NULL mirrors assignee_id's existing pattern — deleting a sprint returns
-- its tasks to the backlog rather than deleting them.
ALTER TABLE tasks
    ADD COLUMN sprint_id     uuid REFERENCES sprints (id) ON DELETE SET NULL,
    ADD COLUMN priority      text NOT NULL DEFAULT 'medium' CHECK (priority IN ('urgent', 'high', 'medium', 'low')),
    ADD COLUMN story_points  integer CHECK (story_points IS NULL OR story_points >= 0),
    ADD COLUMN due_date      date;

CREATE INDEX idx_tasks_sprint_id ON tasks (sprint_id);
