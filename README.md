# GoTuto - Learn Go with Real-World Examples

A practical Go learning project designed for .NET developers transitioning to Go. This project demonstrates Domain-Driven Design and enterprise development practices while following the principle: **start simple, add complexity gradually**.

## 🎯 Learning Philosophy

This project evolved through real development challenges, demonstrating:
- **Start with basics**: Simple unit tests, core functionality, local development
- **Add complexity incrementally**: Tests first, then CI/CD, then advanced patterns
- **Learn from failures**: Document troubleshooting and decisions
- **Focus on working solutions**: Practical over perfect
- **Embrace simplicity**: Remove complexity that doesn't add value

## 📁 Current Project Structure (Single Module + DDD)

```
GoTuto/
├── cmd/                          # Application entry points
│   ├── billing-api/main.go      # Billing API service
│   └── billing-migrator/main.go # Database migration service
├── internal/                     # Domain-Driven Design organization
│   ├── billing/                 # DOMAIN: Business logic
│   │   ├── api/                 # Application services (handlers)
│   │   ├── models/              # Domain entities
│   │   ├── services/            # Domain services
│   │   └── repositories/        # Domain interfaces
│   ├── billing-migration/       # INFRASTRUCTURE: Data management
│   │   ├── database/            # Database connections
│   │   ├── runners/             # Migration execution
│   │   └── scripts/             # Migration scripts
│   └── shared/                  # Shared utilities (config, etc.)
├── scripts/                      # Developer convenience tools
│   ├── build.sh                # Build both binaries
│   ├── test.sh                 # Run tests with coverage
│   ├── test-unit.sh            # Quick unit tests
│   ├── lint.sh                 # Run linting tools
│   ├── docker-build.sh         # Build Docker images
│   ├── clean.sh                # Clean build artifacts
│   └── dev-setup.sh            # Set up development environment
├── billing-api.Dockerfile       # API service container (root = Go convention)
├── billing-migrator.Dockerfile  # Migration service container
├── Makefile                      # Familiar interface (make build, make test)
├── .github/workflows/           # Simplified CI (Claude integration only)
│   ├── claude.yml              # Interactive Claude (@claude mentions)
│   └── claude-code-review.yml  # Automated PR reviews
├── test/                        # Integration & E2E tests (future)
├── config/                      # Configuration files
└── notes/                       # Learning documentation
```

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

# Set up development environment (installs tools, makes scripts executable)
./scripts/dev-setup.sh
```

### 2. Quick Development Workflow

```bash
# Run unit tests (fast, no dependencies)
./scripts/test-unit.sh
# OR
make test

# Build both applications
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

# Clean everything
./scripts/clean.sh
# OR
make clean
```

### 3. Run the API (with Database)

```bash
# Start PostgreSQL (using Docker)
docker run -d \
  --name postgres-dev \
  -e POSTGRES_DB=billing \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  postgres:15

# Set environment variables
export BILLING_DATABASE_HOST=localhost
export BILLING_DATABASE_PASSWORD=password

# Build and run the API
./scripts/build.sh
./bin/billing-api
```

### 4. Test the API

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/clients
```

## 🧪 Testing Strategy (Simplified & Effective)

### Current Approach: Unit Tests Only

We learned that complex testing setups can become blockers. Current philosophy: **start with tests that always work**.

```bash
# Fast feedback loop (runs everywhere)
./scripts/test-unit.sh

# With coverage reports
./scripts/test.sh
```

**Current Tests:**
- **Configuration validation** (`internal/shared/config_test.go`)
- **Domain model logic** (`internal/billing/models/*_test.go`)
- **Business rules** (invoice status transitions, validation)

### Why We Simplified

**❌ Removed (barriers to learning):**
- Integration tests requiring PostgreSQL setup
- E2E tests requiring full application stack  
- Handler tests requiring database mocking
- Complex CI pipelines blocking development

**✅ Kept (always works):**
- Unit tests with no external dependencies
- Domain logic validation
- Business rule testing
- Local development tools

**Lesson**: Start with tests that run everywhere instantly. Add complexity when the foundation is solid.

## 🏗️ Architecture: Domain-Driven Design in Go

### Single Module with Domain Separation

This project demonstrates **DDD in Go** while maintaining simplicity:

**Domain Layer** (`internal/billing/`):
- **Models**: Core business entities (Client, Invoice)
- **Services**: Business logic and rules
- **Repositories**: Domain interfaces (ports)
- **API**: Application services (use cases)

**Infrastructure Layer** (`internal/billing-migration/`):
- **Database**: Connection management
- **Runners**: Migration execution
- **Scripts**: Schema definitions

**Shared Kernel** (`internal/shared/`):
- Configuration management
- Common utilities

### Why Single Module?

**Current Benefits:**
- Simpler dependency management
- Easier local development
- Single test suite
- Go-idiomatic structure

**Future-Proof Design:**
- Clear domain boundaries
- Easy extraction to separate repositories:
  ```
  billing-api/        # Domain + API
  billing-migration/  # Infrastructure
  shared-library/     # Common utilities
  ```

## 🏗️ API Overview

### Billing Service

A RESTful API demonstrating Go web development patterns:

**Core Features:**
- CRUD operations for clients and invoices
- Configuration management with Viper
- Database integration with GORM
- HTTP routing with Gin
- Structured logging
- Health checks

**API Endpoints:**
```
GET    /health                 # Health check
GET    /api/v1/clients        # List clients  
POST   /api/v1/clients        # Create client
GET    /api/v1/clients/{id}   # Get client
PUT    /api/v1/clients/{id}   # Update client
DELETE /api/v1/clients/{id}   # Delete client
GET    /api/v1/invoices       # List invoices
POST   /api/v1/invoices       # Create invoice
GET    /api/v1/invoices/{id}  # Get invoice
PUT    /api/v1/invoices/{id}  # Update invoice
DELETE /api/v1/invoices/{id}  # Delete invoice
```

## 🤖 AI-Powered Development (Simplified)

### Claude Integration Only

We streamlined to focus on what adds real value: **AI-assisted learning**.

**Current Workflows:**

**1. Interactive Claude** (`claude.yml`):
- Mention `@claude` in issues/PRs for help
- Get explanations, suggestions, code reviews
- Available 24/7 for learning support

**2. Automated Code Review** (`claude-code-review.yml`):
- Automatic review on every pull request
- Focuses on Go best practices
- Provides learning feedback
- No CI dependency - pure code review

**Benefits:**
- Learn Go patterns through AI feedback
- Get instant help with `@claude` mentions
- Consistent code quality
- Zero infrastructure maintenance

**Removed Complexity:**
- ❌ Complex CI pipelines
- ❌ Dependency update automation
- ❌ Security scanning workflows
- ❌ Deployment automation

**Lesson**: Focus on AI assistance that helps learning, remove automation that creates maintenance overhead.

## 🛠️ Developer Experience

### Scripts Directory (Local Development)

The `scripts/` directory provides consistent development commands:

**Build & Test:**
```bash
./scripts/build.sh         # Build both binaries to bin/
./scripts/test.sh          # Run tests with coverage
./scripts/test-unit.sh     # Quick unit tests only
```

**Code Quality:**
```bash
./scripts/lint.sh          # Run go fmt, go vet, golangci-lint
```

**Docker:**
```bash
./scripts/docker-build.sh  # Build both Docker images locally
./scripts/clean.sh         # Clean artifacts and images
```

**Setup:**
```bash
./scripts/dev-setup.sh     # Install tools, test build
```

### Makefile (Familiar Interface)

For developers who prefer make:
```bash
make build    # Same as ./scripts/build.sh
make test     # Same as ./scripts/test.sh
make lint     # Same as ./scripts/lint.sh
make docker   # Same as ./scripts/docker-build.sh
make clean    # Same as ./scripts/clean.sh
```

### Why Scripts + Makefile?

**Benefits:**
- **Consistency**: Same commands work for all developers
- **Speed**: No need to push to GitHub to test changes
- **Documentation**: Scripts serve as executable documentation
- **CI Simulation**: Run the same checks locally
- **Onboarding**: New developers can start immediately

## ⚙️ Configuration Management

Hierarchical configuration using Viper:

1. **Default values** (in code)
2. **Config files** (`config/config.yaml`)
3. **Environment variables** (prefixed `BILLING_`)
4. **Command line flags**

### Key Environment Variables

```bash
# Server Configuration
BILLING_SERVER_PORT=8080

# Database Configuration  
BILLING_DATABASE_HOST=localhost
BILLING_DATABASE_PORT=5432
BILLING_DATABASE_USERNAME=postgres
BILLING_DATABASE_PASSWORD=password
BILLING_DATABASE_NAME=billing
BILLING_DATABASE_SSL_MODE=disable

# Application Configuration
BILLING_LOG_LEVEL=info
APP_ENV=development
```

## 🐳 Container Support

### Two Services, Two Images

Following Go conventions with root-level Dockerfiles:

```bash
# Build both images locally
./scripts/docker-build.sh

# Or manually:
docker build -f billing-api.Dockerfile -t billing-api:local .
docker build -f billing-migrator.Dockerfile -t billing-migrator:local .

# Run API container
docker run -p 8080:8080 \
  -e BILLING_DATABASE_HOST=host.docker.internal \
  -e BILLING_DATABASE_PASSWORD=password \
  billing-api:local
```

**Container Features:**
- Red Hat UBI base images (enterprise-ready)
- Non-root user (security)
- Optimized for OpenShift/Kubernetes
- Independent versioning per service

## 🎓 Learning Resources

### Go Concepts Demonstrated

**Core Language:**
- Package organization and Go modules
- Error handling patterns (`if err != nil`)
- Interfaces and composition (no inheritance)
- Struct methods and receivers
- Short variable declaration (`:=`)

**Domain-Driven Design:**
- Domain vs Infrastructure separation
- Entity and value object patterns
- Repository interfaces
- Application services

**Web Development:**
- HTTP handlers and routing (Gin)
- JSON marshaling/unmarshaling
- Middleware patterns
- Request validation

**Database Integration:**
- ORM patterns with GORM
- Connection management
- Configuration-based connections

**Testing:**
- Table-driven tests (Go idiom)
- Test organization and structure
- Coverage measurement

**DevOps:**
- Container deployment
- Configuration management
- AI-powered development workflow

### Notes for .NET Developers

Key differences when coming from .NET:

| .NET Concept | Go Equivalent | Notes |
|-------------|---------------|--------|
| **Solution (.sln)** | **Repository** | Multiple Go modules in one repo |
| **Project (.csproj)** | **Go Module (go.mod)** | Single module in this project |
| **Class** | **Struct + methods** | No inheritance, use composition |
| **try/catch** | **Explicit error returns** | `if err != nil { return err }` |
| **null** | **Zero values** | Prefer zero values over pointers |
| **NuGet** | **Go modules** | Simpler dependency management |
| **MSBuild** | **go build** | No complex build configuration |
| **LINQ** | **Manual iteration** | More explicit, less magic |
| **Entity Framework** | **GORM** | Similar ORM patterns |

### Domain-Driven Design Comparison

| DDD Concept | .NET Implementation | Go Implementation |
|------------|-------------------|------------------|
| **Domain Layer** | `Domain` assembly | `internal/billing/` |
| **Application Layer** | `Application` assembly | `internal/billing/api/` |
| **Infrastructure** | `Infrastructure` assembly | `internal/billing-migration/` |
| **Bounded Context** | Separate projects/solutions | Clear package boundaries |

## 📚 Documentation and Learning Notes

### Learning Documentation
- `notes/packages-vs-modules.md` - Go project organization
- `notes/short-variable-declaration.md` - `:=` operator usage
- `notes/functions-comprehensive.md` - Go functions and error handling
- `notes/claude-github-integration.md` - AI-powered development workflow
- `notes/branch-protection-setup.md` - Manual branch protection guide

### Development Workflow

1. **Create feature branch**:
   ```bash
   git checkout -b feature/new-feature
   ```

2. **Develop and test locally**:
   ```bash
   ./scripts/test-unit.sh  # Fast feedback loop
   ./scripts/build.sh      # Test build
   ./bin/billing-api       # Test manually
   ```

3. **Create pull request**:
   - Claude provides automatic code review
   - Learn from AI feedback
   - Merge after review

4. **Learn from feedback**:
   - Review Claude's suggestions
   - Ask questions with `@claude` mentions
   - Document new patterns learned

## 🔍 Troubleshooting

### Common Issues

**1. Scripts not executable**
```bash
chmod +x scripts/*.sh
```

**2. Tests fail with import errors**
```bash
go mod tidy
go clean -modcache
```

**3. Build fails with "command not found"**
```bash
# Run from project root
cd GoTuto
./scripts/build.sh
```

**4. Claude review doesn't run**
- Claude runs on every PR automatically
- Check GitHub Actions tab for status
- Try mentioning `@claude` for interactive help

**5. Can't find binaries after build**
```bash
# Binaries are in bin/ directory
ls bin/
./bin/billing-api
./bin/billing-migrator
```

### Getting Help

1. **Check scripts output**: Scripts provide helpful error messages
2. **Use `@claude` mentions**: Ask specific questions in issues/PRs
3. **Review notes**: Check `notes/` directory for specific topics
4. **Test locally first**: `./scripts/test-unit.sh` for quick validation

## 🎯 Next Steps

### Immediate Learning Goals
1. **Master Go basics**: Complete examples in `notes/`
2. **Practice DDD**: Add new domain concepts (products, orders)
3. **Understand testing**: Write unit tests for new features
4. **Learn from AI feedback**: Create PRs and review Claude's suggestions

### When Ready to Add Complexity
1. **Integration tests**: When comfortable with unit testing
2. **CI/CD pipeline**: When local development is smooth
3. **Authentication**: JWT or session-based auth
4. **More services**: Split into multiple repositories
5. **Advanced patterns**: CQRS, Event Sourcing

### Architecture Evolution Path
```
Current: Single Module           →    Future: Multiple Repos
├── internal/billing/           →    billing-api/
├── internal/billing-migration/ →    billing-migrations/
└── internal/shared/            →    shared-library/
```

### Advanced Topics to Explore Later
- Goroutines and channels (concurrency)
- Advanced database patterns (CQRS, Event Sourcing)
- Microservices communication
- Performance optimization
- Security hardening

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Write unit tests for new functionality: `./scripts/test-unit.sh`
4. Ensure all checks pass: `./scripts/lint.sh`
5. Commit with descriptive messages
6. Push and create pull request
7. Learn from Claude's automatic code review feedback

## 📖 Philosophy: Simplicity Through Experience

This project embodies a key software development principle:

> **"Make it work, make it right, make it fast"** - Kent Beck

### Our Learning Journey

**Started with enterprise complexity:**
- Complex CI/CD pipelines
- Integration and E2E tests
- Multiple automated workflows
- Over-engineered architecture

**Learned to embrace simplicity:**
- ✅ Working unit tests > Complex integration tests
- ✅ Local scripts > Complex CI pipelines
- ✅ Manual setup > Unreliable automation  
- ✅ Domain separation > Microservices complexity
- ✅ AI assistance > Human code review bottlenecks

**Current sweet spot:**
- Simple structure that scales
- Fast local development
- AI-powered learning
- Domain-driven design foundations

### Key Lessons Learned

1. **Start simple**: Complex solutions often solve problems you don't have yet
2. **Test what matters**: Unit tests provide better ROI than complex integration tests
3. **Local first**: Optimize for local development speed over CI complexity
4. **AI-assisted learning**: Claude provides better, faster feedback than traditional code review
5. **Domain focus**: Good architecture emerges from understanding the domain, not following patterns

**Remember**: Begin with the minimum viable solution, then add complexity incrementally based on real needs, not imagined requirements.

---

## 🚀 Happy Go Learning!

This project demonstrates that learning enterprise development doesn't require enterprise complexity from day one. 

**Start simple → Learn from real challenges → Grow skills gradually → Scale when needed**

The combination of:
- **Domain-Driven Design** (good structure)
- **Go conventions** (idiomatic code)  
- **AI assistance** (continuous learning)
- **Local-first development** (fast feedback)

...creates an excellent environment for mastering Go and modern software development practices.

**Questions?** Create an issue or mention `@claude` in a pull request - AI assistance is always available!