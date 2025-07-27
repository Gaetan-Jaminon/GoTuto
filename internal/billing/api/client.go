package api

import (
	"net/http"
	"strconv"
	
	"gaetanjaminon/GoTuto/internal/billing/models"
	
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetClients retrieves all clients with optional pagination
func GetClients(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var clients []models.Client
		
		// Optional pagination
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset := (page - 1) * limit
		
		// Optional search by name or email
		search := c.Query("search")
		query := db.Limit(limit).Offset(offset)
		
		if search != "" {
			query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
		}
		
		if err := query.Find(&clients).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve clients"})
			return
		}
		
		// Get total count for pagination
		var total int64
		countQuery := db.Model(&models.Client{})
		if search != "" {
			countQuery = countQuery.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
		}
		countQuery.Count(&total)
		
		c.JSON(http.StatusOK, gin.H{
			"clients": clients,
			"pagination": gin.H{
				"page":  page,
				"limit": limit,
				"total": total,
			},
		})
	}
}

// GetClient retrieves a single client by ID
func GetClient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var client models.Client
		
		// Include invoices in the response
		if err := db.Preload("Invoices").First(&client, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		
		c.JSON(http.StatusOK, client)
	}
}

// CreateClient creates a new client
func CreateClient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
		
		if err := db.Create(&client).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
			return
		}
		
		c.JSON(http.StatusCreated, client)
	}
}

// UpdateClient updates an existing client
func UpdateClient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var client models.Client
		
		if err := db.First(&client, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		
		var req models.UpdateClientRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Update only provided fields
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
		
		if err := db.Save(&client).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update client"})
			return
		}
		
		c.JSON(http.StatusOK, client)
	}
}

// DeleteClient soft deletes a client
func DeleteClient(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var client models.Client
		
		if err := db.First(&client, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Client not found"})
			return
		}
		
		// Check if client has invoices
		var invoiceCount int64
		db.Model(&models.Invoice{}).Where("client_id = ?", id).Count(&invoiceCount)
		
		if invoiceCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Cannot delete client with existing invoices",
				"invoice_count": invoiceCount,
			})
			return
		}
		
		if err := db.Delete(&client).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete client"})
			return
		}
		
		c.JSON(http.StatusOK, gin.H{"message": "Client deleted successfully"})
	}
}