-- Enables explicit server-side refresh-token invalidation on logout
-- (adr/0004-jwt-access-refresh-tokens.md). Stores a SHA-256 hash of the
-- refresh token, never the raw token, so a DB read alone can't be used to
-- forge a session.

CREATE TABLE refresh_tokens (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash text NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens (user_id);
