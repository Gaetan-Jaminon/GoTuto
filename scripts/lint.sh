#!/bin/bash
set -e

echo "🔍 Running Go linting tools..."

# Run go fmt
echo "Running go fmt..."
if ! gofmt -l . | grep -q .; then
    echo "✅ Code is properly formatted"
else
    echo "❌ Code formatting issues found:"
    gofmt -l .
    echo "Run 'go fmt ./...' to fix formatting"
    exit 1
fi

# Run go vet
echo "Running go vet..."
go vet ./...

# Run golangci-lint if available
if command -v golangci-lint &> /dev/null; then
    echo "Running golangci-lint..."
    golangci-lint run
else
    echo "⚠️  golangci-lint not available, skipping"
    echo "Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2"
fi

echo "✅ Linting completed successfully!"