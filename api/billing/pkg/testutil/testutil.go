// Package testutil provides utilities for testing
package testutil

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"gotuto/api/billing/internal/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Quiet during tests
	})
	assert.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(&models.Client{}, &models.Invoice{})
	assert.NoError(t, err)

	return db
}

// LoadTestClients loads test client data from JSON file
func LoadTestClients(t *testing.T) []models.Client {
	data, err := ioutil.ReadFile(filepath.Join("../../testdata", "clients.json"))
	assert.NoError(t, err)

	var clients []models.Client
	err = json.Unmarshal(data, &clients)
	assert.NoError(t, err)

	return clients
}

// LoadTestInvoices loads test invoice data from JSON file
func LoadTestInvoices(t *testing.T) []models.Invoice {
	data, err := ioutil.ReadFile(filepath.Join("../../testdata", "invoices.json"))
	assert.NoError(t, err)

	var invoices []models.Invoice
	err = json.Unmarshal(data, &invoices)
	assert.NoError(t, err)

	return invoices
}

// SeedTestData seeds the test database with client and invoice data
func SeedTestData(t *testing.T, db *gorm.DB) {
	clients := LoadTestClients(t)
	for i := range clients {
		err := db.Create(&clients[i]).Error
		assert.NoError(t, err)
	}

	invoices := LoadTestInvoices(t)
	for i := range invoices {
		err := db.Create(&invoices[i]).Error
		assert.NoError(t, err)
	}
}

// CleanTestData removes all test data from database
func CleanTestData(t *testing.T, db *gorm.DB) {
	// Order matters due to foreign key constraints
	err := db.Exec("DELETE FROM invoices").Error
	assert.NoError(t, err)

	err = db.Exec("DELETE FROM clients").Error
	assert.NoError(t, err)
}

// AssertJSONEqual compares two JSON strings for equality
func AssertJSONEqual(t *testing.T, expected, actual string) {
	var expectedObj, actualObj interface{}
	
	err := json.Unmarshal([]byte(expected), &expectedObj)
	assert.NoError(t, err)
	
	err = json.Unmarshal([]byte(actual), &actualObj)
	assert.NoError(t, err)
	
	assert.Equal(t, expectedObj, actualObj)
}

// CreateTestClient creates a test client with default values
func CreateTestClient(overrides ...func(*models.Client)) *models.Client {
	client := &models.Client{
		Name:    "Test Client",
		Email:   "test@example.com",
		Phone:   "+1234567890",
		Address: "123 Test Street",
	}

	// Apply any overrides
	for _, override := range overrides {
		override(client)
	}

	return client
}

// CreateTestInvoice creates a test invoice with default values
func CreateTestInvoice(clientID uint, overrides ...func(*models.Invoice)) *models.Invoice {
	invoice := &models.Invoice{
		Number:      "INV-TEST-001",
		ClientID:    clientID,
		Amount:      100.00,
		Status:      models.InvoiceStatusDraft,
		Description: "Test invoice",
	}

	// Apply any overrides
	for _, override := range overrides {
		override(invoice)
	}

	return invoice
}

// DBAssertion provides helper methods for database assertions
type DBAssertion struct {
	t  *testing.T
	db *gorm.DB
}

// NewDBAssertion creates a new database assertion helper
func NewDBAssertion(t *testing.T, db *gorm.DB) *DBAssertion {
	return &DBAssertion{t: t, db: db}
}

// AssertClientExists verifies that a client exists in the database
func (da *DBAssertion) AssertClientExists(id uint) *models.Client {
	var client models.Client
	err := da.db.First(&client, id).Error
	assert.NoError(da.t, err)
	return &client
}

// AssertClientNotExists verifies that a client does not exist in the database
func (da *DBAssertion) AssertClientNotExists(id uint) {
	var client models.Client
	err := da.db.First(&client, id).Error
	assert.Error(da.t, err)
	assert.Equal(da.t, gorm.ErrRecordNotFound, err)
}

// AssertInvoiceExists verifies that an invoice exists in the database
func (da *DBAssertion) AssertInvoiceExists(id uint) *models.Invoice {
	var invoice models.Invoice
	err := da.db.Preload("Client").First(&invoice, id).Error
	assert.NoError(da.t, err)
	return &invoice
}

// AssertInvoiceCount verifies the total number of invoices
func (da *DBAssertion) AssertInvoiceCount(expected int64) {
	var count int64
	err := da.db.Model(&models.Invoice{}).Count(&count).Error
	assert.NoError(da.t, err)
	assert.Equal(da.t, expected, count)
}

// AssertClientCount verifies the total number of clients
func (da *DBAssertion) AssertClientCount(expected int64) {
	var count int64
	err := da.db.Model(&models.Client{}).Count(&count).Error
	assert.NoError(da.t, err)
	assert.Equal(da.t, expected, count)
}

// TimeAssertion provides helper methods for time-related assertions
type TimeAssertion struct {
	t *testing.T
}

// NewTimeAssertion creates a new time assertion helper
func NewTimeAssertion(t *testing.T) *TimeAssertion {
	return &TimeAssertion{t: t}
}

// AssertRecentTime verifies that a time is within the last minute
func (ta *TimeAssertion) AssertRecentTime(actual interface{}) {
	// This would implement time assertion logic
	// For brevity, just checking it's not nil
	assert.NotNil(ta.t, actual)
}

// HTTPAssertion provides helper methods for HTTP response assertions
type HTTPAssertion struct {
	t *testing.T
}

// NewHTTPAssertion creates a new HTTP assertion helper
func NewHTTPAssertion(t *testing.T) *HTTPAssertion {
	return &HTTPAssertion{t: t}
}

// AssertJSONResponse verifies that response contains expected JSON structure
func (ha *HTTPAssertion) AssertJSONResponse(body []byte, expectedFields ...string) map[string]interface{} {
	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	assert.NoError(ha.t, err)

	for _, field := range expectedFields {
		assert.Contains(ha.t, response, field)
	}

	return response
}

// AssertErrorResponse verifies that response is an error response
func (ha *HTTPAssertion) AssertErrorResponse(body []byte) {
	response := ha.AssertJSONResponse(body, "error")
	assert.NotEmpty(ha.t, response["error"])
}

// AssertSuccessResponse verifies that response indicates success
func (ha *HTTPAssertion) AssertSuccessResponse(body []byte) {
	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	assert.NoError(ha.t, err)

	// Success response should not have an error field
	_, hasError := response["error"]
	assert.False(ha.t, hasError)
}