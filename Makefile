# =========================
# Project config
# =========================
include .env
export

APP_NAME=task-queue
CMD_DIR=cmd/api
BIN_DIR=bin
BIN=$(BIN_DIR)/$(APP_NAME)
MIGRATIONS_PATH=cmd/migrate/migrations
DB_URL?=$(DB_ADDR)

# =========================
# Go settings
# =========================
GO=go
GOFLAGS=-v

# =========================
# Default target
# =========================
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run        - Run app"
	@echo "  make dev        - Run with air"
	@echo "  make build      - Build binary"
	@echo "  make test       - Run tests"
	@echo "  make tidy       - Go mod tidy"
	@echo "  make fmt        - Go fmt"
	@echo "  make lint       - Go vet"
	@echo "  make clean      - Remove binaries"

# =========================
# Run app
# =========================
.PHONY: run
run:
	$(GO) run ./$(CMD_DIR)

# =========================
# Dev mode (Air)
# =========================
.PHONY: dev
dev:
	air -c .air.api.toml & air -c .air.worker.toml
echo:
	echo "hello"

# =========================
# Build binary
# =========================
.PHONY: build
build:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN) ./$(CMD_DIR)

# =========================
# Test
# =========================
.PHONY: test
test:
	$(GO) test ./... -count=1

# =========================
# Format
# =========================
.PHONY: fmt
fmt:
	$(GO) fmt ./...

# =========================
# Lint
# =========================
.PHONY: lint
lint:
	$(GO) vet ./...

# =========================
# Tidy
# =========================
.PHONY: tidy
tidy:
	$(GO) mod tidy

# =========================
# Clean
# =========================
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

MIGRATIONS_PATH=cmd/migrate/migrations
DB_URL?=$(DATABASE_URL)

.PHONY: migrate-up migrate-down migrate-force migrate-version migrate-create

## Apply all up migrations
migrate-up:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DB_URL)" up

## Rollback last migration
migrate-down:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DB_URL)" down 1

## Force migration version (used when dirty)
migrate-force:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DB_URL)" force $(v)

## Show current migration version
migrate-version:
	migrate -path=$(MIGRATIONS_PATH) -database="$(DB_URL)" version

## Create a new migration
## usage: make migrate-create name=add_tasks_table
migrate-create:
	migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(name)
