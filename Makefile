# Makefile

# Load .env file
include .env
export $(shell sed 's/=.*//' .env)

# Load environment variables
DB_NAME ?= $(POSTGRES_DB)
DB_USER ?= $(POSTGRES_USER)
DB_PASSWORD ?= $(POSTGRES_PASSWORD)
DB_HOST ?= $(POSTGRES_HOST)
DB_PORT ?= $(POSTGRES_PORT)
DB_SSLMODE ?= disable

MIGRATIONS_DIR = ./db/migrations
DATABASE_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
export $PATH="bin/:$PATH"
# Migrate Up: Apply all migrations
migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

# Migrate Down: Revert the last migration
migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down

# Migrate Force: Force a specific migration version
migrate-force:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)

# Migrate Create: Create a new migration file
migrate-create:
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

# Help: Display help for Makefile targets
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  migrate-up       Apply all migrations"
	@echo "  migrate-down     Revert the last migration"
	@echo "  migrate-force    Force a specific migration version (e.g., make migrate-force version=1)"
	@echo "  migrate-create   Create a new migration file (e.g., make migrate-create name=create_users_table)"
	@echo "  help             Show this help message"

.PHONY: migrate-up migrate-down migrate-create migrate-force
