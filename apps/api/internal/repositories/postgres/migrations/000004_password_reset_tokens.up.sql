-- Backs AuthUsecase.ForgotPassword/ResetPassword (Phase 0 auth close-out).
-- Same shape as refresh_tokens: stores a SHA-256 hash of the emailed token,
-- never the raw value, so a DB read alone can't be used to reset a
-- password. used_at (rather than a hard delete) lets a reused/expired
-- token still be distinguished from one that was never issued.

CREATE TABLE password_reset_tokens (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token_hash text NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    used_at    timestamptz,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_password_reset_tokens_user_id ON password_reset_tokens (user_id);
