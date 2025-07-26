# Go Testing - Comprehensive Guide

A complete guide to testing in Go, covering patterns, best practices, and comparisons with .NET testing approaches.

## Testing Philosophy in Go

### Built-in Testing Framework
Unlike .NET which requires external testing frameworks (NUnit, xUnit, MSTest), Go has **testing built into the language**:
- No external dependencies needed for basic testing
- `go test` command is part of the standard toolchain
- Simple, consistent conventions

### Core Principles
1. **Simplicity** - Minimal setup, maximum clarity
2. **Convention over Configuration** - Standard naming and structure
3. **Fast Feedback** - Quick test execution
4. **Table-Driven Tests** - Go idiom for testing multiple scenarios

## Go vs .NET Testing Comparison

| Aspect | .NET | Go |
|--------|------|-----|
| **Framework** | NUnit/xUnit/MSTest | Built-in `testing` package |
| **Test Discovery** | Attributes `[Test]`, `[Fact]` | Function naming `TestXxx` |
| **Assertions** | Rich assertion libraries | Simple `t.Error()`, `t.Fatal()` |
| **Mocking** | Moq, NSubstitute | Interfaces + manual mocks |
| **Setup/Teardown** | `[SetUp]`/`[TearDown]` | `TestMain()` function |
| **Parameterized Tests** | `[TestCase]` attributes | Table-driven tests |
| **Test Categories** | `[Category]` attributes | Build tags |

## Test Types and Structure

### 1. Unit Tests

#### Basic Test Structure
```go
// client_test.go (same package as client.go)
package models

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestClient_Validation(t *testing.T) {
    // Arrange
    client := Client{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    // Act
    err := client.Validate()
    
    // Assert
    assert.NoError(t, err)
}
```

#### Table-Driven Tests (Go Idiom)
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name     string
        email    string
        expected bool
    }{
        {"valid email", "test@example.com", true},
        {"invalid email", "invalid", false},
        {"empty email", "", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ValidateEmail(tt.email)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 2. Integration Tests
```go
// tests/integration/database_test.go
package integration

import (
    "testing"
    "github.com/stretchr/testify/suite"
)

type DatabaseIntegrationSuite struct {
    suite.Suite
    db *gorm.DB
}

func (suite *DatabaseIntegrationSuite) SetupSuite() {
    // Setup test database
    db, err := setupTestDatabase()
    suite.NoError(err)
    suite.db = db
}

func (suite *DatabaseIntegrationSuite) TearDownSuite() {
    // Cleanup
    suite.db.Close()
}

func (suite *DatabaseIntegrationSuite) TestClientCRUD() {
    // Test database operations
    client := &Client{Name: "Test", Email: "test@example.com"}
    
    // Create
    err := suite.db.Create(client).Error
    suite.NoError(err)
    suite.NotZero(client.ID)
    
    // Read
    var found Client
    err = suite.db.First(&found, client.ID).Error
    suite.NoError(err)
    suite.Equal(client.Name, found.Name)
}

func TestDatabaseIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration tests")
    }
    suite.Run(t, new(DatabaseIntegrationSuite))
}
```

### 3. End-to-End Tests
```go
// tests/e2e/api_test.go
package e2e

func (suite *E2ETestSuite) TestClientLifecycle() {
    // 1. Create client via API
    clientData := CreateClientRequest{
        Name:  "E2E Test Client",
        Email: "e2e@example.com",
    }
    
    jsonData, _ := json.Marshal(clientData)
    req, _ := http.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(jsonData))
    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)
    
    suite.Equal(http.StatusCreated, w.Code)
    
    // 2. Verify client was created
    var client Client
    json.Unmarshal(w.Body.Bytes(), &client)
    suite.NotZero(client.ID)
    
    // 3. Update client
    // 4. Delete client
    // 5. Verify deletion
}
```

## Testing Libraries and Tools

### 1. **Testify** (Most Popular)
```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
)

// Assertions
assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.Contains(t, slice, item)

// Mocking
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) GetClient(id uint) (*Client, error) {
    args := m.Called(id)
    return args.Get(0).(*Client), args.Error(1)
}

// Test Suites
type ClientTestSuite struct {
    suite.Suite
    mockRepo *MockRepository
}
```

### 2. **Built-in Testing**
```go
// Simple assertions without external dependencies
func TestBasic(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Expected 5, got %d", result)
    }
}

// Subtests
func TestMath(t *testing.T) {
    t.Run("addition", func(t *testing.T) {
        // Test addition
    })
    
    t.Run("subtraction", func(t *testing.T) {
        // Test subtraction
    })
}
```

### 3. **Test Database Strategies**

#### In-Memory SQLite
```go
func SetupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    assert.NoError(t, err)
    
    // Run migrations
    err = db.AutoMigrate(&Client{}, &Invoice{})
    assert.NoError(t, err)
    
    return db
}
```

#### Test Containers (Real Database)
```go
// Using testcontainers-go for real PostgreSQL
func SetupPostgresContainer(t *testing.T) *gorm.DB {
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_PASSWORD": "password",
            "POSTGRES_DB":       "testdb",
        },
    }
    
    postgres, err := testcontainers.GenericContainer(ctx, 
        testcontainers.GenericContainerRequest{
            ContainerRequest: req,
            Started:          true,
        })
    require.NoError(t, err)
    
    // Get connection details and connect
    // ...
}
```

## Testing Patterns and Best Practices

### 1. **Test Organization**
```
api/billing/
├── internal/
│   ├── models/
│   │   ├── client.go
│   │   └── client_test.go        # Unit tests alongside source
│   ├── handlers/
│   │   ├── client.go
│   │   └── client_test.go
│   └── database/
│       ├── repository.go
│       └── repository_test.go
├── tests/
│   ├── integration/              # Integration tests
│   │   └── database_test.go
│   └── e2e/                      # End-to-end tests
│       └── api_test.go
├── testdata/                     # Test fixtures
│   ├── clients.json
│   └── invoices.json
└── pkg/
    └── testutil/                 # Test utilities
        └── testutil.go
```

### 2. **Test Utilities and Helpers**
```go
// pkg/testutil/testutil.go
package testutil

// SetupTestDB creates test database
func SetupTestDB(t *testing.T) *gorm.DB { ... }

// LoadTestData loads fixtures
func LoadTestData(t *testing.T, filename string) []byte { ... }

// CreateTestClient creates test client with overrides
func CreateTestClient(overrides ...func(*Client)) *Client {
    client := &Client{
        Name:  "Test Client",
        Email: "test@example.com",
    }
    
    for _, override := range overrides {
        override(client)
    }
    
    return client
}

// Usage
client := CreateTestClient(func(c *Client) {
    c.Name = "Custom Name"
    c.Email = "custom@example.com"
})
```

### 3. **Mocking Strategies**

#### Interface-Based Mocking
```go
// Define interface
type ClientRepository interface {
    GetClient(id uint) (*Client, error)
    CreateClient(client *Client) error
}

// Production implementation
type DBClientRepository struct {
    db *gorm.DB
}

func (r *DBClientRepository) GetClient(id uint) (*Client, error) {
    var client Client
    err := r.db.First(&client, id).Error
    return &client, err
}

// Test mock
type MockClientRepository struct {
    clients map[uint]*Client
}

func (m *MockClientRepository) GetClient(id uint) (*Client, error) {
    client, exists := m.clients[id]
    if !exists {
        return nil, gorm.ErrRecordNotFound
    }
    return client, nil
}
```

#### Dependency Injection for Testing
```go
type ClientService struct {
    repo ClientRepository
}

func NewClientService(repo ClientRepository) *ClientService {
    return &ClientService{repo: repo}
}

// Test with mock
func TestClientService_GetClient(t *testing.T) {
    mockRepo := &MockClientRepository{
        clients: map[uint]*Client{
            1: {ID: 1, Name: "Test Client"},
        },
    }
    
    service := NewClientService(mockRepo)
    client, err := service.GetClient(1)
    
    assert.NoError(t, err)
    assert.Equal(t, "Test Client", client.Name)
}
```

### 4. **HTTP Handler Testing**
```go
func TestCreateClient(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.POST("/clients", CreateClient)
    
    // Test data
    clientData := CreateClientRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }
    
    jsonData, _ := json.Marshal(clientData)
    
    // Request
    req, _ := http.NewRequest("POST", "/clients", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var client Client
    err := json.Unmarshal(w.Body.Bytes(), &client)
    assert.NoError(t, err)
    assert.Equal(t, clientData.Name, client.Name)
}
```

### 5. **Test Data Management**

#### Test Fixtures
```go
// testdata/clients.json
[
  {
    "id": 1,
    "name": "Acme Corporation",
    "email": "contact@acme.com"
  }
]

// Load in tests
func LoadTestClients(t *testing.T) []Client {
    data, err := ioutil.ReadFile("../../testdata/clients.json")
    assert.NoError(t, err)
    
    var clients []Client
    err = json.Unmarshal(data, &clients)
    assert.NoError(t, err)
    
    return clients
}
```

#### Factory Functions
```go
func NewTestClient() *Client {
    return &Client{
        Name:  "Test Client",
        Email: "test@example.com",
    }
}

func NewTestInvoice(clientID uint) *Invoice {
    return &Invoice{
        Number:   "INV-TEST-001",
        ClientID: clientID,
        Amount:   100.00,
        Status:   InvoiceStatusDraft,
    }
}
```

## Test Execution and CI/CD

### 1. **Test Commands**
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run only fast tests (skip integration)
go test -short ./...

# Run specific test
go test -run TestClientCRUD ./...

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

### 2. **Test Categories with Build Tags**
```go
// +build integration

package integration

func TestDatabaseIntegration(t *testing.T) {
    // Integration test code
}
```

```bash
# Run only integration tests
go test -tags=integration ./...

# Skip integration tests
go test ./...
```

### 3. **Makefile for Testing**
```makefile
# Run all tests
test:
	go test -v -timeout 30s ./...

# Run unit tests only
test-unit:
	go test -v -short ./internal/...

# Run with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run integration tests
test-integration:
	go test -v -tags=integration ./tests/integration/...

# Run in CI
test-ci: test-race test-coverage
	@echo "CI tests completed"
```

## Performance Testing

### 1. **Benchmarks**
```go
func BenchmarkClientCreate(b *testing.B) {
    db := setupTestDB(b)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        client := &Client{
            Name:  fmt.Sprintf("Client %d", i),
            Email: fmt.Sprintf("client%d@example.com", i),
        }
        db.Create(client)
    }
}

func BenchmarkGetClients(b *testing.B) {
    router := setupTestRouter()
    req, _ := http.NewRequest("GET", "/clients", nil)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
    }
}
```

### 2. **Load Testing**
```go
func TestConcurrentClientCreation(t *testing.T) {
    const numGoroutines = 100
    const clientsPerGoroutine = 10
    
    var wg sync.WaitGroup
    wg.Add(numGoroutines)
    
    for i := 0; i < numGoroutines; i++ {
        go func(routineID int) {
            defer wg.Done()
            
            for j := 0; j < clientsPerGoroutine; j++ {
                client := &Client{
                    Name:  fmt.Sprintf("Client %d-%d", routineID, j),
                    Email: fmt.Sprintf("client%d-%d@example.com", routineID, j),
                }
                
                err := db.Create(client).Error
                assert.NoError(t, err)
            }
        }(i)
    }
    
    wg.Wait()
    
    // Verify all clients were created
    var count int64
    db.Model(&Client{}).Count(&count)
    assert.Equal(t, int64(numGoroutines*clientsPerGoroutine), count)
}
```

## Error Testing and Edge Cases

### 1. **Error Scenarios**
```go
func TestClientService_ErrorHandling(t *testing.T) {
    tests := []struct {
        name        string
        clientID    uint
        mockError   error
        expectedErr string
    }{
        {
            name:        "client not found",
            clientID:    999,
            mockError:   gorm.ErrRecordNotFound,
            expectedErr: "client not found",
        },
        {
            name:        "database connection error",
            clientID:    1,
            mockError:   errors.New("connection failed"),
            expectedErr: "database error",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &MockClientRepository{}
            mockRepo.On("GetClient", tt.clientID).Return(nil, tt.mockError)
            
            service := NewClientService(mockRepo)
            client, err := service.GetClient(tt.clientID)
            
            assert.Nil(t, client)
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.expectedErr)
        })
    }
}
```

### 2. **Boundary Testing**
```go
func TestClientValidation_BoundaryValues(t *testing.T) {
    tests := []struct {
        name     string
        client   Client
        valid    bool
    }{
        {"minimum valid name", Client{Name: "A", Email: "a@b.co"}, true},
        {"maximum valid name", Client{Name: strings.Repeat("A", 100), Email: "test@example.com"}, true},
        {"name too long", Client{Name: strings.Repeat("A", 101), Email: "test@example.com"}, false},
        {"empty name", Client{Name: "", Email: "test@example.com"}, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.client.Validate()
            
            if tt.valid {
                assert.NoError(t, err)
            } else {
                assert.Error(t, err)
            }
        })
    }
}
```

## Key Takeaways for .NET Developers

### 1. **Philosophy Differences**
- **Go**: Simple, explicit, minimal setup
- **.NET**: Rich tooling, attributes, extensive frameworks

### 2. **Test Organization**
- **Go**: Tests in same package as source code
- **.NET**: Separate test projects

### 3. **Assertions**
- **Go**: Simple `t.Error()` or testify for rich assertions
- **.NET**: Rich assertion libraries built-in

### 4. **Mocking**
- **Go**: Interface-based, manual mocks or testify/mock
- **.NET**: Reflection-based mocking frameworks

### 5. **Test Discovery**
- **Go**: Function naming convention (`TestXxx`)
- **.NET**: Attributes (`[Test]`, `[Fact]`)

### 6. **Best Practices**
1. **Keep tests simple and focused**
2. **Use table-driven tests for multiple scenarios**
3. **Test interfaces, not implementations**
4. **Use dependency injection for testability**
5. **Separate unit, integration, and e2e tests**
6. **Use build tags for test categories**
7. **Test error cases and edge conditions**
8. **Use test utilities to reduce duplication**

Go testing is simpler but more explicit than .NET testing. The trade-off is less "magic" but more understanding of what your tests are actually doing. The built-in tooling is excellent and provides everything needed for comprehensive testing strategies.