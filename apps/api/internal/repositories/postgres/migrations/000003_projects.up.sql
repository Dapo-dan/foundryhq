-- Closes a schema/doc gap: docs/database.md documents `projects` as a
-- required V1 table (every task belongs to exactly one project) and
-- `tasks.project_id` as a not-null FK, but 000001_init_schema never created
-- either. project_id can go straight to NOT NULL with no backfill — no
-- usecase/repository writes to `tasks` yet, so the table is empty in every
-- environment.

CREATE TABLE projects (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id uuid NOT NULL REFERENCES workspaces (id) ON DELETE CASCADE,
    name         text NOT NULL,
    description  text,
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    deleted_at   timestamptz
);

CREATE INDEX idx_projects_workspace_id ON projects (workspace_id);

ALTER TABLE tasks
    ADD COLUMN project_id uuid NOT NULL REFERENCES projects (id) ON DELETE CASCADE;

CREATE INDEX idx_tasks_project_id ON tasks (project_id);
