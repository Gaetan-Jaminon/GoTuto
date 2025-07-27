#!/bin/bash
set -e

echo "🧪 Running all tests..."

# Run unit tests with coverage
echo "Running unit tests with coverage..."
go test -v -race -coverprofile=coverage.out ./internal/...

# Generate coverage report
echo "Generating coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Show coverage summary
echo "Coverage summary:"
go tool cover -func=coverage.out

echo "✅ Tests completed successfully!"
echo "📊 Coverage report generated: coverage.html"