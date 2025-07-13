# Makefile

-include .env
export
#
MIGRATIONS_DIR = "./db/postgres/migrations"
DATABASE_URL = "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)"
# export $PATH="${PWD}/bin/:$PATH"

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

db-up: 
	docker compose up url_shortener_db dragonfly -d
	sleep 10

local: db-up migrate-up
	air

# dev: implemented future
# 	docker compose up --build -d
clean:
	docker compose down
.PHONY: migrate-up migrate-down migrate-create migrate-force db-up local dev clean
