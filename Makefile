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
local:
	docker compose up url_shortener_db
	go run main.go

.PHONY: migrate-up migrate-down migrate-create migrate-force local
