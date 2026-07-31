ALTER TABLE tasks DROP CONSTRAINT tasks_assignee_id_fkey;
ALTER TABLE tasks
    ADD CONSTRAINT tasks_assignee_id_fkey
    FOREIGN KEY (assignee_id) REFERENCES users (id);

ALTER TABLE tasks DROP CONSTRAINT tasks_workspace_id_fkey;
ALTER TABLE tasks
    ADD CONSTRAINT tasks_workspace_id_fkey
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id);

ALTER TABLE workspace_members DROP CONSTRAINT workspace_members_workspace_id_fkey;
ALTER TABLE workspace_members
    ADD CONSTRAINT workspace_members_workspace_id_fkey
    FOREIGN KEY (workspace_id) REFERENCES workspaces (id);
