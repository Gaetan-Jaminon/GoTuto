#!/bin/bash
set -e

echo "🐳 Building Docker images locally..."

# Build billing API image
echo "Building billing-api Docker image..."
docker build -f billing-api.Dockerfile -t billing-api:local .

# Build catalog API image
echo "Building catalog-api Docker image..."
docker build -f catalog-api.Dockerfile -t catalog-api:local .

# Build billing migrator image
echo "Building billing-migrator Docker image..."
docker build -f billing-migrator.Dockerfile -t billing-migrator:local .

# Build catalog migrator image
echo "Building catalog-migrator Docker image..."
docker build -f catalog-migrator.Dockerfile -t catalog-migrator:local .

echo "✅ Docker images built successfully:"
docker images | grep -E "(billing-api|catalog-api|billing-migrator|catalog-migrator)" | grep local