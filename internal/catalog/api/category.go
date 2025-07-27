package api

import (
	"net/http"
	"strconv"
	
	"gaetanjaminon/GoTuto/internal/catalog/database"
	"gaetanjaminon/GoTuto/internal/catalog/models"
	
	"github.com/gin-gonic/gin"
)

// GetCategories retrieves all categories with optional pagination and filters
func GetCategories(c *gin.Context) {
	var categories []models.Category
	
	// Optional pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit
	
	// Optional filters
	search := c.Query("search")
	parentID := c.Query("parent_id")
	isActive := c.Query("is_active")
	rootOnly := c.Query("root_only") == "true"
	
	query := database.DB.Preload("Parent").Preload("Children").Limit(limit).Offset(offset).Order("sort_order ASC, name ASC")
	
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	
	if parentID != "" {
		query = query.Where("parent_id = ?", parentID)
	} else if rootOnly {
		query = query.Where("parent_id IS NULL")
	}
	
	if isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	
	if err := query.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve categories"})
		return
	}
	
	// Get total count for pagination
	var total int64
	countQuery := database.DB.Model(&models.Category{})
	if search != "" {
		countQuery = countQuery.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if parentID != "" {
		countQuery = countQuery.Where("parent_id = ?", parentID)
	} else if rootOnly {
		countQuery = countQuery.Where("parent_id IS NULL")
	}
	if isActive != "" {
		countQuery = countQuery.Where("is_active = ?", isActive == "true")
	}
	countQuery.Count(&total)
	
	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetCategory retrieves a single category by ID
func GetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	
	if err := database.DB.Preload("Parent").Preload("Children").Preload("Products").First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	
	c.JSON(http.StatusOK, category)
}

// CreateCategory creates a new category
func CreateCategory(c *gin.Context) {
	var req models.CreateCategoryRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Verify parent category exists if provided
	if req.ParentID != nil {
		var parentCategory models.Category
		if err := database.DB.First(&parentCategory, *req.ParentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent category not found"})
			return
		}
	}
	
	category := models.Category{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}
	
	// Set optional fields if provided
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	
	// Load parent data for response
	database.DB.Preload("Parent").First(&category, category.ID)
	
	c.JSON(http.StatusCreated, category)
}

// UpdateCategory updates an existing category
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	
	var req models.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate request
	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Prevent self-reference
	if req.ParentID != nil && *req.ParentID == category.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category cannot be its own parent"})
		return
	}
	
	// Verify parent category exists if provided
	if req.ParentID != nil {
		var parentCategory models.Category
		if err := database.DB.First(&parentCategory, *req.ParentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent category not found"})
			return
		}
	}
	
	// Update only provided fields
	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Description != "" {
		category.Description = req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	
	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}
	
	// Load parent data for response
	database.DB.Preload("Parent").Preload("Children").First(&category, category.ID)
	
	c.JSON(http.StatusOK, category)
}

// DeleteCategory soft deletes a category
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	
	// Check if category has child categories
	var childCount int64
	database.DB.Model(&models.Category{}).Where("parent_id = ?", id).Count(&childCount)
	
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete category with child categories",
			"child_count": childCount,
		})
		return
	}
	
	// Check if category has products
	var productCount int64
	database.DB.Model(&models.Product{}).Where("category_id = ?", id).Count(&productCount)
	
	if productCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete category with existing products",
			"product_count": productCount,
		})
		return
	}
	
	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// MoveCategory moves a category to a different parent
func MoveCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	
	var req models.MoveCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate request
	categoryID, _ := strconv.ParseUint(id, 10, 32)
	if err := req.Validate(uint(categoryID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Verify new parent category exists if provided
	if req.NewParentID != nil {
		var parentCategory models.Category
		if err := database.DB.First(&parentCategory, *req.NewParentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "New parent category not found"})
			return
		}
	}
	
	// Update parent
	category.ParentID = req.NewParentID
	
	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move category"})
		return
	}
	
	// Load updated data for response
	database.DB.Preload("Parent").Preload("Children").First(&category, category.ID)
	
	c.JSON(http.StatusOK, category)
}

// GetCategoryProducts retrieves all products for a specific category
func GetCategoryProducts(c *gin.Context) {
	categoryID := c.Param("category_id")
	
	// Verify category exists
	var category models.Category
	if err := database.DB.First(&category, categoryID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}
	
	// Optional pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit
	
	var products []models.Product
	query := database.DB.Where("category_id = ?", categoryID).Limit(limit).Offset(offset)
	
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}
	
	// Get total count
	var total int64
	database.DB.Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&total)
	
	c.JSON(http.StatusOK, gin.H{
		"category": category,
		"products": products,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}