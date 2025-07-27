package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Category represents a product category in the catalog domain
type Category struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	ParentID    *uint          `json:"parent_id"`
	Parent      *Category      `json:"parent,omitempty"`
	Children    []Category     `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Products    []Product      `json:"products,omitempty"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	SortOrder   int            `json:"sort_order" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Validate validates category business rules
func (c *Category) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("category name is required")
	}

	if len(c.Name) > 100 {
		return fmt.Errorf("category name cannot exceed 100 characters")
	}

	if len(c.Description) > 500 {
		return fmt.Errorf("category description cannot exceed 500 characters")
	}

	// Prevent self-reference
	if c.ParentID != nil && *c.ParentID == c.ID {
		return fmt.Errorf("category cannot be its own parent")
	}

	return nil
}

// GetDepth calculates the depth of the category in the hierarchy
func (c *Category) GetDepth() int {
	depth := 0
	current := c
	for current.Parent != nil {
		depth++
		current = current.Parent
	}
	return depth
}

// IsRoot checks if the category is a root category (no parent)
func (c *Category) IsRoot() bool {
	return c.ParentID == nil
}

// HasChildren checks if the category has child categories
func (c *Category) HasChildren() bool {
	return len(c.Children) > 0
}

// GetFullPath returns the full path of the category (e.g., "Electronics > Computers > Laptops")
func (c *Category) GetFullPath() string {
	if c.Parent == nil {
		return c.Name
	}
	return c.Parent.GetFullPath() + " > " + c.Name
}

// CreateCategoryRequest represents the request to create a new category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   *int   `json:"sort_order"`
}

// Validate validates the create category request
func (r *CreateCategoryRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("category name is required")
	}

	if len(r.Name) > 100 {
		return fmt.Errorf("category name cannot exceed 100 characters")
	}

	if len(r.Description) > 500 {
		return fmt.Errorf("category description cannot exceed 500 characters")
	}

	return nil
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
	IsActive    *bool  `json:"is_active"`
	SortOrder   *int   `json:"sort_order"`
}

// Validate validates the update category request
func (r *UpdateCategoryRequest) Validate() error {
	if r.Name != "" && len(r.Name) > 100 {
		return fmt.Errorf("category name cannot exceed 100 characters")
	}

	if len(r.Description) > 500 {
		return fmt.Errorf("category description cannot exceed 500 characters")
	}

	return nil
}

// CategoryWithProductCount represents a category with its product count
type CategoryWithProductCount struct {
	Category
	ProductCount int `json:"product_count"`
}

// MoveCategoryRequest represents the request to move a category to a different parent
type MoveCategoryRequest struct {
	NewParentID *uint `json:"new_parent_id"`
}

// Validate validates the move category request
func (r *MoveCategoryRequest) Validate(categoryID uint) error {
	// Prevent moving to self
	if r.NewParentID != nil && *r.NewParentID == categoryID {
		return fmt.Errorf("category cannot be moved to itself")
	}

	return nil
}