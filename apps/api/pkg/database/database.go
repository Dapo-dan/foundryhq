// Package database sets up the GORM connection to PostgreSQL.
package database

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Pool bounds are set conservatively for a single small Render instance
// against a Neon endpoint, which caps concurrent connections well below
// what Go's unbounded default would otherwise happily open under load.
const (
	maxOpenConns    = 15
	maxIdleConns    = 5
	connMaxLifetime = 30 * time.Minute
	connMaxIdleTime = 5 * time.Minute
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

	// GORM shares one *sql.DB under the hood, which defaults to unlimited
	// open connections — bounding it here keeps a connection burst (or two
	// Render instances briefly overlapping during a deploy) from exhausting
	// Neon's connection cap.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("unwrapping database handle: %w", err)
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	return db, nil
}
