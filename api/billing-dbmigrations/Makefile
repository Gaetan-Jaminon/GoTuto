# Billing Database Migrations Makefile

BINARY_NAME=billing-migrate
HEALTH_BINARY=billing-health

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the migration tool
build:
	$(GOBUILD) -o bin/$(BINARY_NAME) cmd/migrate/*.go
	$(GOBUILD) -o bin/$(HEALTH_BINARY) cmd/health/*.go

# Build for Linux (useful for containers)
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(BINARY_NAME)-linux cmd/migrate/*.go
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o bin/$(HEALTH_BINARY)-linux cmd/health/*.go

# Run migrations up
up:
	$(GOCMD) run cmd/migrate/*.go up

# Run migrations down
down:
	$(GOCMD) run cmd/migrate/*.go down

# Check migration version
version:
	$(GOCMD) run cmd/migrate/*.go version

# Create a new migration
create:
	@read -p "Enter migration name: " name; \
	$(GOCMD) run cmd/migrate/*.go create $$name

# Run health check server
health:
	$(GOCMD) run cmd/health/main.go

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin/

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Test (when we add tests)
test:
	$(GOTEST) -v ./...

# Container build with Podman
container-build:
	podman build -t billing-migrations:latest -f Containerfile .

# Run container with Podman
container-run:
	podman run --rm \
		--network host \
		-e BILLING_MIGRATE_DATABASE_HOST=localhost \
		-e BILLING_MIGRATE_DATABASE_PASSWORD=$${DB_PASSWORD} \
		billing-migrations:latest up

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the migration tools"
	@echo "  build-linux    - Build for Linux (container deployment)"
	@echo "  up             - Apply all pending migrations"
	@echo "  down           - Rollback last migration"
	@echo "  version        - Show current migration version"
	@echo "  create         - Create a new migration"
	@echo "  health         - Run health check server"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Download dependencies"
	@echo "  test           - Run tests"
	@echo "  container-build - Build container with Podman"
	@echo "  container-run  - Run migrations in container"
	@echo "  help           - Show this help"

.PHONY: build build-linux up down version create health clean deps test container-build container-run help