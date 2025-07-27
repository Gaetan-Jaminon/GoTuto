#!/bin/bash
set -e

echo "🧪 Running unit tests only..."

# Run unit tests without coverage for speed
go test -v ./internal/...

echo "✅ Unit tests completed successfully!"