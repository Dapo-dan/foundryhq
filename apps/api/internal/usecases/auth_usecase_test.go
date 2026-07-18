package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/foundryhq/foundryhq/apps/api/internal/apperrors"
	"github.com/foundryhq/foundryhq/apps/api/internal/domain"
	"github.com/foundryhq/foundryhq/apps/api/pkg/jwt"
)

// mockUserRepo is a hand-written in-memory domain.UserRepository, per
// CONTRIBUTING.md's "mock repositories via interfaces" convention.
type mockUserRepo struct {
	byEmail map[string]*domain.User
	byID    map[uuid.UUID]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{byEmail: map[string]*domain.User{}, byID: map[uuid.UUID]*domain.User{}}
}

func (m *mockUserRepo) Create(_ context.Context, user *domain.User) error {
	if _, exists := m.byEmail[user.Email]; exists {
		return domain.ErrEmailAlreadyExists
	}
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	cp := *user
	m.byEmail[user.Email] = &cp
	m.byID[user.ID] = &cp
	return nil
}

func (m *mockUserRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.User, error) {
	u, ok := m.byID[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	cp := *u
	return &cp, nil
}

func (m *mockUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	u, ok := m.byEmail[email]
	if !ok {
		return nil, domain.ErrUserNotFound
	}
	cp := *u
	return &cp, nil
}

// mockRefreshTokenRepo is a hand-written in-memory domain.RefreshTokenRepository.
type mockRefreshTokenRepo struct {
	byHash map[string]*domain.RefreshToken
}

func newMockRefreshTokenRepo() *mockRefreshTokenRepo {
	return &mockRefreshTokenRepo{byHash: map[string]*domain.RefreshToken{}}
}

func (m *mockRefreshTokenRepo) Create(_ context.Context, token *domain.RefreshToken) error {
	token.ID = uuid.New()
	token.CreatedAt = time.Now()
	cp := *token
	m.byHash[token.TokenHash] = &cp
	return nil
}

func (m *mockRefreshTokenRepo) GetByTokenHash(_ context.Context, tokenHash string) (*domain.RefreshToken, error) {
	t, ok := m.byHash[tokenHash]
	if !ok {
		return nil, domain.ErrRefreshTokenNotFound
	}
	cp := *t
	return &cp, nil
}

func (m *mockRefreshTokenRepo) Revoke(_ context.Context, tokenHash string) error {
	t, ok := m.byHash[tokenHash]
	if !ok {
		return nil
	}
	now := time.Now()
	t.RevokedAt = &now
	return nil
}

func newTestAuthUsecase() (*AuthUsecase, *mockUserRepo, *mockRefreshTokenRepo) {
	users := newMockUserRepo()
	tokens := newMockRefreshTokenRepo()
	manager := jwt.NewManager("access-secret", "refresh-secret", 15*time.Minute, 168*time.Hour)
	return NewAuthUsecase(users, tokens, manager), users, tokens
}

func asAppError(t *testing.T, err error) *apperrors.Error {
	t.Helper()
	var appErr *apperrors.Error
	if !errors.As(err, &appErr) {
		t.Fatalf("error = %v, want an *apperrors.Error", err)
	}
	return appErr
}

func TestRegister_Success(t *testing.T) {
	u, _, tokens := newTestAuthUsecase()

	result, err := u.Register(context.Background(), "User@Example.com", "password123")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if result.User.Email != "user@example.com" {
		t.Errorf("User.Email = %q, want lowercased/trimmed email", result.User.Email)
	}
	if result.AccessToken == "" || result.RefreshToken == "" {
		t.Error("Register() should return both an access and a refresh token")
	}
	if len(tokens.byHash) != 1 {
		t.Errorf("expected 1 persisted refresh token, got %d", len(tokens.byHash))
	}
}

func TestRegister_ValidationErrors(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		field    string
	}{
		{"empty email", "  ", "password123", "email"},
		{"short password", "user@example.com", "short", "password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, _, _ := newTestAuthUsecase()
			_, err := u.Register(context.Background(), tt.email, tt.password)
			appErr := asAppError(t, err)
			if appErr.Code != apperrors.CodeValidation {
				t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeValidation)
			}
			if appErr.Field != tt.field {
				t.Errorf("Field = %q, want %q", appErr.Field, tt.field)
			}
		})
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	u, _, _ := newTestAuthUsecase()
	ctx := context.Background()

	if _, err := u.Register(ctx, "user@example.com", "password123"); err != nil {
		t.Fatalf("first Register() error = %v", err)
	}

	_, err := u.Register(ctx, "user@example.com", "password456")
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeConflict {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeConflict)
	}
}

func TestLogin_Success(t *testing.T) {
	u, _, _ := newTestAuthUsecase()
	ctx := context.Background()

	registered, err := u.Register(ctx, "user@example.com", "password123")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	result, err := u.Login(ctx, "user@example.com", "password123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if result.User.ID != registered.User.ID {
		t.Errorf("Login() User.ID = %v, want %v", result.User.ID, registered.User.ID)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	u, _, _ := newTestAuthUsecase()
	ctx := context.Background()

	if _, err := u.Register(ctx, "user@example.com", "password123"); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	_, err := u.Login(ctx, "user@example.com", "wrong-password")
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeUnauthorized {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeUnauthorized)
	}
	if appErr.Message != invalidCredentialsMessage {
		t.Errorf("Message = %q, want the generic invalid-credentials message", appErr.Message)
	}
}

func TestLogin_UnknownEmail(t *testing.T) {
	u, _, _ := newTestAuthUsecase()

	_, err := u.Login(context.Background(), "nobody@example.com", "password123")
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeUnauthorized {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeUnauthorized)
	}
	// Deliberately the same message as a wrong password, so login can't be
	// used to enumerate registered emails.
	if appErr.Message != invalidCredentialsMessage {
		t.Errorf("Message = %q, want the generic invalid-credentials message", appErr.Message)
	}
}

func TestRefresh_Success(t *testing.T) {
	u, _, _ := newTestAuthUsecase()
	ctx := context.Background()

	registered, err := u.Register(ctx, "user@example.com", "password123")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	refreshed, err := u.Refresh(ctx, registered.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh() error = %v", err)
	}
	if refreshed.User.ID != registered.User.ID {
		t.Errorf("Refresh() User.ID = %v, want %v", refreshed.User.ID, registered.User.ID)
	}
	if refreshed.RefreshToken == registered.RefreshToken {
		t.Error("Refresh() should rotate to a new refresh token, not reuse the old one")
	}

	// The rotated-out token must no longer work.
	if _, err := u.Refresh(ctx, registered.RefreshToken); err == nil {
		t.Error("Refresh() with the rotated-out token should fail")
	}
}

func TestRefresh_InvalidToken(t *testing.T) {
	u, _, _ := newTestAuthUsecase()

	_, err := u.Refresh(context.Background(), "not-a-real-token")
	appErr := asAppError(t, err)
	if appErr.Code != apperrors.CodeUnauthorized {
		t.Errorf("Code = %v, want %v", appErr.Code, apperrors.CodeUnauthorized)
	}
}

func TestLogout_RevokesToken(t *testing.T) {
	u, _, _ := newTestAuthUsecase()
	ctx := context.Background()

	registered, err := u.Register(ctx, "user@example.com", "password123")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if err := u.Logout(ctx, registered.RefreshToken); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}

	if _, err := u.Refresh(ctx, registered.RefreshToken); err == nil {
		t.Error("Refresh() with a logged-out token should fail")
	}
}

func TestLogout_UnknownTokenIsNotAnError(t *testing.T) {
	u, _, _ := newTestAuthUsecase()

	if err := u.Logout(context.Background(), "never-issued"); err != nil {
		t.Errorf("Logout() error = %v, want nil for an already-gone session", err)
	}
}
