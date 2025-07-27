#!/bin/bash
set -e

echo "🔨 Building Go binaries..."

# Create bin directory if it doesn't exist
mkdir -p bin/

# Build all applications
echo "Building billing API..."
go build -o bin/billing-api ./cmd/billing-api

echo "Building catalog API..."
go build -o bin/catalog-api ./cmd/catalog-api

echo "Building billing migrator..."
go build -o bin/billing-migrator ./cmd/billing-migrator

echo "Building catalog migrator..."
go build -o bin/catalog-migrator ./cmd/catalog-migrator

echo "✅ Binaries built successfully:"
ls -la bin/