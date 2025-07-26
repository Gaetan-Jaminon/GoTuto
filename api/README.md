# GoTuto API Services

This directory contains Go-based microservices for learning enterprise Go development patterns. The project demonstrates real-world practices including CRUD operations, database migrations, comprehensive testing, and CI/CD workflows.

## 📁 Project Structure

```
api/
├── billing/                    # Main API service
│   ├── cmd/main.go            # Application entry point
│   ├── internal/              # Private application code
│   │   ├── config/            # Configuration management
│   │   ├── database/          # Database connection
│   │   ├── handlers/          # HTTP handlers
│   │   └── models/            # Data models
│   ├── tests/                 # Test suites
│   │   ├── integration/       # Integration tests
│   │   └── e2e/              # End-to-end tests
│   ├── deployments/          # Deployment manifests
│   └── Dockerfile            # Container build
├── billing-dbmigrations/      # Database migration service
│   ├── cmd/                  # CLI tools (migrate, health)
│   ├── migrations/           # SQL migration files
│   └── Dockerfile           # Container build
└── README.md                # This file
```

## 🚀 Quick Start

### Prerequisites

- **Go 1.22+**
- **PostgreSQL 15+**
- **Podman** (for local development)
- **Git**

### Local Development Setup

1. **Clone and navigate**:
   ```bash
   git clone https://github.com/Gaetan-Jaminon/GoTuto.git
   cd GoTuto/api
   ```

2. **Start PostgreSQL** (using Podman):
   ```bash
   # Start PostgreSQL container
   podman run -d \
     --name postgres-dev \
     -e POSTGRES_DB=billing \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=password \
     -p 5432:5432 \
     postgres:15

   # Create test database
   podman exec postgres-dev createdb -U postgres billing_test
   ```

3. **Run database migrations**:
   ```bash
   cd billing-dbmigrations
   go run cmd/migrate/*.go up
   ```

4. **Start the API service**:
   ```bash
   cd ../billing
   go run cmd/main.go
   ```

5. **Test the API**:
   ```bash
   curl http://localhost:8080/health
   curl http://localhost:8080/api/v1/clients
   ```

## 🏗️ Services Overview

### Billing Service (`/billing`)

A RESTful API service providing CRUD operations for clients and invoices.

**Key Features:**
- REST API endpoints for clients and invoices
- Hierarchical configuration management (Viper)
- Structured logging
- Health checks and metrics
- Comprehensive error handling
- Input validation

**API Endpoints:**
```
GET    /health                    # Health check
GET    /api/v1/clients           # List clients
POST   /api/v1/clients           # Create client
GET    /api/v1/clients/{id}      # Get client
PUT    /api/v1/clients/{id}      # Update client
DELETE /api/v1/clients/{id}      # Delete client
GET    /api/v1/invoices          # List invoices
POST   /api/v1/invoices          # Create invoice
GET    /api/v1/invoices/{id}     # Get invoice
PUT    /api/v1/invoices/{id}     # Update invoice
DELETE /api/v1/invoices/{id}     # Delete invoice
```

### Database Migrations Service (`/billing-dbmigrations`)

A dedicated service for managing database schema changes with CLI tools.

**Key Features:**
- Database migration management (up/down)
- Migration status tracking
- Health check server
- Cobra CLI interface
- PostgreSQL support with SSL

**CLI Commands:**
```bash
billing-migrate up           # Apply all pending migrations
billing-migrate down         # Rollback one migration
billing-migrate force <v>    # Force set version without running migrations
billing-migrate version      # Show current migration version
billing-migrate create <name> # Create new migration files
billing-health              # Run health check server
```

## 🧪 Testing Strategy

The project implements comprehensive testing following Go best practices:

### Test Types

1. **Unit Tests** (`*_test.go`)
   - Model validation
   - Handler logic
   - Configuration loading
   - Table-driven tests

2. **Integration Tests** (`tests/integration/`)
   - Database interactions
   - Full request/response cycles
   - Service integration

3. **End-to-End Tests** (`tests/e2e/`)
   - Complete workflow testing
   - External service interactions
   - Production-like scenarios

### Running Tests

```bash
# Unit tests
cd billing
go test ./...

# Integration tests (requires database)
go test -tags=integration ./tests/integration/...

# E2E tests (requires running service)
go test -tags=e2e ./tests/e2e/...

# All tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 🐳 Container Images

Both services are containerized using multi-stage builds with Red Hat UBI base images for OpenShift compatibility.

### Build Images Locally

```bash
# Build billing service
cd billing
podman build -t billing:latest .

# Build migrations service
cd ../billing-dbmigrations
podman build -t billing-dbmigrations:latest .
```

### Run with Podman

```bash
# Run migrations
podman run --rm \
  --network host \
  -e APP_ENV=dev \
  -e BILLING_MIGRATE_DATABASE_HOST=localhost \
  -e BILLING_MIGRATE_DATABASE_PASSWORD=password \
  billing-dbmigrations:latest

# Run API service
podman run -d \
  --network host \
  --name billing-api \
  -e APP_ENV=dev \
  -e BILLING_DATABASE_HOST=localhost \
  -e BILLING_DATABASE_PASSWORD=password \
  billing:latest
```

## ⚙️ Configuration

Configuration is managed hierarchically using Viper:

1. **Default values** (in code)
2. **Configuration files** (`config/config.yaml`)
3. **Environment variables** (prefixed with `BILLING_`)
4. **Command line flags**

### Environment Variables

**Billing Service:**
```bash
BILLING_SERVER_PORT=8080
BILLING_DATABASE_HOST=localhost
BILLING_DATABASE_PORT=5432
BILLING_DATABASE_USERNAME=postgres
BILLING_DATABASE_PASSWORD=password
BILLING_DATABASE_NAME=billing
BILLING_DATABASE_SSL_MODE=disable
BILLING_LOG_LEVEL=info
```

**Migrations Service:**
```bash
BILLING_MIGRATE_DATABASE_HOST=localhost
BILLING_MIGRATE_DATABASE_PORT=5432
BILLING_MIGRATE_DATABASE_USER=postgres
BILLING_MIGRATE_DATABASE_PASSWORD=password
BILLING_MIGRATE_DATABASE_NAME=billing
BILLING_MIGRATE_DATABASE_SSLMODE=disable
```

## 🚀 CI/CD Pipeline

The project includes comprehensive GitHub Actions workflows:

### Workflows

1. **Continuous Integration** (`.github/workflows/ci.yml`)
   - Multi-module build strategy
   - Lint and test with PostgreSQL
   - Security scanning with Trivy
   - Container image building
   - Quality gates

2. **Continuous Deployment** (`.github/workflows/cd.yml`)
   - Staging deployment after CI
   - E2E tests against staging
   - Production deployment with approval
   - OpenShift integration

3. **Release Management** (`.github/workflows/release.yml`)
   - Semantic versioning
   - Container image tagging
   - SBOM generation
   - GitHub releases

4. **Security Scanning** (`.github/workflows/security.yml`)
   - Daily vulnerability scans
   - SAST with Gosec
   - Container image scanning
   - License compliance

### Container Registry

Images are published to GitHub Container Registry:
- `ghcr.io/gaetan-jaminon/gotuto/billing:latest`
- `ghcr.io/gaetan-jaminon/gotuto/billing-dbmigrations:latest`

## 🔄 Development Workflow

### Feature Development

1. **Create feature branch**:
   ```bash
   git checkout -b feature/new-feature
   ```

2. **Make changes and test**:
   ```bash
   # Run tests frequently
   go test ./...
   
   # Run linting (if available)
   golangci-lint run
   ```

3. **Create pull request**:
   - CI pipeline runs automatically
   - Security scans execute
   - Code coverage reports generated

4. **Deploy to staging**:
   - Merge to `develop` branch
   - Automatic staging deployment
   - E2E tests execute

5. **Deploy to production**:
   - Merge to `main` branch
   - Manual approval required
   - Blue-green deployment

### Adding New Features

1. **Database changes**:
   ```bash
   cd billing-dbmigrations
   go run cmd/migrate/*.go create add_new_table
   # Edit the generated .up.sql and .down.sql files
   go run cmd/migrate/*.go up
   ```

2. **API endpoints**:
   - Add model in `internal/models/`
   - Add handler in `internal/handlers/`
   - Add routes in `cmd/main.go`
   - Write tests for each component

3. **Configuration**:
   - Add config struct in `internal/config/`
   - Update config files in `config/`
   - Document environment variables

### Code Quality Standards

- **Go modules**: Use `go mod tidy` frequently
- **Formatting**: Use `go fmt` and `goimports`
- **Linting**: Address `golangci-lint` issues
- **Testing**: Maintain >80% code coverage
- **Documentation**: Update README and godoc comments
- **Security**: Run `gosec` and address findings

## 🏭 Production Deployment

### OpenShift

The services are designed for OpenShift deployment with:

- **Non-root containers** (UID 1001)
- **Security contexts** for OpenShift SCC
- **Health checks** for liveness/readiness
- **Resource limits** and requests
- **ConfigMaps** and Secrets for configuration

Deployment manifests are in `billing/deployments/openshift/`:

```bash
# Deploy to staging
oc apply -f billing/deployments/openshift/staging/

# Deploy to production
oc apply -f billing/deployments/openshift/production/
```

### Environment Management

- **Staging**: `billing-staging` namespace, debug enabled
- **Production**: `billing-production` namespace, optimized settings

## 📚 Learning Resources

This project demonstrates:

- **Go project structure** and organization
- **Database patterns** with PostgreSQL
- **Testing strategies** (unit, integration, E2E)
- **Configuration management** with Viper
- **Container best practices** with multi-stage builds
- **CI/CD pipelines** with GitHub Actions
- **OpenShift deployment** patterns
- **Security practices** and scanning

### Key Go Patterns

- **Error handling**: Explicit error returns
- **Interface usage**: Dependency injection
- **Table-driven tests**: Go testing idiom
- **Context usage**: Request cancellation
- **Struct embedding**: Code reuse
- **Package organization**: Clean architecture

## 🤝 Contributing

1. **Fork the repository**
2. **Create feature branch**: `git checkout -b feature/amazing-feature`
3. **Write tests** for new functionality
4. **Ensure all tests pass**: `go test ./...`
5. **Commit changes**: `git commit -m 'Add amazing feature'`
6. **Push to branch**: `git push origin feature/amazing-feature`
7. **Open pull request**

### Code Review Checklist

- [ ] Tests written and passing
- [ ] Code formatted with `go fmt`
- [ ] No security vulnerabilities
- [ ] Documentation updated
- [ ] No breaking changes (or properly documented)
- [ ] Performance implications considered

## 📝 Notes for .NET Developers

Coming from .NET? Here are key differences:

- **No exceptions**: Use explicit error returns
- **No inheritance**: Use composition and interfaces
- **No null**: Use pointers sparingly, zero values instead
- **No generics** (in older Go versions): Use interfaces
- **Package-level organization**: Not class-based
- **Build process**: `go build`, no MSBuild/NuGet complexity

This project serves as a practical learning tool for transitioning from .NET to Go while maintaining enterprise development practices.

## 🔗 Links

- **Go Documentation**: https://golang.org/doc/
- **OpenShift Documentation**: https://docs.openshift.com/
- **GitHub Actions**: https://docs.github.com/en/actions
- **PostgreSQL**: https://www.postgresql.org/docs/

---

**Happy coding! 🚀**

This project is part of the GoTuto learning series, designed to help developers transition to Go with real-world, enterprise-grade examples.