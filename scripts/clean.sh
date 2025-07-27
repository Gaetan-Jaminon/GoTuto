#!/bin/bash

echo "🧹 Cleaning build artifacts..."

# Remove binaries
if [ -d "bin" ]; then
    echo "Removing binaries..."
    rm -rf bin/
fi

# Remove coverage files
if [ -f "coverage.out" ]; then
    echo "Removing coverage files..."
    rm -f coverage.out coverage.html
fi

# Remove Go build cache (optional)
if [ "$1" == "--all" ]; then
    echo "Cleaning Go build cache..."
    go clean -cache
    go clean -modcache
fi

# Remove Docker images if they exist
if docker images | grep -q "billing.*local"; then
    echo "Removing local Docker images..."
    docker images | grep "billing.*local" | awk '{print $3}' | xargs docker rmi --force 2>/dev/null || true
fi

echo "✅ Cleanup completed!"