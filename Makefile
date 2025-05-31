# Makefile

-include .env
export

# Load environment variables
DB_NAME ?= $(DB_NAME)
DB_USER ?= $(DB_USER)
DB_PASSWORD ?= $(DB_PASSWORD)
DB_HOST ?= $(DB_HOST)
DB_PORT ?= $(DB_PORT)
DB_SSLMODE ?= ${DB_SSLMODE}
DB_HOST ?= ${DB_HOST}
MIGRATIONS_DIR = "./db/postgres/migrations"
DATABASE_URL = "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"
export $PATH="${PWD}/bin/:$PATH"

OS := $(shell uname | tr A-Z a-z)
ARCH := $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')

# .PHONY: install-deps deps
# install-deps: migrate air mockery golangci-lint
# deps: $(MIGRATE) $(AIR) ${MOKERY} $(GOLANGCI) ## Checks for Global Development Dependencies.
# # Migrate Up: Apply all migrations
migrate-up:
	echo $(DATABASE_URL)
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
	@echo "migarte	Download latest verison of migrate golang-migrate"
	@echo "  migrate-up       Apply all migrations"
	@echo "  migrate-down     Revert the last migration"
	@echo "  migrate-force    Force a specific migration version (e.g., make migrate-force version=1)"
	@echo "  migrate-create   Create a new migration file (e.g., make migrate-create name=create_users_table)"
	@echo "  help             Show this help message"

.PHONY: migrate-up migrate-down migrate-create migrate-force
