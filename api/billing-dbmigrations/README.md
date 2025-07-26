# Billing Database Migrations

A dedicated module for managing database migrations for the billing service. This separation follows enterprise patterns where migrations are managed independently from the application runtime.

## Features

- **CLI-based migration management** using Cobra
- **Health check endpoint** for init containers
- **Configuration via file or environment variables**
- **Support for Podman/OpenShift deployment**
- **Rollback capabilities**
- **Migration creation helpers**

## Installation

```bash
cd api/billing-dbmigrations
go build -o billing-migrate cmd/migrate/*.go
```

## Usage

### Basic Commands

```bash
# Apply all pending migrations
./billing-migrate up

# Rollback one migration
./billing-migrate down

# Rollback N migrations
./billing-migrate down --steps=3

# Rollback all migrations
./billing-migrate down --all

# Check current version
./billing-migrate version

# Create new migration
./billing-migrate create add_customer_index

# Force version (for dirty state)
./billing-migrate force 5
```

### Configuration

#### Via Config File
```yaml
database:
  host: localhost
  port: 5432
  user: postgres
  password: secretpassword
  name: billing
  sslmode: require

migrations:
  path: ./migrations
```

#### Via Environment Variables
```bash
export BILLING_MIGRATE_DATABASE_HOST=localhost
export BILLING_MIGRATE_DATABASE_PORT=5432
export BILLING_MIGRATE_DATABASE_USER=postgres
export BILLING_MIGRATE_DATABASE_PASSWORD=secretpassword
export BILLING_MIGRATE_DATABASE_NAME=billing
export BILLING_MIGRATE_DATABASE_SSLMODE=require
```

#### Via Command Line Flags
```bash
./billing-migrate up \
  --db-host=localhost \
  --db-port=5432 \
  --db-user=postgres \
  --db-password=secretpassword \
  --db-name=billing \
  --db-sslmode=require
```

## Health Check Server

For Kubernetes/OpenShift init containers:

```bash
# Run health check server
go run cmd/health/main.go

# Endpoints:
# GET /health - Detailed health status
# GET /ready  - Simple readiness check
```

## Docker/Podman Build

### Multi-stage build for migrations
```dockerfile
FROM registry.access.redhat.com/ubi9/go-toolset:latest AS builder
WORKDIR /opt/app-root/src
COPY . .
RUN go build -o billing-migrate cmd/migrate/*.go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
COPY --from=builder /opt/app-root/src/billing-migrate /usr/local/bin/
COPY migrations /migrations
CMD ["billing-migrate", "up"]
```

## OpenShift Deployment

### Init Container Pattern
```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      initContainers:
      - name: migrate
        image: billing-migrations:latest
        command: ["billing-migrate", "up"]
        env:
        - name: BILLING_MIGRATE_DATABASE_HOST
          value: postgres-service
        - name: BILLING_MIGRATE_DATABASE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: postgres-secret
              key: password
```

### Job for One-time Migration
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: billing-db-migration
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: migrate
        image: billing-migrations:latest
        command: ["billing-migrate", "up"]
```

## Migration Files

Migrations follow the pattern: `{timestamp}_{description}.{up|down}.sql`

Example:
- `20240115143022_create_clients_table.up.sql`
- `20240115143022_create_clients_table.down.sql`

## Best Practices

1. **Always test migrations** in a non-production environment first
2. **Keep migrations idempotent** when possible
3. **Include rollback scripts** for every migration
4. **Use transactions** for multi-statement migrations
5. **Version control** all migration files
6. **Document breaking changes** in migration files

## Troubleshooting

### Dirty State
If migrations fail and leave the database in a dirty state:
```bash
# Check current state
./billing-migrate version

# Force to last known good version
./billing-migrate force 5

# Then reapply migrations
./billing-migrate up
```

### Connection Issues
```bash
# Test with explicit connection info
./billing-migrate version --db-host=localhost --db-port=5432
```

## Security Notes

- Never commit passwords to version control
- Use Kubernetes secrets or environment variables for sensitive data
- Consider using service accounts for database access in production
- Enable SSL/TLS for database connections in production