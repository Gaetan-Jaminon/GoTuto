# GoTuto - Learn Go with Real-World Domain-Driven Design

A practical Go learning project designed for .NET developers transitioning to Go. This project demonstrates **Domain-Driven Design** with **proper separation of concerns** and enterprise development practices while following the principle: **start simple, add complexity gradually**.

## 🎯 Learning Philosophy

This project evolved through real development challenges, demonstrating:
- **Start with basics**: Simple unit tests, core functionality, local development
- **Add complexity incrementally**: Tests first, then CI/CD, then advanced patterns
- **Learn from failures**: Document troubleshooting and decisions
- **Focus on working solutions**: Practical over perfect
- **Embrace simplicity**: Remove complexity that doesn't add value
- **True DDD separation**: Each domain owns its complete context

## 📁 Current Project Structure (Domain-First Architecture)

```
GoTuto/
├── cmd/                                    # Application entry points
│   ├── billing-api/main.go                # Billing API service
│   ├── billing-migrator/main.go           # Billing migration tool
│   └── catalog-migrator/main.go           # Catalog migration tool
├── config/                                 # Domain-first configuration
│   ├── base/                              # Shared infrastructure config
│   │   ├── base.yaml                      # Common server, database, logging
│   │   ├── dev.yaml                       # Development overrides
│   │   ├── qua.yaml                       # QA environment overrides
│   │   └── prd.yaml                       # Production overrides
│   ├── billing/                           # Billing domain configuration
│   │   ├── billing.yaml                   # Billing-specific settings
│   │   ├── dev.yaml                       # Billing dev overrides
│   │   ├── qua.yaml                       # Billing QA overrides
│   │   └── prd.yaml                       # Billing prod overrides
│   └── catalog/                           # Catalog domain configuration
│       ├── catalog.yaml                   # Catalog-specific settings
│       ├── dev.yaml                       # Catalog dev overrides
│       ├── qua.yaml                       # Catalog QA overrides
│       └── prd.yaml                       # Catalog prod overrides
├── internal/                               # Domain-Driven Design organization
│   ├── shared/
│   │   └── infrastructure/                # Truly shared utilities
│   │       ├── config.go                  # Generic config loader
│   │       ├── server.go                  # ServerConfig struct
│   │       ├── database.go                # DatabaseConfig struct + schema support
│   │       ├── logging.go                 # LoggingConfig struct
│   │       └── cors.go                    # CORSConfig struct
│   ├── billing/                           # BILLING DOMAIN (complete isolation)
│   │   ├── config/config.go               # Billing config (BILLING_ env prefix)
│   │   ├── migrations/                    # Billing schema migrations
│   │   │   ├── 001_create_billing_schema.up.sql
│   │   │   ├── 001_create_billing_schema.down.sql
│   │   │   ├── 002_billing_tables.up.sql
│   │   │   └── 002_billing_tables.down.sql
│   │   ├── database/connection.go         # Billing database (billing schema)
│   │   ├── models/                        # Billing domain entities
│   │   │   ├── client.go                  # Client entity + DTOs
│   │   │   └── invoice.go                 # Invoice entity + DTOs
│   │   ├── api/                           # Billing application services
│   │   │   ├── client.go                  # Client HTTP handlers
│   │   │   └── invoice.go                 # Invoice HTTP handlers
│   │   ├── services/                      # Billing domain services
│   │   └── repositories/                  # Billing data interfaces
│   └── catalog/                           # CATALOG DOMAIN (complete isolation)
│       ├── config/config.go               # Catalog config (CATALOG_ env prefix)
│       ├── migrations/                    # Catalog schema migrations
│       │   ├── 001_create_catalog_schema.up.sql
│       │   ├── 001_create_catalog_schema.down.sql
│       │   ├── 002_catalog_tables.up.sql
│       │   └── 002_catalog_tables.down.sql
│       ├── database/connection.go         # Catalog database (catalog schema)
│       ├── models/                        # Catalog domain entities
│       ├── api/                           # Catalog application services
│       ├── services/                      # Catalog domain services
│       └── repositories/                  # Catalog data interfaces
├── scripts/                                # Developer convenience tools
│   ├── build.sh                          # Build all binaries
│   ├── test.sh                           # Run tests with coverage
│   ├── test-unit.sh                      # Quick unit tests
│   ├── lint.sh                           # Run linting tools
│   ├── docker-build.sh                   # Build Docker images
│   ├── clean.sh                          # Clean build artifacts
│   └── dev-setup.sh                      # Set up development environment
├── billing-api.Dockerfile                 # API service container
├── billing-migrator.Dockerfile            # Billing migration container
├── Makefile                                # Familiar interface (make build, make test)
├── .github/workflows/                      # Simplified CI (Claude integration)
│   └── claude-code-review.yml            # Automated PR reviews
├── test/                                   # Integration & E2E tests (future)
└── notes/                                  # Learning documentation
```

## 🏗️ Architecture: True Domain-Driven Design

### Domain-First Principles

This project demonstrates **enterprise-grade DDD** with complete domain separation:

#### 🎯 **Developer Cognitive Load Reduction**
When working on billing features, developers only need to focus on:
- **Code**: `internal/billing/`
- **Config**: `config/billing/`
- **Migrations**: `internal/billing/migrations/`

Everything is co-located, reducing context switching and mental overhead.

#### 🔒 **Database Schema Isolation**
```sql
-- Single database, multiple schemas
Database: gotuto_dev / gotuto_qua / gotuto_prd

-- Domain-specific schemas
├── billing schema
│   ├── clients table
│   └── invoices table
└── catalog schema
    ├── products table
    ├── categories table
    └── product_categories table
```

#### 🛡️ **RBAC at Database Level**
```sql
-- Domain-specific users prevent cross-domain access
billing_app     → USAGE + DML on billing schema only
catalog_app     → USAGE + DML on catalog schema only

billing_migrator → CREATE + DDL on billing schema only
catalog_migrator → CREATE + DDL on catalog schema only
```

### Configuration Architecture

#### Domain-First Configuration Loading
```
config/billing/billing.yaml  # Domain defaults
         ├── dev.yaml         # Environment overrides
         ├── qua.yaml         # QA overrides
         └── prd.yaml         # Production overrides

config/catalog/catalog.yaml  # Domain defaults
         ├── dev.yaml         # Environment overrides
         ├── qua.yaml         # QA overrides
         └── prd.yaml         # Production overrides
```

#### Environment Variable Strategy
```bash
# Billing domain (BILLING_ prefix)
BILLING_DATABASE_HOST=localhost
BILLING_DATABASE_PASSWORD=secret
BILLING_SERVER_PORT=8080

# Catalog domain (CATALOG_ prefix)  
CATALOG_DATABASE_HOST=localhost
CATALOG_DATABASE_PASSWORD=secret
CATALOG_SERVER_PORT=8081
```

### Migration Architecture

Each domain manages its own schema:

```bash
# Billing migrations (billing schema only)
./bin/billing-migrator up

# Catalog migrations (catalog schema only) 
./bin/catalog-migrator up
```

**Key Benefits:**
- ✅ No cross-schema migrations possible
- ✅ Database-level enforcement of boundaries
- ✅ Independent domain evolution
- ✅ Easy microservice extraction

## 🚀 Quick Start

### Prerequisites
- **Go 1.22+**
- **Git**
- **Docker** (optional, for containers)
- **PostgreSQL 15+** (optional, for full API functionality)

### 1. Clone and Set Up

```bash
git clone https://github.com/Gaetan-Jaminon/GoTuto.git
cd GoTuto

# Set up development environment
./scripts/dev-setup.sh
```

### 2. Quick Development Workflow

```bash
# Run unit tests (fast, no dependencies)
./scripts/test-unit.sh
# OR
make test

# Build all applications
./scripts/build.sh
# OR
make build

# Run linting
./scripts/lint.sh
# OR
make lint

# Build Docker images
./scripts/docker-build.sh
# OR
make docker
```

### 3. Database Setup (Schema-Based Separation)

```bash
# Start PostgreSQL
docker run -d \
  --name postgres-dev \
  -e POSTGRES_DB=gotuto_dev \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:15

# Set up database users and schemas
psql -h localhost -U postgres -d gotuto_dev -c "
  -- Create schemas
  CREATE SCHEMA IF NOT EXISTS billing;
  CREATE SCHEMA IF NOT EXISTS catalog;
  
  -- Create domain users
  CREATE USER billing_app WITH PASSWORD 'billing_pass';
  CREATE USER catalog_app WITH PASSWORD 'catalog_pass';
  CREATE USER billing_migrator WITH PASSWORD 'billing_migrate_pass';
  CREATE USER catalog_migrator WITH PASSWORD 'catalog_migrate_pass';
  
  -- Grant schema permissions (RBAC)
  GRANT USAGE ON SCHEMA billing TO billing_app;
  GRANT ALL ON ALL TABLES IN SCHEMA billing TO billing_app;
  GRANT CREATE ON SCHEMA billing TO billing_migrator;
  
  GRANT USAGE ON SCHEMA catalog TO catalog_app;
  GRANT ALL ON ALL TABLES IN SCHEMA catalog TO catalog_app;
  GRANT CREATE ON SCHEMA catalog TO catalog_migrator;
"
```

### 4. Run Migrations

```bash
# Migrate billing schema
APP_ENV=dev ./bin/billing-migrator up

# Migrate catalog schema  
APP_ENV=dev ./bin/catalog-migrator up
```

### 5. Start the API

```bash
# Build and run billing API
./scripts/build.sh
APP_ENV=dev ./bin/billing-api
```

### 6. Test Domain Separation

```bash
# Test billing API
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/clients

# Each domain has its own:
# - Database schema (billing vs catalog)
# - Configuration files (config/billing/ vs config/catalog/)
# - Migration tools (billing-migrator vs catalog-migrator)
# - Environment prefixes (BILLING_ vs CATALOG_)
```

## 🧪 Testing Strategy (Domain-Focused)

### Current Approach: Domain-Specific Unit Tests

Each domain has its own test suite with no cross-domain dependencies:

```bash
# Test billing domain only
go test ./internal/billing/...

# Test catalog domain only  
go test ./internal/catalog/...

# Test all domains
./scripts/test-unit.sh
```

**Domain-Specific Tests:**
- **Billing Models** (`internal/billing/models/*_test.go`)
- **Catalog Models** (`internal/catalog/models/*_test.go`)
- **Configuration Loading** (per domain)
- **Business Rules** (domain-specific validation)

## 🏗️ API Overview

### Billing Service (Port 8080)

RESTful API demonstrating Go web development with billing domain:

**Health Check:**
```bash
curl http://localhost:8080/health
# Returns: {"status":"healthy","service":"billing-api","domain":"billing"}
```

**Billing Endpoints:**
```
GET    /health                      # Health check with domain info
GET    /api/v1/clients             # List billing clients  
POST   /api/v1/clients             # Create billing client
GET    /api/v1/clients/{id}        # Get billing client
PUT    /api/v1/clients/{id}        # Update billing client
DELETE /api/v1/clients/{id}        # Delete billing client
GET    /api/v1/invoices            # List invoices
POST   /api/v1/invoices            # Create invoice
GET    /api/v1/invoices/{id}       # Get invoice
PUT    /api/v1/invoices/{id}       # Update invoice
DELETE /api/v1/invoices/{id}       # Delete invoice
```

### Future: Catalog Service (Port 8081)

When implemented, the catalog service will have its own API:

**Catalog Endpoints:**
```
GET    /health                      # Health check (catalog domain)
GET    /api/v1/products            # List catalog products
POST   /api/v1/products            # Create product
GET    /api/v1/categories          # List categories
POST   /api/v1/categories          # Create category
```

## ⚙️ Configuration Management (Domain-First)

### Hierarchical Configuration Loading

Each domain loads configuration in this order:
1. **Base defaults** (`config/base/base.yaml`)
2. **Base environment** (`config/base/{env}.yaml`)
3. **Domain defaults** (`config/{domain}/{domain}.yaml`)
4. **Domain environment** (`config/{domain}/{env}.yaml`)
5. **Environment variables** (`{DOMAIN}_*`)

### Billing Domain Configuration

```yaml
# config/billing/billing.yaml
database:
  name: "gotuto"
  schema: "billing"              # Schema isolation
  username: "billing_app"        # Domain-specific user

pagination:
  default_limit: 10
  max_limit: 100

invoice:
  number_prefix: "INV"
  default_currency: "USD"
```

### Catalog Domain Configuration

```yaml
# config/catalog/catalog.yaml  
database:
  name: "gotuto"
  schema: "catalog"              # Schema isolation
  username: "catalog_app"        # Domain-specific user

pagination:
  default_limit: 20
  max_limit: 50

product:
  sku_prefix: "SKU"
  default_currency: "USD"
```

### Environment Variables

**Billing Domain:**
```bash
BILLING_DATABASE_HOST=localhost
BILLING_DATABASE_PASSWORD=secret
BILLING_SERVER_PORT=8080
BILLING_PAGINATION_DEFAULT_LIMIT=10
```

**Catalog Domain:**
```bash
CATALOG_DATABASE_HOST=localhost
CATALOG_DATABASE_PASSWORD=secret
CATALOG_SERVER_PORT=8081
CATALOG_PAGINATION_DEFAULT_LIMIT=20
```

## 🐳 Container Support (Domain-Aware)

### Domain-Specific Images

```bash
# Build domain-specific images
./scripts/docker-build.sh

# Or manually:
docker build -f billing-api.Dockerfile -t billing-api:local .
docker build -f billing-migrator.Dockerfile -t billing-migrator:local .

# Run billing API
docker run -p 8080:8080 \
  -e BILLING_DATABASE_HOST=host.docker.internal \
  -e BILLING_DATABASE_PASSWORD=secret \
  billing-api:local

# Run billing migrations
docker run \
  -e BILLING_DATABASE_HOST=host.docker.internal \
  -e BILLING_DATABASE_PASSWORD=secret \
  billing-migrator:local up
```

## 🤖 AI-Powered Development

### Claude Integration

**Interactive Claude:**
- Mention `@claude` in issues/PRs for domain-specific help
- Get explanations about DDD patterns
- Learn Go idioms and best practices

**Automated Code Review:**
- Reviews focus on domain separation
- Validates DDD principles
- Suggests improvements for configuration architecture

## 🎓 Learning Resources

### Domain-Driven Design Concepts

**Key DDD Patterns Demonstrated:**

| DDD Concept | Implementation | Location |
|------------|----------------|----------|
| **Bounded Context** | Complete domain isolation | `internal/billing/` vs `internal/catalog/` |
| **Ubiquitous Language** | Domain-specific models | `internal/billing/models/` |
| **Aggregate Root** | Client, Invoice entities | `models/client.go`, `models/invoice.go` |
| **Repository Pattern** | Data access interfaces | `repositories/` (planned) |
| **Domain Services** | Business logic | `services/` (planned) |
| **Application Services** | HTTP handlers | `api/client.go`, `api/invoice.go` |
| **Infrastructure** | Database, config | `database/`, `config/` |

### Configuration Architecture Benefits

**For .NET Developers:**

| .NET Pattern | Go Equivalent | Notes |
|-------------|---------------|--------|
| **appsettings.json** | `config/{domain}/{domain}.yaml` | Domain-specific configs |
| **IConfiguration** | `config.Load()` | Type-safe config loading |
| **Environment-specific configs** | `config/{domain}/{env}.yaml` | Hierarchical overrides |
| **IOptions<T>** | Domain config structs | Strongly-typed configuration |

### Database Schema Strategy

**Coming from .NET/SQL Server:**

| .NET Pattern | Go/PostgreSQL Pattern | Benefits |
|-------------|---------------------|----------|
| **Separate Databases** | **Single DB, Multiple Schemas** | Simpler ops, logical separation |
| **ConnectionStrings** | **Schema-aware DSN** | search_path enforces boundaries |
| **EF DbContext per domain** | **GORM DB per schema** | Domain isolation maintained |
| **SQL Users per app** | **PostgreSQL users per domain** | Database-level security |

## 📚 Learning Path

### 1. Master Domain Separation
```bash
# Explore billing domain
ls internal/billing/
cat config/billing/billing.yaml

# See how config loading works
go run -c "
  cfg, _ := billing.Load()
  fmt.Printf('Schema: %s', cfg.Database.Schema)
"
```

### 2. Understand Schema Isolation
```bash
# Connect to billing schema
BILLING_DATABASE_SCHEMA=billing ./bin/billing-migrator version

# Connect to catalog schema  
CATALOG_DATABASE_SCHEMA=catalog ./bin/catalog-migrator version
```

### 3. Practice Adding Domains
```bash
# Add a new domain (e.g., inventory)
mkdir -p internal/inventory/{config,migrations,models,api}
mkdir -p config/inventory

# Follow the billing domain pattern
cp -r internal/billing/* internal/inventory/
# Adapt for inventory domain...
```

## 🔍 Troubleshooting

### Domain Configuration Issues

**1. Wrong schema accessed**
```bash
# Check DSN includes search_path
APP_ENV=dev go run -c "
  cfg, _ := billing.Load()
  fmt.Println(cfg.Database.GetDSN())
"
# Should include: search_path=billing
```

**2. Config not loading**
```bash
# Verify config file structure
ls config/billing/
cat config/billing/billing.yaml

# Test environment loading
APP_ENV=dev go run -c "
  cfg, _ := billing.Load()
  fmt.Printf('Loaded from: %s\n', cfg.Database.Name)
"
```

**3. Cross-domain access**
```bash
# This should fail (good!)
psql -U billing_app -d gotuto_dev -c "SELECT * FROM catalog.products;"
# ERROR: permission denied for schema catalog
```

## 🎯 Next Steps

### Adding the Catalog Domain (Complete DDD Exercise)

1. **Implement Catalog Models**:
   ```bash
   # Create product and category entities
   touch internal/catalog/models/{product,category,brand}.go
   ```

2. **Add Catalog API**:
   ```bash
   # Implement CRUD endpoints
   touch internal/catalog/api/{product,category}.go
   ```

3. **Create Catalog API Service**:
   ```bash
   # New binary for catalog domain
   mkdir cmd/catalog-api
   touch cmd/catalog-api/main.go
   ```

4. **Test Domain Isolation**:
   ```bash
   # Verify no cross-domain imports
   go mod graph | grep "billing.*catalog\|catalog.*billing"
   # Should return nothing
   ```

### Architecture Evolution Path

```
Current: Domain-First Monolith    →    Future: Domain-Based Microservices
├── config/billing/              →    billing-service/config/
├── internal/billing/            →    billing-service/internal/
├── config/catalog/              →    catalog-service/config/
└── internal/catalog/            →    catalog-service/internal/
```

### Advanced DDD Patterns to Explore

1. **Domain Events**: Inter-domain communication
2. **CQRS**: Separate read/write models
3. **Event Sourcing**: Event-based state changes
4. **Saga Pattern**: Cross-domain transactions

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/catalog-products`
3. **Follow domain boundaries**: Keep changes within single domain
4. Write domain-specific tests: `go test ./internal/catalog/...`
5. Ensure configuration loading works: Test with different environments
6. Commit with domain prefix: `git commit -m "catalog: add product model"`
7. Create pull request and learn from Claude's DDD feedback

## 📖 Philosophy: Domain-First Development

This project demonstrates enterprise-grade Domain-Driven Design in Go:

### Core Principles Applied

✅ **Bounded Contexts**: Each domain has complete isolation
✅ **Ubiquitous Language**: Domain-specific models and terminology  
✅ **Schema Separation**: Database-level domain boundaries
✅ **Configuration Isolation**: Domain-first config management
✅ **Developer Experience**: Cognitive load reduction through co-location

### Key Lessons for .NET Developers

1. **Go's simplicity enables DDD**: Less ceremony, clearer domain focus
2. **Schema separation > microservices**: Start with logical boundaries
3. **Configuration co-location**: Domain owns all its concerns
4. **Database-level RBAC**: Security through isolation
5. **Migration per domain**: Independent evolution paths

### When to Extract to Microservices

🟢 **Stay monolithic when:**
- Domains fit in single database
- Team is small (< 10 developers)
- Deployment complexity isn't justified

🔴 **Extract to microservices when:**
- Domain teams are independent
- Different scaling requirements
- Technology diversity needed

**Remember**: Good domain boundaries in a monolith become good service boundaries in microservices.

---

## 🚀 Happy Domain-Driven Go Learning!

This project demonstrates that learning DDD doesn't require microservices complexity from day one. 

**Start with domains → Add boundaries → Scale when needed**

The combination of:
- **True Domain Separation** (clear boundaries)
- **Schema-Based Isolation** (database-level security)
- **Domain-First Configuration** (developer experience)
- **Go Simplicity** (focus on domain, not framework)

...creates an excellent environment for mastering both Go and Domain-Driven Design.

**Questions about DDD patterns?** Create an issue or mention `@claude` in a pull request - AI assistance with domain modeling is always available!