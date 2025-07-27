#!/bin/bash
set -e

echo "🚀 Setting up local development environment..."

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.22 or later"
    exit 1
fi

echo "✅ Go version: $(go version)"

# Check Go version (requires 1.22+)
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [ "$(printf '%s\n' "1.22" "$GO_VERSION" | sort -V | head -n1)" != "1.22" ]; then
    echo "⚠️  Go 1.22+ recommended, you have $GO_VERSION"
fi

# Download dependencies
echo "📦 Downloading Go dependencies..."
go mod download
go mod verify

# Install useful development tools (optional)
echo "🔧 Installing development tools..."

# golangci-lint for linting
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2
fi

# Check Docker installation
if command -v docker &> /dev/null; then
    echo "✅ Docker is available: $(docker --version)"
else
    echo "⚠️  Docker not found - needed for container builds"
fi

# Make scripts executable
echo "🔧 Making scripts executable..."
chmod +x scripts/*.sh

# Test build
echo "🔨 Testing build..."
./scripts/build.sh

echo "✅ Development environment setup complete!"
echo ""
echo "Available commands:"
echo "  ./scripts/build.sh         - Build binaries"
echo "  ./scripts/test.sh          - Run tests with coverage"
echo "  ./scripts/test-unit.sh     - Run unit tests only"
echo "  ./scripts/lint.sh          - Run linting"
echo "  ./scripts/docker-build.sh  - Build Docker images"
echo "  ./scripts/clean.sh         - Clean build artifacts"
echo ""
echo "Or use the Makefile:"
echo "  make build"
echo "  make test"
echo "  make lint"
echo "  make docker"