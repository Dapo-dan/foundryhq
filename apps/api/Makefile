# Requires the golang-migrate CLI: https://github.com/golang-migrate/migrate
#
#   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
#
# Reads DB_* from apps/api/.env (see .env.example), matching the connection
# details config.go loads for the running server.

include .env
export

DATABASE_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATIONS_PATH := internal/repositories/postgres/migrations

.PHONY: migrate-up migrate-down

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1
