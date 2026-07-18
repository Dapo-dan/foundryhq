// Package database sets up the GORM connection to PostgreSQL.
package database

import (
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config describes the connection parameters needed to reach Postgres.
// It's a separate type from internal/config.Config so this package doesn't
// depend on the whole app's config shape — it only needs these six fields.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// Connect opens a GORM connection to Postgres built from cfg.
func Connect(cfg Config) (*gorm.DB, error) {
	// Building the DSN via url.URL (instead of fmt.Sprintf) lets the
	// stdlib handle escaping of special characters in the user/password,
	// avoiding a malformed or unsafely-injected connection string.
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:     "/" + cfg.Name,
		RawQuery: url.Values{"sslmode": {cfg.SSLMode}}.Encode(),
	}

	// TranslateError converts driver-specific errors (e.g. Postgres' unique
	// violation code) into GORM's portable sentinel errors (gorm.ErrDuplicatedKey),
	// so repositories can use errors.Is without depending on the pgx driver's
	// error types directly.
	db, err := gorm.Open(postgres.Open(dsn.String()), &gorm.Config{TranslateError: true})
	if err != nil {
		// %w wraps the original error so callers can still use errors.Is/As
		// on it, while adding context about what operation failed.
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return db, nil
}
