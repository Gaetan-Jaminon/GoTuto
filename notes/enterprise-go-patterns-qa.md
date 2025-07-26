# Enterprise Go Patterns - Q&A Summary

This document summarizes key enterprise patterns and best practices learned during the development of a Go billing service project.

## 1. Configuration Management in Go

### Question: Configuration Files for Multiple Environments

**Q**: "I don't see any configuration files. How will that happen in multiple environments? I suppose it will be injected during deployment but we will still need a template. What is the recommended way and best practice to do that?"

### Answer: Viper + Hierarchical Configuration

**Best Practice**: Use **Viper** library with hierarchical configuration management.

#### Configuration Hierarchy (highest to lowest priority):
1. **Command-line flags**
2. **Environment variables**
3. **Config files** (environment-specific)
4. **Default values**

#### File Structure:
```
config/
├── config.yaml          # Base defaults
├── config.dev.yaml      # Development overrides
├── config.prod.yaml     # Production overrides
└── config.test.yaml     # Test environment
```

#### Environment Variable Pattern:
- **Prefix**: `BILLING_`
- **Mapping**: `BILLING_DATABASE_HOST` → `database.host`
- **Replacement**: Dots become underscores

#### Example Implementation:
```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Logging  LoggingConfig  `mapstructure:"logging"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    
    // Load base config
    viper.ReadInConfig()
    
    // Load environment-specific overrides
    env := viper.GetString("APP_ENV")
    viper.SetConfigName(fmt.Sprintf("config.%s", env))
    viper.MergeInConfig()
    
    // Enable environment variable overrides
    viper.SetEnvPrefix("BILLING")
    viper.AutomaticEnv()
    
    var config Config
    return &config, viper.Unmarshal(&config)
}
```

#### Deployment Patterns:
- **Development**: Config files + `.env` for local overrides
- **Production**: ConfigMaps + Secrets for sensitive data
- **Container**: Environment variables injected at runtime

---

## 2. Database Migrations as Separate Modules

### Question: Dedicated DB Migration Packages

**Q**: "Don't DB migrations have a dedicated package/module sometimes?"

### Answer: Yes - Enterprise Best Practice

**Why Separate Migration Modules?**

#### Benefits:
1. **Clean Separation** - Migrations aren't part of runtime code
2. **CI/CD Friendly** - Run migrations as separate deployment step
3. **Version Control** - Track schema changes independently
4. **Different Permissions** - Migration tool needs DDL rights, API needs DML only
5. **Rollback Safety** - Can version and test migrations separately

#### Enterprise Patterns:
- **Init Containers** in K8s/OpenShift - Run migrations before app starts
- **Separate Deployment** - Deploy migrations first, then API
- **Schema Registry** - Central place for all DB schemas
- **GitOps** - Migrations triggered by Git commits

#### Implementation Structure:
```
api/
├── billing/                    # Main API service
│   ├── internal/
│   ├── cmd/
│   └── ...
└── billing-dbmigrations/       # Separate migration module
    ├── go.mod                  # Independent module
    ├── cmd/migrate/            # CLI tool with Cobra
    ├── cmd/health/             # Health check server
    ├── migrations/             # SQL files
    └── ...
```

#### Migration CLI Features:
```bash
billing-migrate up              # Apply migrations
billing-migrate down --steps=3 # Rollback 3 migrations
billing-migrate create add_index # Create new migration
billing-migrate version         # Check current version
billing-migrate force 5         # Force version (dirty state)
```

---

## 3. Project Structure: From Demo to Enterprise

### Transformation: demo01 → billing

#### What Changed:
- **Module Name**: `gotuto/api/demo01` → `gotuto/api/billing`
- **Environment Prefix**: `DEMO01_` → `BILLING_`
- **Database Name**: `demo01` → `billing`
- **Service Context**: Generic demo → Specific billing domain

#### Enterprise Structure Achieved:
```
api/
├── billing/                    # Main service
│   ├── cmd/main.go            # Entry point
│   ├── internal/              # Business logic
│   │   ├── config/           # Configuration management
│   │   ├── models/           # Data models
│   │   ├── handlers/         # HTTP handlers
│   │   └── database/         # DB connection
│   ├── config/               # Config files
│   ├── deployments/          # Deployment manifests
│   │   └── openshift/        # OpenShift specific
│   └── Containerfile         # Podman build
└── billing-dbmigrations/      # Migration module
    ├── cmd/                  # CLI tools
    ├── migrations/           # SQL files
    └── Containerfile         # Migration container
```

#### Why This Structure?
- **Domain-Driven** - Clear business context (billing)
- **Microservice Ready** - Independent deployable units
- **Enterprise Standard** - Matches real-world Go projects
- **Platform Specific** - Optimized for OpenShift/RHEL

---

## 4. Podman/OpenShift vs Docker/Kubernetes

### Environment Choice: RHEL + Podman + OpenShift

#### Key Differences:

| Aspect | Docker/K8s | Podman/OpenShift |
|--------|------------|------------------|
| **Base Images** | Alpine/Ubuntu | Red Hat UBI |
| **Container Runtime** | Docker | Podman (rootless) |
| **Build Files** | Dockerfile | Containerfile |
| **Orchestration** | Kubernetes | OpenShift (K8s + more) |
| **Networking** | Ingress | Routes |
| **Deployments** | Deployment | DeploymentConfig |
| **Security** | Manual setup | Built-in SCCs |

#### OpenShift-Specific Features:
- **Routes** instead of Ingress
- **DeploymentConfig** with rolling updates
- **BuildConfig** for Source-to-Image builds
- **ImageStreams** for image management
- **Security Context Constraints** (SCC)

#### Example Containerfile (vs Dockerfile):
```dockerfile
# Red Hat UBI instead of Alpine
FROM registry.access.redhat.com/ubi9/go-toolset:1.20 AS builder
FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

# OpenShift-compatible user (group 0)
RUN useradd -u 1001 -r -g 0 -s /sbin/nologin appuser
RUN chown -R 1001:0 /app && chmod -R g=u /app
USER 1001
```

#### Deployment Patterns:
```yaml
# OpenShift DeploymentConfig with init containers
apiVersion: apps.openshift.io/v1
kind: DeploymentConfig
spec:
  template:
    spec:
      initContainers:
      - name: migrate
        image: billing-migrations:latest
        command: ["billing-migrate", "up"]
      containers:
      - name: api
        image: billing-api:latest
```

---

## 5. Health Endpoints in Migration Tools

### Question: Why Health in DB Migrations?

**Q**: "Why Health in the dbmigration?"

### Answer: Enterprise Observability Pattern

#### Why Health Endpoints Everywhere?

#### 1. **Init Container Pattern**
```yaml
initContainers:
- name: migrate
  image: billing-migrations:latest
  # Health endpoint allows external monitoring
  # of migration progress and status
```

#### 2. **Enterprise Operations**
- **Monitor everything** from centralized dashboards
- **Debug without SSH access** (security policy)
- **Automate incident response** based on health status
- **Track SLOs** across all components

#### 3. **Real-World Scenarios**
```bash
# Operations team checks migration status during deployment
curl https://migration-health.billing.svc/health

# CI/CD pipeline waits for migration completion
while ! curl -f migration-health:8081/ready; do
  echo "Waiting for migrations..."
  sleep 5
done

# Monitoring alerts if migrations are stuck
if migration_health_status != "ready" for 5m:
  alert: "Database migrations stuck in billing service"
```

#### 4. **Health Endpoint Implementation**
```go
type HealthResponse struct {
    Status   string `json:"status"`
    Database struct {
        Connected bool   `json:"connected"`
        Version   string `json:"version,omitempty"`
        Error     string `json:"error,omitempty"`
    } `json:"database"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    response := HealthResponse{Status: "healthy"}
    
    // Check database connection
    db, err := getDB()
    if err != nil {
        response.Database.Connected = false
        response.Database.Error = err.Error()
    } else {
        // Get current migration version
        var version int
        db.QueryRow("SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1").Scan(&version)
        response.Database.Version = fmt.Sprintf("%d", version)
    }
    
    json.NewEncoder(w).Encode(response)
}
```

#### 5. **Industry Standards**
- **12-Factor Apps** - Treat logs and health as first-class citizens
- **Cloud Native** - Everything should be observable
- **DevOps Culture** - "If you build it, you monitor it"
- **Platform Conventions** - Kubernetes/OpenShift expect health endpoints

---

## Key Takeaways

### 1. **Think Enterprise From Day One**
- Configuration management is not optional
- Every component needs observability
- Follow platform conventions (OpenShift patterns)
- Design for operations teams

### 2. **Separation of Concerns**
- Migrations separate from application code
- Configuration separate from business logic
- Platform-specific deployments in dedicated folders

### 3. **Production Patterns**
- Use industry-standard libraries (Viper, Cobra)
- Health endpoints on everything
- Proper security contexts for containers
- Environment-specific configuration

### 4. **Go Enterprise Architecture**
```
project/
├── notes/                     # Learning documentation
├── api/
│   ├── billing/              # Main service (domain-focused)
│   │   ├── internal/         # Business logic
│   │   ├── config/           # Configuration
│   │   ├── deployments/      # Platform-specific deploys
│   │   └── cmd/             # Entry points
│   └── billing-dbmigrations/ # Separate migration module
│       ├── cmd/             # CLI tools
│       └── migrations/      # SQL files
```

### 5. **Learning Philosophy**
- Learn enterprise patterns, not just language syntax
- Understand why practices exist (observability, security, operations)
- Focus on real-world deployment scenarios
- Think about the full software lifecycle

This structure and these patterns reflect what you'd see in production Go microservices at major enterprises using OpenShift/RHEL environments.