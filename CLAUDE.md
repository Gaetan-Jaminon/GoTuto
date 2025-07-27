# Claude Memory - GoTuto Project

## Project Overview
**Goal:** Learn Go as a .NET developer through practical Domain-Driven Design (DDD) implementation  
**Architecture:** Single module with complete domain separation (billing + catalog)  
**Current Status:** Phase 1 Complete - Billing domain fully implemented and tested

## Project Structure
```
GoTuto/
├── cmd/                          # Application entry points
│   ├── billing-api/             # Billing REST API server
│   ├── catalog-api/             # Catalog REST API server  
│   ├── billing-migrator/        # Billing database migrations
│   └── catalog-migrator/        # Catalog database migrations
├── internal/                    # Internal application code
│   ├── billing/                 # Billing domain (COMPLETE ✅)
│   │   ├── api/                # REST API handlers with dependency injection
│   │   ├── models/             # Domain models with comprehensive tests
│   │   │   └── testdata/       # Kubernetes-style external test data
│   │   ├── database/           # Database connection and config
│   │   ├── migrations/         # SQL migration files
│   │   ├── repositories/       # Data access layer (placeholder)
│   │   └── services/           # Business logic layer (placeholder)
│   ├── catalog/                # Catalog domain (NEEDS TESTING 🧪)
│   │   ├── api/               # REST API handlers implemented
│   │   ├── models/            # Domain models implemented
│   │   ├── database/          # Database connection and config
│   │   ├── migrations/        # SQL migration files
│   │   ├── repositories/      # Data access layer (placeholder)
│   │   └── services/          # Business logic layer (placeholder)
│   └── shared/                # Shared infrastructure
│       └── infrastructure/    # Config, logging, CORS, database, server
├── config/                    # Hierarchical configuration
│   ├── base/                 # Base configuration for all environments
│   ├── billing/              # Billing-specific configuration
│   └── catalog/              # Catalog-specific configuration
├── docs/                     # Project documentation
│   ├── USE_CASES_REPORT.md   # Comprehensive use case analysis
│   └── GITHUB_PROJECT_SETUP.md # Project management setup guide
└── scripts/                  # Build and development scripts
```

## Implementation Status

### ✅ COMPLETED - Billing Domain (Phase 1)
**11 Use Cases Fully Implemented and Tested:**

#### Client Management (UC-B-001 to UC-B-005)
- Create Client - Full validation, error handling
- Get Client - Include related invoices  
- Update Client - Partial updates with validation
- Delete Client - Business rule: cannot delete if has invoices
- Search Clients - Pagination, name/email search

#### Invoice Management (UC-B-006 to UC-B-011)  
- Create Invoice - Auto-numbering, client validation
- Get Invoice - Include client details
- Update Invoice - Business rule validation
- Delete Invoice - Business rule: cannot delete paid invoices
- List Invoices - Pagination, filtering by client/status
- Get Client Invoices - All invoices for specific client

**Key Features:**
- **API Layer:** Complete REST endpoints with dependency injection
- **Models:** Full validation with business rules (status transitions, overdue logic)
- **Testing:** Kubernetes-style external test data organization
- **Database:** GORM integration with schema-based isolation
- **Business Rules:** Enforced state machines and constraints

### 🧪 NEEDS TESTING - Catalog Domain (Phase 2)
**4 Use Cases Implemented, Tests Pending:**
- UC-C-001: Create Category - Basic CRUD implemented
- UC-C-002: Get Category - With product relationships  
- UC-C-003: Create Product - With category relationships
- UC-C-004: Get Product - With category details

**Missing:** Test data organization, business rule validation, comprehensive testing

### 📋 PLANNED - Cross-Domain Features (Phase 3)
**3 Use Cases Identified:**
- UC-X-001: Product Invoice - High priority cross-domain feature
- UC-X-002: Customer Product History - Medium priority analytics
- UC-X-003: Product Revenue Report - Low priority reporting

## Technical Achievements

### Architecture Decisions
- **Single Module:** Simplified from multiple modules for better learning
- **Domain-Driven Design:** Complete separation between billing and catalog
- **Dependency Injection:** API handlers receive database connections
- **Schema-based Isolation:** Each domain has its own database schema
- **Configuration Hierarchy:** Base + domain-specific + environment configs

### Testing Strategy  
- **External Test Data:** Kubernetes-style testdata/ directories
- **Path Resolution:** runtime.Caller() for reliable file loading
- **Type Safety:** Duplicate types in testdata to avoid circular imports
- **Comprehensive Coverage:** 100% billing domain test coverage
- **Business Rule Testing:** Status transitions, validation edge cases

### Build System
- **Multi-Binary:** Separate binaries for APIs and migrators
- **Docker Support:** Complete containerization with proper stages
- **OpenShift Compatible:** Non-root users, proper permissions
- **Scripts:** Automated build, test, lint, and Docker workflows

## Recent Session Summary (January 27, 2025)

### 🎯 Major Accomplishments
1. **Fixed Test Data Path Resolution** - Resolved runtime.Caller() issues for external JSON test data
2. **Completed Kubernetes-Style Testing** - Full external test data organization working
3. **Created Comprehensive Documentation** - Use case reports and project management guides
4. **Set Up GitHub Project Management** - 18 issues, labels, milestones for tracking

### 🔧 Technical Fixes
- **Path Resolution:** Fixed loadJSONFile() function to use runtime.Caller() properly
- **Test Data Validation:** Corrected "name too long" test case to exceed 100 characters
- **Binary Management:** Removed accidentally committed binary, updated .gitignore
- **All Tests Passing:** Complete billing domain test suite working

### 📊 Project Management Setup
**GitHub Organization Created:**
- **Labels:** 9 labels for domains, types, priorities, status tracking
- **Milestones:** 3 phases for release planning
- **Issues:** 18 issues (#26-43) tracking all use cases
- **Documentation:** Comprehensive reports for stakeholders

**Status Tracking:**
- Billing: 11/11 use cases complete with tests ✅
- Catalog: 4/4 use cases implemented, 0/4 tested 🧪  
- Cross-domain: 0/3 use cases implemented 📋
- **Overall Progress:** 73% use case coverage

## Next Steps Identified

### Immediate (Next Session)
1. **Catalog Domain Testing** - Add comprehensive test coverage following billing patterns
2. **Create Project Board** - Manual setup via GitHub web interface (CLI permissions needed)
3. **Business Rule Definition** - Establish catalog validation and constraints

### Short Term
1. **Repository Layer** - Abstract data access patterns
2. **Service Layer** - Complex business orchestration  
3. **Cross-Domain Use Cases** - Product invoicing implementation

### Long Term  
1. **Authentication/Authorization** - Secure API endpoints
2. **Event-Driven Architecture** - Domain events and handlers
3. **Advanced Features** - Reporting, analytics, performance optimization

## Development Patterns Established

### Code Organization
- **Domain Isolation:** Complete separation of billing and catalog
- **Dependency Injection:** Clean, testable API handlers
- **External Test Data:** JSON files in testdata/ directories
- **Business Rules:** Enforced through model validation and state machines

### Testing Patterns
- **Test Data Loading:** runtime.Caller() for reliable path resolution
- **Type Duplication:** Avoid circular imports in test packages
- **Edge Case Coverage:** Comprehensive validation scenarios
- **Business Logic Testing:** Status transitions, overdue detection

### Git Workflow
- **Feature Branches:** Clean separation of work
- **Protected Main:** Pull request workflow for quality
- **Descriptive Commits:** Clear history with Claude Code attribution
- **Documentation:** Comprehensive tracking of decisions and progress

## Key Learning Outcomes (Go for .NET Developer)

### Go-Specific Concepts Learned
- **testdata/ Directories:** Go convention for external test data (excluded from builds)
- **Package Organization:** Clean domain separation without circular dependencies
- **GORM Integration:** Database ORM patterns in Go ecosystem
- **Dependency Injection:** Function-based DI pattern for handlers
- **Build Tags:** Conditional compilation for different environments

### DDD Implementation in Go
- **Domain Models:** Rich entities with business logic
- **Value Objects:** Status enums with validation
- **Business Rules:** Enforced through model methods
- **Repository Pattern:** Interface-based data access (planned)
- **Service Layer:** Orchestration of complex business operations (planned)

---

*Last Updated: January 27, 2025*  
*Session: Kubernetes-style test data completion + GitHub project management setup*  
*Next Focus: Catalog domain testing and project board creation*