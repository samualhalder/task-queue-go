# =========================
# Project config
# =========================
APP_NAME=task-queue
CMD_DIR=cmd/api
BIN_DIR=bin
BIN=$(BIN_DIR)/$(APP_NAME)

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
	air
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
