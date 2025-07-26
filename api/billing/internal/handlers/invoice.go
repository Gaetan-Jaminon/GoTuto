package handlers

import (
	"net/http"
	"strconv"
	"time"
	
	"gotuto/api/billing/internal/database"
	"gotuto/api/billing/internal/models"
	
	"github.com/gin-gonic/gin"
)

// GetInvoices retrieves all invoices with optional filters
func GetInvoices(c *gin.Context) {
	var invoices []models.Invoice
	
	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	
	// Filters
	clientID := c.Query("client_id")
	status := c.Query("status")
	
	query := database.DB.Preload("Client").Limit(limit).Offset(offset)
	
	if clientID != "" {
		query = query.Where("client_id = ?", clientID)
	}
	
	if status != "" {
		query = query.Where("status = ?", status)
	}
	
	if err := query.Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoices"})
		return
	}
	
	// Get total count
	var total int64
	countQuery := database.DB.Model(&models.Invoice{})
	if clientID != "" {
		countQuery = countQuery.Where("client_id = ?", clientID)
	}
	if status != "" {
		countQuery = countQuery.Where("status = ?", status)
	}
	countQuery.Count(&total)
	
	c.JSON(http.StatusOK, gin.H{
		"invoices": invoices,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetInvoice retrieves a single invoice by ID
func GetInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice
	
	if err := database.DB.Preload("Client").First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	
	c.JSON(http.StatusOK, invoice)
}

// CreateInvoice creates a new invoice
func CreateInvoice(c *gin.Context) {
	var req models.CreateInvoiceRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Verify client exists
	var client models.Client
	if err := database.DB.First(&client, req.ClientID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Client not found"})
		return
	}
	
	// Generate invoice number (simple format: INV-YYYYMMDD-XXXX)
	var count int64
	database.DB.Model(&models.Invoice{}).Where("DATE(created_at) = ?", time.Now().Format("2006-01-02")).Count(&count)
	invoiceNumber := time.Now().Format("INV-20060102-") + strconv.FormatInt(count+1, 10)
	
	invoice := models.Invoice{
		Number:      invoiceNumber,
		ClientID:    req.ClientID,
		Amount:      req.Amount,
		Status:      req.Status,
		IssueDate:   req.IssueDate,
		DueDate:     req.DueDate,
		Description: req.Description,
	}
	
	// Set default status if not provided
	if invoice.Status == "" {
		invoice.Status = models.InvoiceStatusDraft
	}
	
	if err := database.DB.Create(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}
	
	// Load client data for response
	database.DB.Preload("Client").First(&invoice, invoice.ID)
	
	c.JSON(http.StatusCreated, invoice)
}

// UpdateInvoice updates an existing invoice
func UpdateInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice
	
	if err := database.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	
	var req models.UpdateInvoiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Update only provided fields
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
	
	if err := database.DB.Save(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice"})
		return
	}
	
	// Load client data for response
	database.DB.Preload("Client").First(&invoice, invoice.ID)
	
	c.JSON(http.StatusOK, invoice)
}

// DeleteInvoice soft deletes an invoice
func DeleteInvoice(c *gin.Context) {
	id := c.Param("id")
	var invoice models.Invoice
	
	if err := database.DB.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	
	// Prevent deletion of paid invoices
	if invoice.Status == models.InvoiceStatusPaid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete paid invoice"})
		return
	}
	
	if err := database.DB.Delete(&invoice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invoice"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}

// GetInvoicesByClient retrieves all invoices for a specific client
func GetInvoicesByClient(c *gin.Context) {
	clientID := c.Param("client_id")
	
	// Verify client exists
	var client models.Client
	if err := database.DB.First(&client, clientID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
		return
	}
	
	var invoices []models.Invoice
	if err := database.DB.Where("client_id = ?", clientID).Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoices"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"client":   client,
		"invoices": invoices,
	})
}