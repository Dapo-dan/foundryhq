-- V1 MVP schema: Auth + Workspace/Team + Tasks (see docs/mvp.md).
-- Full reference schema (all modules) is documented in docs/database.md;
-- tables/columns not needed until v1.1+ (sprints, priorities, due dates, etc.)
-- are added in a later migration when those features ship.

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE workspaces (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name       text NOT NULL,
    slug       text NOT NULL UNIQUE,
    logo_url   text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE users (
    id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email          text NOT NULL UNIQUE,
    password_hash  text,
    oauth_provider text,
    oauth_id       text,
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE workspace_members (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id uuid NOT NULL REFERENCES workspaces (id),
    user_id      uuid NOT NULL REFERENCES users (id),
    role         text NOT NULL CHECK (role IN ('owner', 'admin', 'member', 'viewer')),
    invited_at   timestamptz NOT NULL DEFAULT now(),
    joined_at    timestamptz,
    UNIQUE (workspace_id, user_id)
);

CREATE INDEX idx_workspace_members_workspace_id ON workspace_members (workspace_id);

CREATE TABLE tasks (
    id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id uuid NOT NULL REFERENCES workspaces (id),
    title        text NOT NULL,
    status       text NOT NULL DEFAULT 'todo' CHECK (status IN ('todo', 'in_progress', 'done')),
    assignee_id  uuid REFERENCES users (id),
    created_at   timestamptz NOT NULL DEFAULT now(),
    updated_at   timestamptz NOT NULL DEFAULT now(),
    deleted_at   timestamptz
);

CREATE INDEX idx_tasks_workspace_id ON tasks (workspace_id);
CREATE INDEX idx_tasks_assignee_id ON tasks (assignee_id);
