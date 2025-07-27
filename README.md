# GoTuto - Learn Go with Real-World Examples

A practical Go learning project designed for .NET developers transitioning to Go. This project demonstrates enterprise development practices while following the principle: **start simple, add complexity gradually**.

## 🎯 Learning Philosophy

This project evolved through real development challenges, demonstrating:
- **Start with basics**: Simple unit tests, core functionality
- **Add complexity incrementally**: CI/CD, then advanced testing
- **Learn from failures**: Document troubleshooting and decisions
- **Focus on working solutions**: Practical over perfect

## 📁 Current Project Structure

```
GoTuto/
├── api/
│   └── billing/                 # Main CRUD API service
│       ├── cmd/main.go         # Application entry point
│       ├── internal/           # Private application code
│       │   ├── config/         # Configuration management (Viper)
│       │   ├── database/       # Database connection
│       │   ├── handlers/       # HTTP handlers (Gin)
│       │   └── models/         # Data models (GORM)
│       ├── deployments/        # OpenShift deployment manifests
│       └── Dockerfile          # Container build
├── .github/workflows/          # CI/CD automation
│   ├── ci.yml                 # Core CI pipeline
│   ├── claude-code-review.yml # AI-powered code review
│   └── dependabot-ci.yml      # Dependency management
├── notes/                      # Learning documentation
└── README.md                   # This file
```

## 🚀 Quick Start

### Prerequisites
- **Go 1.22+**
- **PostgreSQL 15+** (for full API functionality)
- **Git**

### 1. Clone and Test

```bash
git clone https://github.com/Gaetan-Jaminon/GoTuto.git
cd GoTuto/api/billing

# Run unit tests (no dependencies required)
go test ./...

# Build the application
go build -o billing cmd/main.go
```

### 2. Run with Database (Optional)

```bash
# Start PostgreSQL (using Docker/Podman)
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

# Run the API
./billing
```

### 3. Test the API

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/clients
```

## 🧪 Testing Strategy (Simplified)

We learned the hard way that complex testing setups can become blockers. Current approach:

### Unit Tests Only (3 files)
- **`internal/config/config_test.go`**: Configuration validation
- **`internal/models/client_test.go`**: Client model logic  
- **`internal/models/invoice_test.go`**: Invoice model logic

```bash
# Run all tests (fast, no dependencies)
go test ./...

# With coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Why We Removed Complex Tests
- ❌ **Integration tests**: Required PostgreSQL setup (barrier to entry)
- ❌ **E2E tests**: Required full application stack (complex)
- ❌ **Handler tests**: Required database mocking (over-engineering)

**Lesson**: Start with unit tests that work everywhere, add complexity when the core is stable.

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

## 🤖 AI-Powered Development

### Claude Code Review Integration

One of the key innovations in this project is automated AI code review:

**How it works:**
1. Create pull request
2. CI runs tests first (saves tokens)
3. If tests pass → Claude reviews code automatically
4. Provides suggestions, catches issues, explains patterns

**Benefits:**
- Learn Go best practices through AI feedback
- Catch issues early in development
- Get explanations for complex patterns
- Available 24/7 without human reviewers

**Token Optimization:**
- Claude only runs after CI success (saves costs)
- Focuses review on working code
- Provides meaningful feedback vs debugging compilation errors

See `notes/claude-github-integration.md` for detailed setup and troubleshooting.

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

## 🔄 CI/CD Pipeline (Simplified)

### Current Workflows

**1. Continuous Integration (`ci.yml`)**
- Build and test Go code
- Security scanning with Trivy
- Quality gates
- Triggers Claude review on success

**2. Claude Code Review (`claude-code-review.yml`)**
- AI-powered code review
- Only runs after CI passes (token optimization)
- Provides learning feedback

**3. Dependabot CI (`dependabot-ci.yml`)**
- Automated dependency updates
- Simplified checks for dependency PRs

### What We Removed (And Why)
- ❌ **Complex deployment workflows**: Too early in learning process
- ❌ **Multiple environments**: Start with simple local development
- ❌ **Release automation**: Focus on core functionality first
- ❌ **Security scanning workflows**: Basic security in main CI is sufficient

**Lesson**: Start with essential CI (build/test), add deployment complexity later.

## 🛡️ Branch Protection

We learned the hard way that automated GitHub settings can be unreliable. Current approach:

### Manual Branch Protection Setup

**Protection Rules:**
- Require pull request reviews (1 approval)
- Require status checks: `Continuous Integration`, `claude-review`
- Allow administrators to bypass (for learning)

**Why Manual?**
- GitHub Settings app had reliability issues
- Manual setup works immediately
- Simpler for learning projects
- Can automate later when the project matures

See `notes/branch-protection-setup.md` for step-by-step setup.

## 🐳 Container Support

### Multi-stage Dockerfile

Optimized for OpenShift deployment:
- Red Hat UBI base images
- Non-root user (UID 1001)
- Security-conscious build process

```bash
# Build image
docker build -t billing:latest .

# Run container
docker run -p 8080:8080 \
  -e BILLING_DATABASE_HOST=host.docker.internal \
  -e BILLING_DATABASE_PASSWORD=password \
  billing:latest
```

## 🎓 Learning Resources

### Go Concepts Demonstrated

**Core Language:**
- Package organization and Go modules
- Error handling (no exceptions)
- Interfaces and composition
- Struct methods and receivers
- Short variable declaration (`:=`)

**Web Development:**
- HTTP handlers and routing (Gin)
- JSON marshaling/unmarshaling
- Middleware patterns
- Request validation

**Database Integration:**
- ORM patterns with GORM
- Database connection management
- Configuration-based connection strings

**Testing:**
- Table-driven tests (Go idiom)
- Test organization and structure
- Coverage measurement

**DevOps:**
- Configuration management
- Container deployment
- CI/CD automation
- AI-powered code review

### Notes for .NET Developers

Key differences when coming from .NET:

| .NET Concept | Go Equivalent | Notes |
|-------------|---------------|--------|
| `try/catch` | Explicit error returns | `if err != nil { return err }` |
| Classes | Structs + methods | No inheritance, use composition |
| `null` | Zero values | Prefer zero values over pointers |
| NuGet | Go modules | Simpler dependency management |
| MSBuild | `go build` | No complex build configuration |
| LINQ | Manual iteration | More explicit, less magic |

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
   go test ./...  # Fast feedback loop
   go run cmd/main.go  # Test manually
   ```

3. **Create pull request**:
   - CI runs automatically
   - Claude provides code review
   - Merge after approval

4. **Learn from feedback**:
   - Review Claude's suggestions
   - Update code based on feedback
   - Document new patterns learned

## 🔍 Troubleshooting

### Common Issues

**1. Tests fail with "connection refused"**
- Solution: Tests are now unit tests only, no database required
- If you see this, you might have old integration tests

**2. CI fails with "unused import"**
- Solution: Run `go mod tidy` and remove unused imports
- Use `goimports` to auto-manage imports

**3. Claude review doesn't run**
- Check that CI passed first (Claude only runs after success)
- Verify GitHub Actions are enabled in repository settings

**4. Branch protection not working**
- Use manual setup instead of GitHub Settings app
- Follow guide in `notes/branch-protection-setup.md`

### Getting Help

1. Check the `notes/` directory for specific topics
2. Review AI feedback in pull requests
3. Test locally first: `go test ./...`
4. Use `git bisect` to find breaking changes

## 🎯 Next Steps

### Immediate Learning Goals
1. **Master Go basics**: Complete all examples in `notes/`
2. **Practice API development**: Add new endpoints
3. **Understand testing**: Write unit tests for new features
4. **Learn from AI feedback**: Create PRs and review Claude's suggestions

### Future Enhancements (When Ready)
1. **Add integration tests**: When comfortable with unit testing
2. **Implement authentication**: JWT or session-based auth
3. **Add more complex business logic**: Validation rules, calculations
4. **Deploy to cloud**: OpenShift, Kubernetes, or cloud providers
5. **Add monitoring**: Metrics, logging, tracing

### Advanced Topics to Explore Later
- Goroutines and channels (concurrency)
- Advanced database patterns
- Microservices communication
- Performance optimization
- Security hardening

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Write tests for new functionality
4. Ensure tests pass: `go test ./...`
5. Commit with descriptive messages
6. Push and create pull request
7. Learn from Claude's code review feedback

## 📖 Philosophy: Start Simple, Iterate

This project embodies a key software development principle:

> **"Make it work, make it right, make it fast"** - Kent Beck

We started with enterprise complexity and learned to simplify:
- ✅ Working unit tests > Complex integration tests
- ✅ Manual setup > Unreliable automation  
- ✅ Essential CI > Over-engineered pipelines
- ✅ Local development > Complex deployment

**Key Lesson**: Begin with the minimum viable solution, then add complexity incrementally based on real needs, not imagined requirements.

---

**Happy Go learning! 🚀**

This project demonstrates that learning enterprise development doesn't require enterprise complexity from day one. Start simple, learn from real challenges, and grow your skills gradually.

For questions or suggestions, create an issue or pull request - Claude will help review your contributions!