# Demo01 API Makefile

# Variables
BINARY_NAME=demo01-api
MIGRATE_BINARY=migrate-tool

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the application
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) cmd/main.go
	$(GOBUILD) -o bin/$(MIGRATE_BINARY) cmd/migrate/main.go

# Run the application
run:
	$(GOCMD) run cmd/main.go

# Run migrations
migrate-up:
	$(GOCMD) run cmd/migrate/main.go up

migrate-down:
	$(GOCMD) run cmd/migrate/main.go down

migrate-version:
	$(GOCMD) run cmd/migrate/main.go version

# Test the application
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin/

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

# Database commands (requires running PostgreSQL)
db-create:
	createdb demo01

db-drop:
	dropdb demo01

db-reset: db-drop db-create migrate-up

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  migrate-up    - Apply database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  migrate-version - Show current migration version"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download dependencies"
	@echo "  dev           - Run with hot reload (requires air)"
	@echo "  db-create     - Create database"
	@echo "  db-drop       - Drop database"
	@echo "  db-reset      - Reset database (drop, create, migrate)"
	@echo "  help          - Show this help"

.PHONY: build run migrate-up migrate-down migrate-version test clean deps dev docker-build docker-run db-create db-drop db-reset help