package models

import (
	"time"
	"gorm.io/gorm"
)

type Client struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Phone     string         `json:"phone"`
	Address   string         `json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationship
	Invoices []Invoice `json:"invoices,omitempty" gorm:"foreignKey:ClientID"`
}

type CreateClientRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone" binding:"max=20"`
	Address string `json:"address" binding:"max=255"`
}

type UpdateClientRequest struct {
	Name    string `json:"name" binding:"omitempty,min=2,max=100"`
	Email   string `json:"email" binding:"omitempty,email"`
	Phone   string `json:"phone" binding:"omitempty,max=20"`
	Address string `json:"address" binding:"omitempty,max=255"`
}