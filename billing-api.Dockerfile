# Build stage - Using Red Hat UBI Go toolset
FROM registry.access.redhat.com/ubi9/go-toolset:1.20 AS builder

# Switch to root for package installation
USER 0

# Install ca-certificates
RUN dnf install -y ca-certificates && dnf clean all

# Set working directory
WORKDIR /opt/app-root/src

# Copy go mod files first for better Docker layer caching
COPY go.mod go.sum ./

# Download dependencies (cached layer if go.mod/go.sum unchanged)
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/billing-api

# Final stage - Using Red Hat UBI minimal
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

# Install ca-certificates for HTTPS
RUN microdnf install -y ca-certificates && microdnf clean all

# Create app directory
RUN mkdir -p /app/config

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /opt/app-root/src/main .

# Copy config files
COPY --from=builder /opt/app-root/src/config ./config

# Create non-root user (OpenShift compatible UID)
RUN useradd -u 1001 -r -g 0 -s /sbin/nologin appuser

# Set ownership (group 0 for OpenShift compatibility)
RUN chown -R 1001:0 /app && \
    chmod -R g=u /app

# Switch to non-root user
USER 1001

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Run the binary
CMD ["./main"]