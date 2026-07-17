// Package jwt issues and validates the access/refresh token pair described in
// adr/0004-jwt-access-refresh-tokens.md. Tokens carry only the user's
// identity (sub) — workspace-scoped authorization is resolved per-request
// against workspace_members, not embedded in the token, since a user's
// workspace memberships can change without requiring a new login.
package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const issuer = "foundryhq-api"

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

// ErrInvalidToken is returned for any validation failure — expired,
// malformed, wrong signature, or wrong token type presented for the
// expected use. The underlying jwt library's parse error is deliberately
// not exposed, since callers (e.g. middleware) forward this straight into
// an HTTP error response.
var ErrInvalidToken = errors.New("invalid token")

// Claims is the payload carried by both access and refresh tokens. TokenType
// distinguishes the two so a refresh token can't be replayed against an
// endpoint expecting an access token (or vice versa) purely by shape, as
// defense-in-depth on top of the two token types already using distinct
// signing secrets.
type Claims struct {
	UserID    uuid.UUID `json:"sub"`
	TokenType string    `json:"type"`
	jwt.RegisteredClaims
}

// Manager generates and validates access/refresh tokens using the secrets
// and expiries loaded from Config (see internal/config). Construct one
// instance in main and share it across the auth middleware and, later, the
// auth usecase.
type Manager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewManager constructs a Manager from the given secrets and expiries.
func NewManager(accessSecret, refreshSecret string, accessExpiry, refreshExpiry time.Duration) *Manager {
	return &Manager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateAccessToken issues a short-lived access token for userID.
func (m *Manager) GenerateAccessToken(userID uuid.UUID) (string, error) {
	return m.generate(userID, tokenTypeAccess, m.accessSecret, m.accessExpiry)
}

// GenerateRefreshToken issues a long-lived refresh token for userID,
// returning its raw string and expiry so the caller (a future auth
// usecase) can persist a hash of it for logout invalidation — Manager
// itself has no storage and never sees the persisted form.
func (m *Manager) GenerateRefreshToken(userID uuid.UUID) (token string, expiresAt time.Time, err error) {
	expiresAt = time.Now().Add(m.refreshExpiry)
	token, err = m.generate(userID, tokenTypeRefresh, m.refreshSecret, m.refreshExpiry)
	if err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}

// ValidateAccessToken parses and verifies tokenString as an access token.
func (m *Manager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return m.validate(tokenString, tokenTypeAccess, m.accessSecret)
}

// ValidateRefreshToken parses and verifies tokenString as a refresh token.
func (m *Manager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return m.validate(tokenString, tokenTypeRefresh, m.refreshSecret)
}

func (m *Manager) generate(userID uuid.UUID, tokenType string, secret []byte, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:    userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (m *Manager) validate(tokenString, expectedType string, secret []byte) (*Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if claims.TokenType != expectedType {
		return nil, ErrInvalidToken
	}
	return &claims, nil
}
