package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// ErrInviteTokenNotFound is returned by InviteTokenRepository when no token
// matches the given hash.
var ErrInviteTokenNotFound = errors.New("invite token not found")

// InviteToken is a persisted record of an emailed invite-activation link,
// keyed by a hash of the raw token (never the raw value itself, same
// rationale as PasswordResetToken) so it can be looked up and consumed
// without the DB ever holding the value that was emailed. It gates
// WorkspaceMemberID's placeholder account (see WorkspaceUsecase.Invite) —
// only whoever received the link can activate that membership.
type InviteToken struct {
	ID                uuid.UUID
	WorkspaceMemberID uuid.UUID
	TokenHash         string
	ExpiresAt         time.Time
	UsedAt            *time.Time
	CreatedAt         time.Time
}

// IsValid reports whether the token is still usable — neither used nor
// expired.
func (t *InviteToken) IsValid() bool {
	return t.UsedAt == nil && time.Now().Before(t.ExpiresAt)
}

// InviteTokenRepository persists and retrieves InviteToken records.
type InviteTokenRepository interface {
	Create(ctx context.Context, token *InviteToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*InviteToken, error)
	MarkUsed(ctx context.Context, tokenHash string) error
}
