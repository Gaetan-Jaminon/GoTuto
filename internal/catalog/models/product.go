package models

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Product represents a product in the catalog domain
type Product struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	SKU         string         `json:"sku" gorm:"uniqueIndex;not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	Currency    string         `json:"currency" gorm:"default:'USD'"`
	CategoryID  *uint          `json:"category_id"`
	Category    *Category      `json:"category,omitempty"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusActive       ProductStatus = "active"
	ProductStatusInactive     ProductStatus = "inactive"
	ProductStatusDiscontinued ProductStatus = "discontinued"
)

// Validate validates product business rules
func (p *Product) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("product name is required")
	}

	if len(p.Name) > 200 {
		return fmt.Errorf("product name cannot exceed 200 characters")
	}

	if strings.TrimSpace(p.SKU) == "" {
		return fmt.Errorf("product SKU is required")
	}

	if len(p.SKU) > 50 {
		return fmt.Errorf("product SKU cannot exceed 50 characters")
	}

	if p.Price < 0 {
		return fmt.Errorf("product price cannot be negative")
	}

	if len(p.Description) > 1000 {
		return fmt.Errorf("product description cannot exceed 1000 characters")
	}

	// Validate currency code (basic validation)
	if p.Currency != "" && len(p.Currency) != 3 {
		return fmt.Errorf("currency must be a 3-letter code")
	}

	return nil
}

// IsValidStatus checks if the product status is valid
func IsValidProductStatus(status ProductStatus) bool {
	switch status {
	case ProductStatusActive, ProductStatusInactive, ProductStatusDiscontinued:
		return true
	default:
		return false
	}
}

// FormatPrice formats the price with currency
func (p *Product) FormatPrice() string {
	if p.Currency == "" {
		return fmt.Sprintf("%.2f", p.Price)
	}
	return fmt.Sprintf("%.2f %s", p.Price, p.Currency)
}

// CreateProductRequest represents the request to create a new product
type CreateProductRequest struct {
	SKU         string  `json:"sku" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0"`
	Currency    string  `json:"currency"`
	CategoryID  *uint   `json:"category_id"`
	IsActive    *bool   `json:"is_active"`
}

// Validate validates the create product request
func (r *CreateProductRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("product name is required")
	}

	if strings.TrimSpace(r.SKU) == "" {
		return fmt.Errorf("product SKU is required")
	}

	if r.Price < 0 {
		return fmt.Errorf("product price cannot be negative")
	}

	if len(r.Name) > 200 {
		return fmt.Errorf("product name cannot exceed 200 characters")
	}

	if len(r.SKU) > 50 {
		return fmt.Errorf("product SKU cannot exceed 50 characters")
	}

	if len(r.Description) > 1000 {
		return fmt.Errorf("product description cannot exceed 1000 characters")
	}

	return nil
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       *float64 `json:"price"`
	Currency    string   `json:"currency"`
	CategoryID  *uint    `json:"category_id"`
	IsActive    *bool    `json:"is_active"`
}

// Validate validates the update product request
func (r *UpdateProductRequest) Validate() error {
	if r.Name != "" && len(r.Name) > 200 {
		return fmt.Errorf("product name cannot exceed 200 characters")
	}

	if len(r.Description) > 1000 {
		return fmt.Errorf("product description cannot exceed 1000 characters")
	}

	if r.Price != nil && *r.Price < 0 {
		return fmt.Errorf("product price cannot be negative")
	}

	if r.Currency != "" && len(r.Currency) != 3 {
		return fmt.Errorf("currency must be a 3-letter code")
	}

	return nil
}