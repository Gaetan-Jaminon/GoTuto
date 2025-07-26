package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"gotuto/api/billing/internal/config"
	"gotuto/api/billing/internal/database"
	"gotuto/api/billing/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// E2ETestSuite tests the complete API functionality end-to-end
type E2ETestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *gorm.DB
	config *config.Config
}

func (suite *E2ETestSuite) SetupSuite() {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Load test configuration
	suite.config = &config.Config{
		Database: config.DatabaseConfig{
			Host:     getEnv("TEST_DB_HOST", "localhost"),
			Port:     5432,
			Username: getEnv("TEST_DB_USER", "postgres"),
			Password: getEnv("TEST_DB_PASSWORD", "password"),
			Name:     getEnv("TEST_DB_NAME", "billing_e2e_test"),
			SSLMode:  "disable",
		},
		Server: config.ServerConfig{
			Port: 8080,
			Mode: "test",
		},
		Pagination: config.PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	// Connect to database
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

	// Setup router (in a real app, this would come from your main setup)
	suite.router = suite.setupTestRouter()
}

func (suite *E2ETestSuite) TearDownSuite() {
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

func (suite *E2ETestSuite) SetupTest() {
	// Clean database before each test
	suite.db.Exec("TRUNCATE TABLE invoices, clients RESTART IDENTITY CASCADE")
}

func (suite *E2ETestSuite) setupTestRouter() *gin.Engine {
	// This is a simplified setup - in real app, import your main router setup
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Note: In a real implementation, you'd import your actual handlers
	// For this example, we'll create simplified versions
	v1 := router.Group("/api/v1")
	{
		v1.GET("/clients", suite.mockGetClients)
		v1.GET("/clients/:id", suite.mockGetClient)
		v1.POST("/clients", suite.mockCreateClient)
		v1.PUT("/clients/:id", suite.mockUpdateClient)
		v1.DELETE("/clients/:id", suite.mockDeleteClient)
		
		v1.GET("/invoices", suite.mockGetInvoices)
		v1.GET("/invoices/:id", suite.mockGetInvoice)
		v1.POST("/invoices", suite.mockCreateInvoice)
		v1.PUT("/invoices/:id", suite.mockUpdateInvoice)
		v1.DELETE("/invoices/:id", suite.mockDeleteInvoice)
	}

	return router
}

func (suite *E2ETestSuite) TestHealthEndpoint() {
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal("healthy", response["status"])
}

func (suite *E2ETestSuite) TestClientLifecycle() {
	// 1. Create a client
	clientData := models.CreateClientRequest{
		Name:    "E2E Test Client",
		Email:   "e2e@example.com",
		Phone:   "+1234567890",
		Address: "123 Test Street",
	}

	jsonData, _ := json.Marshal(clientData)
	req, _ := http.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusCreated, w.Code)
	
	var createdClient models.Client
	err := json.Unmarshal(w.Body.Bytes(), &createdClient)
	suite.NoError(err)
	suite.NotZero(createdClient.ID)
	suite.Equal(clientData.Name, createdClient.Name)
	suite.Equal(clientData.Email, createdClient.Email)

	clientID := createdClient.ID

	// 2. Get the client
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/clients/%d", clientID), nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	
	var fetchedClient models.Client
	err = json.Unmarshal(w.Body.Bytes(), &fetchedClient)
	suite.NoError(err)
	suite.Equal(clientID, fetchedClient.ID)
	suite.Equal(clientData.Name, fetchedClient.Name)

	// 3. Update the client
	updateData := models.UpdateClientRequest{
		Name:  "Updated E2E Client",
		Phone: "+9876543210",
	}

	jsonData, _ = json.Marshal(updateData)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/clients/%d", clientID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	
	var updatedClient models.Client
	err = json.Unmarshal(w.Body.Bytes(), &updatedClient)
	suite.NoError(err)
	suite.Equal(updateData.Name, updatedClient.Name)
	suite.Equal(updateData.Phone, updatedClient.Phone)
	suite.Equal(clientData.Email, updatedClient.Email) // Should remain unchanged

	// 4. List clients (should include our client)
	req, _ = http.NewRequest("GET", "/api/v1/clients", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	
	var listResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &listResponse)
	suite.NoError(err)
	
	clients := listResponse["clients"].([]interface{})
	suite.Len(clients, 1)

	// 5. Delete the client
	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/api/v1/clients/%d", clientID), nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	// 6. Verify client is deleted
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/clients/%d", clientID), nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *E2ETestSuite) TestInvoiceWorkflow() {
	// First create a client
	client := &models.Client{
		Name:  "Invoice Test Client",
		Email: "invoice@example.com",
	}
	err := suite.db.Create(client).Error
	suite.NoError(err)

	// 1. Create an invoice
	invoiceData := models.CreateInvoiceRequest{
		ClientID:    client.ID,
		Amount:      299.99,
		Status:      models.InvoiceStatusDraft,
		IssueDate:   time.Now(),
		DueDate:     time.Now().AddDate(0, 1, 0), // 1 month from now
		Description: "E2E Test Invoice",
	}

	jsonData, _ := json.Marshal(invoiceData)
	req, _ := http.NewRequest("POST", "/api/v1/invoices", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusCreated, w.Code)
	
	var createdInvoice models.Invoice
	err = json.Unmarshal(w.Body.Bytes(), &createdInvoice)
	suite.NoError(err)
	suite.NotZero(createdInvoice.ID)
	suite.Equal(invoiceData.Amount, createdInvoice.Amount)
	suite.Equal(invoiceData.Status, createdInvoice.Status)

	invoiceID := createdInvoice.ID

	// 2. Update invoice status (draft -> sent)
	updateData := models.UpdateInvoiceRequest{
		Status: models.InvoiceStatusSent,
	}

	jsonData, _ = json.Marshal(updateData)
	req, _ = http.NewRequest("PUT", fmt.Sprintf("/api/v1/invoices/%d", invoiceID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	// 3. Get invoice with client details
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/invoices/%d", invoiceID), nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	
	var fetchedInvoice models.Invoice
	err = json.Unmarshal(w.Body.Bytes(), &fetchedInvoice)
	suite.NoError(err)
	suite.Equal(models.InvoiceStatusSent, fetchedInvoice.Status)
	suite.Equal(client.Name, fetchedInvoice.Client.Name)

	// 4. List invoices with filters
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/invoices?client_id=%d&status=%s", client.ID, models.InvoiceStatusSent), nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	
	var listResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &listResponse)
	suite.NoError(err)
	
	invoices := listResponse["invoices"].([]interface{})
	suite.Len(invoices, 1)
}

func (suite *E2ETestSuite) TestPaginationAndSearch() {
	// Create multiple clients
	for i := 0; i < 25; i++ {
		client := &models.Client{
			Name:  fmt.Sprintf("Test Client %d", i+1),
			Email: fmt.Sprintf("test%d@example.com", i+1),
		}
		err := suite.db.Create(client).Error
		suite.NoError(err)
	}

	// Test pagination
	req, _ := http.NewRequest("GET", "/api/v1/clients?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	
	clients := response["clients"].([]interface{})
	suite.Len(clients, 10)
	
	pagination := response["pagination"].(map[string]interface{})
	suite.Equal(float64(1), pagination["page"])
	suite.Equal(float64(10), pagination["limit"])
	suite.Equal(float64(25), pagination["total"])

	// Test search
	req, _ = http.NewRequest("GET", "/api/v1/clients?search=Client 1", nil)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
	
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	
	clients = response["clients"].([]interface{})
	// Should find clients with "Client 1" in name (Client 1, Client 10-19, etc.)
	suite.Greater(len(clients), 0)
}

func (suite *E2ETestSuite) TestErrorHandling() {
	// Test 404 for non-existent client
	req, _ := http.NewRequest("GET", "/api/v1/clients/999999", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusNotFound, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response, "error")

	// Test validation error
	invalidClient := models.CreateClientRequest{
		Name:  "", // Empty name should fail validation
		Email: "invalid-email", // Invalid email format
	}

	jsonData, _ := json.Marshal(invalidClient)
	req, _ = http.NewRequest("POST", "/api/v1/clients", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
	
	err = json.Unmarshal(w.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Contains(response, "error")
}

func TestE2EAPISuite(t *testing.T) {
	// Skip if not in integration test mode
	if testing.Short() {
		t.Skip("Skipping E2E tests in short mode")
	}

	// Check if test database is available
	if os.Getenv("SKIP_E2E_TESTS") == "true" {
		t.Skip("E2E tests skipped")
	}

	suite.Run(t, new(E2ETestSuite))
}

// Mock handlers (in real implementation, these would be your actual handlers)
func (suite *E2ETestSuite) mockGetClients(c *gin.Context) {
	// Simplified implementation for testing
	var clients []models.Client
	query := suite.db

	// Handle search
	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Handle pagination
	page := 1
	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	
	limit := 10
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	
	offset := (page - 1) * limit

	// Get total count
	var total int64
	query.Model(&models.Client{}).Count(&total)

	// Get clients
	query.Limit(limit).Offset(offset).Find(&clients)

	c.JSON(http.StatusOK, gin.H{
		"clients": clients,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (suite *E2ETestSuite) mockGetClient(c *gin.Context) {
	id := c.Param("id")
	var client models.Client
	
	if err := suite.db.Preload("Invoices").First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	
	c.JSON(http.StatusOK, client)
}

func (suite *E2ETestSuite) mockCreateClient(c *gin.Context) {
	var req models.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	client := models.Client{
		Name:    req.Name,
		Email:   req.Email,
		Phone:   req.Phone,
		Address: req.Address,
	}
	
	if err := suite.db.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
		return
	}
	
	c.JSON(http.StatusCreated, client)
}

func (suite *E2ETestSuite) mockUpdateClient(c *gin.Context) {
	id := c.Param("id")
	var client models.Client
	
	if err := suite.db.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	
	var req models.UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if req.Name != "" {
		client.Name = req.Name
	}
	if req.Email != "" {
		client.Email = req.Email
	}
	if req.Phone != "" {
		client.Phone = req.Phone
	}
	if req.Address != "" {
		client.Address = req.Address
	}
	
	if err := suite.db.Save(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client"})
		return
	}
	
	c.JSON(http.StatusOK, client)
}

func (suite *E2ETestSuite) mockDeleteClient(c *gin.Context) {
	id := c.Param("id")
	var client models.Client
	
	if err := suite.db.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	
	if err := suite.db.Delete(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
}

// Simplified invoice handlers (similar pattern)
func (suite *E2ETestSuite) mockGetInvoices(c *gin.Context) {
	var invoices []models.Invoice
	query := suite.db.Preload("Client")
	
	if clientID := c.Query("client_id"); clientID != "" {
		query = query.Where("client_id = ?", clientID)
	}
	
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	
	query.Find(&invoices)
	
	c.JSON(http.StatusOK, gin.H{"invoices": invoices})
}

func (suite *E2ETestSuite) mockGetInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice
	
	if err := suite.db.Preload("Client").First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	
	c.JSON(http.StatusOK, invoice)
}

func (suite *E2ETestSuite) mockCreateInvoice(c *gin.Context) {
	var req models.CreateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	invoice := models.Invoice{
		Number:      fmt.Sprintf("INV-%d", time.Now().Unix()),
		ClientID:    req.ClientID,
		Amount:      req.Amount,
		Status:      req.Status,
		IssueDate:   req.IssueDate,
		DueDate:     req.DueDate,
		Description: req.Description,
	}
	
	if invoice.Status == "" {
		invoice.Status = models.InvoiceStatusDraft
	}
	
	if err := suite.db.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}
	
	suite.db.Preload("Client").First(&invoice, invoice.ID)
	c.JSON(http.StatusCreated, invoice)
}

func (suite *E2ETestSuite) mockUpdateInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice
	
	if err := suite.db.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	
	var req models.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if req.Amount > 0 {
		invoice.Amount = req.Amount
	}
	if req.Status != "" {
		invoice.Status = req.Status
	}
	if !req.IssueDate.IsZero() {
		invoice.IssueDate = req.IssueDate
	}
	if !req.DueDate.IsZero() {
		invoice.DueDate = req.DueDate
	}
	if req.Description != "" {
		invoice.Description = req.Description
	}
	
	if err := suite.db.Save(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice"})
		return
	}
	
	suite.db.Preload("Client").First(&invoice, invoice.ID)
	c.JSON(http.StatusOK, invoice)
}

func (suite *E2ETestSuite) mockDeleteInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice
	
	if err := suite.db.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	
	if err := suite.db.Delete(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invoice"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}