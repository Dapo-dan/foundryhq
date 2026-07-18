package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrUserNotFound is returned by UserRepository when no user matches the
	// given lookup.
	ErrUserNotFound = errors.New("user not found")
	// ErrEmailAlreadyExists is returned by UserRepository.Create when the
	// email is already registered.
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// User is a registered account. OAuth-only users (not built until v1.1+ per
// docs/mvp.md) will have a nil PasswordHash — every user created by the
// current email/password registration flow always sets one.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// UserRepository persists and retrieves User entities.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
}
