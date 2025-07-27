.PHONY: build test test-unit lint docker clean dev-setup help

# Default target
help:
	@echo "Available commands:"
	@echo "  make build        - Build all binaries (APIs and migrators)"
	@echo "  make test         - Run tests with coverage"
	@echo "  make test-unit    - Run unit tests only"
	@echo "  make lint         - Run linting tools"
	@echo "  make docker       - Build Docker images"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make dev-setup    - Set up development environment"
	@echo ""
	@echo "Or use scripts directly:"
	@echo "  ./scripts/build.sh"
	@echo "  ./scripts/test.sh"
	@echo "  ./scripts/lint.sh"

build:
	@./scripts/build.sh

test:
	@./scripts/test.sh

test-unit:
	@./scripts/test-unit.sh

lint:
	@./scripts/lint.sh

docker:
	@./scripts/docker-build.sh

clean:
	@./scripts/clean.sh

dev-setup:
	@./scripts/dev-setup.sh