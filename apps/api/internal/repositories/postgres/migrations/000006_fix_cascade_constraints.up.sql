-- Fixes a cascade/doc mismatch left over from 000001_init_schema:
-- docs/database.md documents workspace_members.workspace_id and
-- tasks.workspace_id as ON DELETE CASCADE, and tasks.assignee_id as
-- ON DELETE SET NULL, but the original migration created all three FKs with
-- no ON DELETE clause at all (Postgres default: RESTRICT). This went
-- unnoticed while those tables were empty; it surfaced once Workspaces,
-- Projects, Tasks, and Sprints started writing real rows through them — a
-- hard-deleted workspace today would fail instead of cascading to its
-- members/tasks as documented, and a deleted user would fail instead of
-- clearing tasks.assignee_id.

ALTER TABLE workspace_members DROP CONSTRAINT workspace_members_workspace_id_fkey;
ALTER TABLE workspace_members
    ADD CONSTRAINT workspace_members_workspace_id_fkey
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE;

ALTER TABLE tasks DROP CONSTRAINT tasks_workspace_id_fkey;
ALTER TABLE tasks
    ADD CONSTRAINT tasks_workspace_id_fkey
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id) ON DELETE CASCADE;

ALTER TABLE tasks DROP CONSTRAINT tasks_assignee_id_fkey;
ALTER TABLE tasks
    ADD CONSTRAINT tasks_assignee_id_fkey
    FOREIGN KEY (assignee_id) REFERENCES users (id) ON DELETE SET NULL;
