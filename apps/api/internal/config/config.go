// Package config loads apps/api's environment configuration.
// See apps/api/.env.example for the full list of supported keys.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all runtime settings the API needs. Fields are plain values
// (not pointers) since every setting is required to have a default or be
// resolved by Load — callers never need to distinguish "unset" from "zero".
type Config struct {
	Env  string
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
}

// Load reads configuration from a local .env file, falling back to real
// environment variables (e.g. in production, where no .env file exists).
// It returns an error instead of panicking so main can decide how to fail
// startup — idiomatic Go leaves error handling to the caller.
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	// AutomaticEnv lets real env vars override .env values, which is what
	// lets the same code run locally (via .env) and in containers (via
	// injected env vars) without a code change.
	v.AutomaticEnv()
	// Viper's internal key lookup uses "." as a nesting delimiter, but env
	// vars conventionally use "_"; this replacer bridges the two so
	// AutomaticEnv can match keys like DB_HOST correctly.
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetDefault("ENV", "development")
	v.SetDefault("PORT", "8080")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("JWT_ACCESS_EXPIRY", "15m")
	v.SetDefault("JWT_REFRESH_EXPIRY", "168h")

	if err := v.ReadInConfig(); err != nil {
		// A missing .env file is expected outside local dev (e.g. in
		// containers, where config comes from real env vars), so only other
		// read errors (e.g. malformed file) should fail startup. Since
		// SetConfigFile points at an explicit path rather than a searched
		// one, a missing file surfaces as a plain *fs.PathError here, not
		// viper.ConfigFileNotFoundError (that type is only returned by
		// viper's own search-path resolution).
		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("reading .env: %w", err)
		}
	}

	// Expiries are stored as duration strings (e.g. "15m") rather than raw
	// numbers so config files stay human-readable; parse them eagerly here
	// so any typo fails fast at startup instead of at first use.
	accessExpiry, err := time.ParseDuration(v.GetString("JWT_ACCESS_EXPIRY"))
	if err != nil {
		return nil, fmt.Errorf("parsing JWT_ACCESS_EXPIRY: %w", err)
	}
	refreshExpiry, err := time.ParseDuration(v.GetString("JWT_REFRESH_EXPIRY"))
	if err != nil {
		return nil, fmt.Errorf("parsing JWT_REFRESH_EXPIRY: %w", err)
	}

	return &Config{
		Env:  v.GetString("ENV"),
		Port: v.GetString("PORT"),

		DBHost:     v.GetString("DB_HOST"),
		DBPort:     v.GetString("DB_PORT"),
		DBUser:     v.GetString("DB_USER"),
		DBPassword: v.GetString("DB_PASSWORD"),
		DBName:     v.GetString("DB_NAME"),
		DBSSLMode:  v.GetString("DB_SSLMODE"),

		JWTAccessSecret:  v.GetString("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: v.GetString("JWT_REFRESH_SECRET"),
		JWTAccessExpiry:  accessExpiry,
		JWTRefreshExpiry: refreshExpiry,
	}, nil
}
