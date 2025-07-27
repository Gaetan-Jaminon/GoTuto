# Build stage - Using Red Hat UBI Go toolset
FROM registry.access.redhat.com/ubi9/go-toolset:1.20 AS builder

# Switch to root for package installation
USER 0

# Set working directory
WORKDIR /opt/app-root/src

# Copy go mod files first for better Docker layer caching
COPY go.mod go.sum ./

# Download dependencies (cached layer if go.mod/go.sum unchanged)
RUN go mod download

# Copy source code
COPY . .

# Build the migrator binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o catalog-migrator ./cmd/catalog-migrator

# Final stage - Using Red Hat UBI minimal
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

# Install PostgreSQL client for debugging and ca-certificates
RUN microdnf install -y postgresql ca-certificates && microdnf clean all

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /opt/app-root/src/catalog-migrator /usr/local/bin/

# Copy catalog migrations and config
COPY --from=builder /opt/app-root/src/internal/catalog/migrations /app/migrations
COPY --from=builder /opt/app-root/src/config ./config

# Create non-root user (OpenShift compatible UID)
RUN useradd -u 1001 -r -g 0 -s /sbin/nologin migrator

# Set ownership (group 0 for OpenShift compatibility)
RUN chown -R 1001:0 /app && \
    chmod -R g=u /app && \
    chmod +x /usr/local/bin/catalog-migrator

# Switch to non-root user
USER 1001

# Default command runs catalog migrations
CMD ["catalog-migrator", "up"]