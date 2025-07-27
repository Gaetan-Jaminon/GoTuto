package integration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"gotuto/api/billing/internal/config"
	"gotuto/api/billing/internal/database"
	"gotuto/api/billing/internal/models"

	"github.com/stretchr/testify/suite"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

// DatabaseIntegrationSuite tests database operations with a real database
type DatabaseIntegrationSuite struct {
	suite.Suite
	db     *gorm.DB
	config *config.Config
}

func (suite *DatabaseIntegrationSuite) SetupSuite() {
	// Setup test database configuration
	suite.config = &config.Config{
		Database: config.DatabaseConfig{
			Host:     getEnv("TEST_DB_HOST", "localhost"),
			Port:     5432,
			Username: getEnv("TEST_DB_USER", "postgres"),
			Password: getEnv("TEST_DB_PASSWORD", "password"),
			Name:     getEnv("TEST_DB_NAME", "billing_test"),
			SSLMode:  "disable",
		},
	}

	// Create test database if it doesn't exist
	err := suite.createTestDatabase()
	if err != nil {
		suite.T().Fatalf("Failed to create test database: %v", err)
	}

	// Connect to test database
	db, err := database.Connect(suite.config)
	if err != nil {
		suite.T().Fatalf("Failed to connect to test database: %v", err)
	}
	suite.db = db

	// Run migrations
	err = database.AutoMigrate(db)
	if err != nil {
		suite.T().Fatalf("Failed to run migrations: %v", err)
	}
}

func (suite *DatabaseIntegrationSuite) TearDownSuite() {
	// Clean up
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

func (suite *DatabaseIntegrationSuite) SetupTest() {
	// Clean up data before each test
	suite.db.Exec("TRUNCATE TABLE invoices, clients RESTART IDENTITY CASCADE")
}

func (suite *DatabaseIntegrationSuite) TestClientCRUD() {
	// Create
	client := &models.Client{
		Name:    "John Doe",
		Email:   "john@example.com",
		Phone:   "+1234567890",
		Address: "123 Main St",
	}

	err := suite.db.Create(client).Error
	suite.NoError(err)
	suite.NotZero(client.ID)
	suite.NotZero(client.CreatedAt)

	// Read
	var foundClient models.Client
	err = suite.db.First(&foundClient, client.ID).Error
	suite.NoError(err)
	suite.Equal(client.Name, foundClient.Name)
	suite.Equal(client.Email, foundClient.Email)

	// Update
	foundClient.Name = "Jane Doe"
	err = suite.db.Save(&foundClient).Error
	suite.NoError(err)

	var updatedClient models.Client
	err = suite.db.First(&updatedClient, client.ID).Error
	suite.NoError(err)
	suite.Equal("Jane Doe", updatedClient.Name)

	// Delete (soft delete)
	err = suite.db.Delete(&foundClient).Error
	suite.NoError(err)

	// Verify soft delete
	var deletedClient models.Client
	err = suite.db.First(&deletedClient, client.ID).Error
	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)

	// Verify record still exists with Unscoped
	err = suite.db.Unscoped().First(&deletedClient, client.ID).Error
	suite.NoError(err)
	suite.True(deletedClient.DeletedAt.Valid)
}

func (suite *DatabaseIntegrationSuite) TestInvoiceCRUD() {
	// First create a client
	client := &models.Client{
		Name:  "Test Client",
		Email: "test@example.com",
	}
	err := suite.db.Create(client).Error
	suite.NoError(err)

	// Create invoice
	invoice := &models.Invoice{
		Number:      "INV-001",
		ClientID:    client.ID,
		Amount:      150.75,
		Status:      models.InvoiceStatusDraft,
		Description: "Test invoice",
	}

	err = suite.db.Create(invoice).Error
	suite.NoError(err)
	suite.NotZero(invoice.ID)

	// Read with preload
	var foundInvoice models.Invoice
	err = suite.db.Preload("Client").First(&foundInvoice, invoice.ID).Error
	suite.NoError(err)
	suite.Equal(invoice.Amount, foundInvoice.Amount)
	suite.Equal(client.Name, foundInvoice.Client.Name)

	// Update status
	foundInvoice.Status = models.InvoiceStatusSent
	err = suite.db.Save(&foundInvoice).Error
	suite.NoError(err)

	// Verify update
	var updatedInvoice models.Invoice
	err = suite.db.First(&updatedInvoice, invoice.ID).Error
	suite.NoError(err)
	suite.Equal(models.InvoiceStatusSent, updatedInvoice.Status)
}

func (suite *DatabaseIntegrationSuite) TestClientInvoiceRelationship() {
	// Create client
	client := &models.Client{
		Name:  "Business Client",
		Email: "business@example.com",
	}
	err := suite.db.Create(client).Error
	suite.NoError(err)

	// Create multiple invoices for the client
	invoices := []models.Invoice{
		{
			Number:   "INV-001",
			ClientID: client.ID,
			Amount:   100.00,
			Status:   models.InvoiceStatusDraft,
		},
		{
			Number:   "INV-002",
			ClientID: client.ID,
			Amount:   200.00,
			Status:   models.InvoiceStatusSent,
		},
	}

	for i := range invoices {
		err = suite.db.Create(&invoices[i]).Error
		suite.NoError(err)
	}

	// Load client with invoices
	var clientWithInvoices models.Client
	err = suite.db.Preload("Invoices").First(&clientWithInvoices, client.ID).Error
	suite.NoError(err)
	suite.Len(clientWithInvoices.Invoices, 2)

	// Test invoice totals
	var totalAmount float64
	for _, invoice := range clientWithInvoices.Invoices {
		totalAmount += invoice.Amount
	}
	suite.Equal(300.00, totalAmount)
}

func (suite *DatabaseIntegrationSuite) TestDatabaseConstraints() {
	// Test unique email constraint
	client1 := &models.Client{
		Name:  "Client 1",
		Email: "duplicate@example.com",
	}
	err := suite.db.Create(client1).Error
	suite.NoError(err)

	client2 := &models.Client{
		Name:  "Client 2",
		Email: "duplicate@example.com", // Same email
	}
	err = suite.db.Create(client2).Error
	suite.Error(err) // Should fail due to unique constraint

	// Test foreign key constraint
	invalidInvoice := &models.Invoice{
		Number:   "INV-999",
		ClientID: 99999, // Non-existent client
		Amount:   100.00,
		Status:   models.InvoiceStatusDraft,
	}
	err = suite.db.Create(invalidInvoice).Error
	suite.Error(err) // Should fail due to foreign key constraint
}

func (suite *DatabaseIntegrationSuite) TestDatabaseTransactions() {
	// Test transaction rollback
	tx := suite.db.Begin()
	
	client := &models.Client{
		Name:  "Transaction Test",
		Email: "transaction@example.com",
	}
	err := tx.Create(client).Error
	suite.NoError(err)
	
	// Rollback transaction
	tx.Rollback()
	
	// Verify client was not saved
	var foundClient models.Client
	err = suite.db.First(&foundClient, "email = ?", "transaction@example.com").Error
	suite.Error(err)
	suite.Equal(gorm.ErrRecordNotFound, err)

	// Test transaction commit
	tx = suite.db.Begin()
	
	err = tx.Create(client).Error
	suite.NoError(err)
	
	// Commit transaction
	err = tx.Commit().Error
	suite.NoError(err)
	
	// Verify client was saved
	err = suite.db.First(&foundClient, "email = ?", "transaction@example.com").Error
	suite.NoError(err)
	suite.Equal(client.Name, foundClient.Name)
}

func (suite *DatabaseIntegrationSuite) TestDatabasePagination() {
	// Create multiple clients
	for i := 0; i < 15; i++ {
		client := &models.Client{
			Name:  fmt.Sprintf("Client %d", i+1),
			Email: fmt.Sprintf("client%d@example.com", i+1),
		}
		err := suite.db.Create(client).Error
		suite.NoError(err)
	}

	// Test pagination
	var clients []models.Client
	
	// Page 1: first 10 clients
	err := suite.db.Limit(10).Offset(0).Find(&clients).Error
	suite.NoError(err)
	suite.Len(clients, 10)

	// Page 2: next 5 clients
	err = suite.db.Limit(10).Offset(10).Find(&clients).Error
	suite.NoError(err)
	suite.Len(clients, 5)

	// Test total count
	var count int64
	err = suite.db.Model(&models.Client{}).Count(&count).Error
	suite.NoError(err)
	suite.Equal(int64(15), count)
}

func (suite *DatabaseIntegrationSuite) createTestDatabase() error {
	// Connect to postgres database to create test database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=%s",
		suite.config.Database.Host,
		suite.config.Database.Port,
		suite.config.Database.Username,
		suite.config.Database.Password,
		suite.config.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Drop test database if exists
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", suite.config.Database.Name))
	if err != nil {
		log.Printf("Warning: Could not drop test database: %v", err)
	}

	// Create test database
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", suite.config.Database.Name))
	return err
}

func TestDatabaseIntegration(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// Check if test database is available
	if os.Getenv("SKIP_DB_TESTS") == "true" {
		t.Skip("Database tests skipped")
	}

	suite.Run(t, new(DatabaseIntegrationSuite))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Benchmark database operations
func BenchmarkClientCreate(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping benchmark in short mode")
	}

	// This would need proper setup similar to the test suite
	b.Skip("Requires database setup")
	
	// Example benchmark structure:
	// b.ResetTimer()
	// for i := 0; i < b.N; i++ {
	//     client := &models.Client{
	//         Name:  fmt.Sprintf("Benchmark Client %d", i),
	//         Email: fmt.Sprintf("bench%d@example.com", i),
	//     }
	//     db.Create(client)
	// }
}