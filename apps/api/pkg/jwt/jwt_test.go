package jwt

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func testManager() *Manager {
	return NewManager("access-secret", "refresh-secret", 15*time.Minute, 168*time.Hour)
}

func TestAccessTokenRoundTrip(t *testing.T) {
	m := testManager()
	userID := uuid.New()

	token, err := m.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	claims, err := m.ValidateAccessToken(token)
	if err != nil {
		t.Fatalf("ValidateAccessToken() error = %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.TokenType != tokenTypeAccess {
		t.Errorf("TokenType = %q, want %q", claims.TokenType, tokenTypeAccess)
	}
}

func TestRefreshTokenRoundTrip(t *testing.T) {
	m := testManager()
	userID := uuid.New()

	token, expiresAt, err := m.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}
	if !expiresAt.After(time.Now()) {
		t.Errorf("expiresAt = %v, want a future time", expiresAt)
	}

	claims, err := m.ValidateRefreshToken(token)
	if err != nil {
		t.Fatalf("ValidateRefreshToken() error = %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
}

func TestValidateAccessToken_Expired(t *testing.T) {
	m := NewManager("access-secret", "refresh-secret", -1*time.Minute, 168*time.Hour)
	token, err := m.GenerateAccessToken(uuid.New())
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	if _, err := m.ValidateAccessToken(token); err != ErrInvalidToken {
		t.Errorf("ValidateAccessToken() error = %v, want %v", err, ErrInvalidToken)
	}
}

func TestValidateAccessToken_WrongSecret(t *testing.T) {
	issuerManager := testManager()
	token, err := issuerManager.GenerateAccessToken(uuid.New())
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	verifier := NewManager("different-secret", "refresh-secret", 15*time.Minute, 168*time.Hour)
	if _, err := verifier.ValidateAccessToken(token); err != ErrInvalidToken {
		t.Errorf("ValidateAccessToken() error = %v, want %v", err, ErrInvalidToken)
	}
}

func TestValidateAccessToken_TamperedSignature(t *testing.T) {
	m := testManager()
	token, err := m.GenerateAccessToken(uuid.New())
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}

	tampered := token[:len(token)-1] + "x"
	if _, err := m.ValidateAccessToken(tampered); err != ErrInvalidToken {
		t.Errorf("ValidateAccessToken() error = %v, want %v", err, ErrInvalidToken)
	}
}

// TestTokenTypeMismatchRejected uses the same secret for both access and
// refresh so the TokenType check is what's actually being exercised, not
// just a signature mismatch from the two Config secrets normally differing.
func TestTokenTypeMismatchRejected(t *testing.T) {
	m := NewManager("shared-secret", "shared-secret", 15*time.Minute, 168*time.Hour)

	accessToken, err := m.GenerateAccessToken(uuid.New())
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	refreshToken, _, err := m.GenerateRefreshToken(uuid.New())
	if err != nil {
		t.Fatalf("GenerateRefreshToken() error = %v", err)
	}

	if _, err := m.ValidateRefreshToken(accessToken); err != ErrInvalidToken {
		t.Errorf("ValidateRefreshToken(accessToken) error = %v, want %v", err, ErrInvalidToken)
	}
	if _, err := m.ValidateAccessToken(refreshToken); err != ErrInvalidToken {
		t.Errorf("ValidateAccessToken(refreshToken) error = %v, want %v", err, ErrInvalidToken)
	}
}
