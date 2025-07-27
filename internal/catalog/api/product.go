package api

import (
	"net/http"
	"strconv"
	
	"gaetanjaminon/GoTuto/internal/catalog/database"
	"gaetanjaminon/GoTuto/internal/catalog/models"
	
	"github.com/gin-gonic/gin"
)

// GetProducts retrieves all products with optional pagination and filters
func GetProducts(c *gin.Context) {
	var products []models.Product
	
	// Optional pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit
	
	// Optional filters
	search := c.Query("search")
	categoryID := c.Query("category_id")
	isActive := c.Query("is_active")
	
	query := database.DB.Preload("Category").Limit(limit).Offset(offset)
	
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	
	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}
	
	// Get total count for pagination
	var total int64
	countQuery := database.DB.Model(&models.Product{})
	if search != "" {
		countQuery = countQuery.Where("name ILIKE ? OR description ILIKE ? OR sku ILIKE ?", 
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if categoryID != "" {
		countQuery = countQuery.Where("category_id = ?", categoryID)
	}
	if isActive != "" {
		countQuery = countQuery.Where("is_active = ?", isActive == "true")
	}
	countQuery.Count(&total)
	
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetProduct retrieves a single product by ID
func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	
	if err := database.DB.Preload("Category").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	
	c.JSON(http.StatusOK, product)
}

// CreateProduct creates a new product
func CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Check if SKU already exists
	var existingProduct models.Product
	if err := database.DB.Where("sku = ?", req.SKU).First(&existingProduct).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product with this SKU already exists"})
		return
	}
	
	// Verify category exists if provided
	if req.CategoryID != nil {
		var category models.Category
		if err := database.DB.First(&category, *req.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
	}
	
	product := models.Product{
		SKU:         req.SKU,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Currency:    req.Currency,
		CategoryID:  req.CategoryID,
	}
	
	// Set default currency if not provided
	if product.Currency == "" {
		product.Currency = "USD"
	}
	
	// Set is_active if provided
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	
	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}
	
	// Load category data for response
	database.DB.Preload("Category").First(&product, product.ID)
	
	c.JSON(http.StatusCreated, product)
}

// UpdateProduct updates an existing product
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	
	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Verify category exists if provided
	if req.CategoryID != nil {
		var category models.Category
		if err := database.DB.First(&category, *req.CategoryID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found"})
			return
		}
	}
	
	// Update only provided fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Currency != "" {
		product.Currency = req.Currency
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	
	if err := database.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	
	// Load category data for response
	database.DB.Preload("Category").First(&product, product.ID)
	
	c.JSON(http.StatusOK, product)
}

// DeleteProduct soft deletes a product
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	
	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}