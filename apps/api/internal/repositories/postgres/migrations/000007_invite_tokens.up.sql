-- Closes a security gap surfaced by a pre-deployment audit: WorkspaceUsecase.Invite
-- creates a placeholder user (no password) for an email that hasn't signed up yet,
-- and until now, POST /auth/register let ANYONE claim that placeholder by simply
-- registering with that email — no proof of ownership required. This table backs
-- a token-gated accept-invite flow (mirroring password_reset_tokens) so a
-- placeholder can only be activated by whoever received the emailed link.

CREATE TABLE invite_tokens (
    id                  uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_member_id uuid NOT NULL REFERENCES workspace_members (id) ON DELETE CASCADE,
    token_hash          text NOT NULL UNIQUE,
    expires_at          timestamptz NOT NULL,
    used_at             timestamptz,
    created_at          timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_invite_tokens_workspace_member_id ON invite_tokens (workspace_member_id);
