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

# Build both binaries
RUN CGO_ENABLED=0 GOOS=linux go build -a -o billing-migrate cmd/migrate/*.go && \
    CGO_ENABLED=0 GOOS=linux go build -a -o billing-health cmd/health/main.go

# Final stage - Using Red Hat UBI minimal
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

# Install PostgreSQL client for debugging and ca-certificates
RUN microdnf install -y postgresql ca-certificates && microdnf clean all

WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /opt/app-root/src/billing-migrate /usr/local/bin/
COPY --from=builder /opt/app-root/src/billing-health /usr/local/bin/

# Copy migrations and config
COPY migrations /app/migrations
COPY config.yaml /app/config.yaml

# Create non-root user (OpenShift compatible UID)
RUN useradd -u 1001 -r -g 0 -s /sbin/nologin migrator

# Set ownership (group 0 for OpenShift compatibility)
RUN chown -R 1001:0 /app && \
    chmod -R g=u /app && \
    chmod +x /usr/local/bin/billing-migrate /usr/local/bin/billing-health

# Switch to non-root user
USER 1001

# Health check for migration status
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD billing-health || exit 1

# Default command runs migrations
CMD ["billing-migrate", "up"]