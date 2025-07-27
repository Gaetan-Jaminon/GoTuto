#!/bin/bash
set -e

echo "🐳 Building Docker images locally..."

# Build billing API image
echo "Building billing-api Docker image..."
docker build -f billing-api.Dockerfile -t billing-api:local .

# Build billing migrator image
echo "Building billing-migrator Docker image..."
docker build -f billing-migrator.Dockerfile -t billing-migrator:local .

echo "✅ Docker images built successfully:"
docker images | grep -E "(billing-api|billing-migrator)" | grep local