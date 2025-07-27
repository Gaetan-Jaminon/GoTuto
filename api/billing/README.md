# Billing Service - Go CRUD API

A complete CRUD API built with Go, demonstrating real-world patterns for client and invoice management.

## Tech Stack

- **Language**: Go 1.22
- **Framework**: Gin (HTTP router)
- **ORM**: GORM
- **Database**: PostgreSQL
- **Migrations**: golang-migrate

## Project Structure

```
api/demo01/
├── cmd/
│   ├── main.go           # API server entry point
│   └── migrate/
│       └── main.go       # Migration tool
├── internal/
│   ├── models/           # Data models
│   ├── handlers/         # HTTP handlers (controllers)
│   └── database/         # DB connection & config
├── migrations/           # SQL migration files
├── .env.example         # Environment variables template
├── go.mod               # Go module definition
└── README.md
```

## Features

### Client Management
- ✅ Create, Read, Update, Delete clients
- ✅ Search clients by name/email
- ✅ Pagination support
- ✅ Prevent deletion of clients with invoices

### Invoice Management
- ✅ Create, Read, Update, Delete invoices
- ✅ Auto-generate invoice numbers
- ✅ Filter by client, status
- ✅ Status workflow (draft → sent → paid → overdue)
- ✅ Prevent deletion of paid invoices

### Technical Features
- ✅ Database migrations
- ✅ Soft deletes
- ✅ Request validation
- ✅ Error handling
- ✅ CORS support
- ✅ Health check endpoint

## Configuration Management

The application uses **Viper** for configuration with the following hierarchy:
1. Default values in `config/config.yaml`
2. Environment-specific overrides (`config/config.{env}.yaml`)
3. Environment variables (prefix: `DEMO01_`)
4. Command-line flags (if implemented)

### Configuration Files
- `config/config.yaml` - Base configuration
- `config/config.dev.yaml` - Development overrides
- `config/config.prod.yaml` - Production overrides
- `config/config.test.yaml` - Test environment

### Environment Variables
Environment variables override config file values. Use the prefix `DEMO01_` and replace dots with underscores:
- `DEMO01_DATABASE_HOST` → `database.host`
- `DEMO01_SERVER_PORT` → `server.port`

## Getting Started

### Prerequisites
- Go 1.22+
- PostgreSQL 12+

### 1. Setup Database
```bash
# Using Docker Compose (recommended)
docker-compose up -d postgres

# Or manually
createdb demo01
```

### 2. Configuration
```bash
# Copy example env file
cp .env.example .env

# Edit .env or config files as needed
# Default config uses localhost PostgreSQL
```

### 3. Install Dependencies
```bash
cd api/demo01
go mod tidy
```

### 4. Run Migrations
```bash
# Apply migrations
go run cmd/migrate/main.go up

# Check migration status
go run cmd/migrate/main.go version
```

### 5. Start the Server
```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Health Check
```
GET /health
```

### Clients
```
GET    /api/v1/clients              # List clients (with pagination/search)
GET    /api/v1/clients/:id          # Get client by ID
POST   /api/v1/clients              # Create new client
PUT    /api/v1/clients/:id          # Update client
DELETE /api/v1/clients/:id          # Delete client
GET    /api/v1/clients/:id/invoices # Get client's invoices
```

### Invoices
```
GET    /api/v1/invoices             # List invoices (with filters)
GET    /api/v1/invoices/:id         # Get invoice by ID
POST   /api/v1/invoices             # Create new invoice
PUT    /api/v1/invoices/:id         # Update invoice
DELETE /api/v1/invoices/:id         # Delete invoice
```

## Example Usage

### Create a Client
```bash
curl -X POST http://localhost:8080/api/v1/clients \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+1234567890",
    "address": "123 Main St, City, Country"
  }'
```

### Create an Invoice
```bash
curl -X POST http://localhost:8080/api/v1/invoices \
  -H "Content-Type: application/json" \
  -d '{
    "client_id": 1,
    "amount": 1500.50,
    "status": "draft",
    "issue_date": "2024-01-15T00:00:00Z",
    "due_date": "2024-02-15T00:00:00Z",
    "description": "Website development services"
  }'
```

### Search Clients
```bash
curl "http://localhost:8080/api/v1/clients?search=john&page=1&limit=10"
```

### Filter Invoices
```bash
curl "http://localhost:8080/api/v1/invoices?client_id=1&status=paid"
```

## Migration Commands

```bash
# Apply all pending migrations
go run cmd/migrate/main.go up

# Rollback all migrations
go run cmd/migrate/main.go down

# Check current version
go run cmd/migrate/main.go version

# Force specific version (if migrations are dirty)
go run cmd/migrate/main.go force 1
```

## Database Schema

### Clients Table
- `id` (Primary Key)
- `name` (Required)
- `email` (Required, Unique)
- `phone` (Optional)
- `address` (Optional)
- `created_at`, `updated_at`, `deleted_at`

### Invoices Table
- `id` (Primary Key)
- `number` (Unique, Auto-generated)
- `client_id` (Foreign Key)
- `amount` (Required)
- `status` (draft/sent/paid/overdue/cancelled)
- `issue_date`, `due_date` (Required)
- `description` (Optional)
- `created_at`, `updated_at`, `deleted_at`

## Deployment Examples

### Docker Compose
```bash
# Start all services (PostgreSQL, API, pgAdmin)
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop services
docker-compose down
```

### Kubernetes
```bash
# Apply configuration
kubectl apply -f deployments/kubernetes/

# Check deployment
kubectl get pods -l app=demo01-api
kubectl logs -l app=demo01-api
```

## Key Go Concepts Demonstrated

1. **Project Structure** - Standard Go layout
2. **Modules & Packages** - Clean separation of concerns
3. **HTTP Handlers** - RESTful API design
4. **Database Integration** - GORM for ORM operations
5. **Error Handling** - Proper error responses
6. **Validation** - Request validation with Gin
7. **Configuration Management** - Viper with multi-environment support
8. **Database Migrations** - Schema versioning
9. **Relationships** - Foreign keys and preloading
10. **Middleware** - CORS, logging, recovery
11. **Docker Multi-stage Builds** - Optimized container images
12. **Kubernetes Deployment** - Production-ready manifests

## Configuration Best Practices

### Development
- Use `.env` file for local overrides
- Config files in `config/` directory
- Hot-reload friendly

### Production
- Environment variables for secrets
- ConfigMaps for non-sensitive config
- External secret management (Vault, AWS Secrets Manager)
- Immutable config (rebuild on changes)

This project demonstrates real-world Go patterns for production applications!