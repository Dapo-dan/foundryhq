package postgres

import (
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/foundryhq/foundryhq/apps/api/pkg/database"
)

// openTestDB connects to the dev Postgres (see docker-compose.yml), reading
// connection details from the environment with the docker-compose defaults
// as fallback. Tests skip rather than fail when Postgres isn't reachable,
// per CONTRIBUTING.md's "integration tests require a running DB".
func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := database.Connect(database.Config{
		Host:     envOrDefault("DB_HOST", "localhost"),
		Port:     envOrDefault("DB_PORT", "5434"),
		User:     envOrDefault("DB_USER", "foundryhq"),
		Password: envOrDefault("DB_PASSWORD", "foundryhq"),
		Name:     envOrDefault("DB_NAME", "foundryhq"),
		SSLMode:  envOrDefault("DB_SSLMODE", "disable"),
	})
	if err != nil {
		t.Skipf("postgres not reachable, skipping integration test: %v", err)
	}
	return db
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// withTestTx runs fn against a transaction that's always rolled back
// afterward, so integration tests never leave rows behind in the shared
// dev database.
func withTestTx(t *testing.T, fn func(tx *gorm.DB)) {
	t.Helper()
	db := openTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()
	fn(tx)
}
