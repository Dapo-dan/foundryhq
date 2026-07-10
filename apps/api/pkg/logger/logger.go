// Package logger builds the application's structured logger.
package logger

import "go.uber.org/zap"

// New builds a zap.Logger suited to env. zap ships two presets so callers
// don't need to hand-tune log config: NewProduction emits JSON for log
// aggregators, NewDevelopment emits human-readable, colorized output.
func New(env string) (*zap.Logger, error) {
	if env == "production" {
		return zap.NewProduction()
	}
	return zap.NewDevelopment()
}
