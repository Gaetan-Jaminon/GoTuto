package models

import (
	"time"
	"gorm.io/gorm"
)

type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusSent      InvoiceStatus = "sent"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusOverdue   InvoiceStatus = "overdue"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

type Invoice struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Number      string         `json:"number" gorm:"uniqueIndex;not null"`
	ClientID    uint           `json:"client_id" gorm:"not null"`
	Amount      float64        `json:"amount" gorm:"not null"`
	Status      InvoiceStatus  `json:"status" gorm:"default:'draft'"`
	IssueDate   time.Time      `json:"issue_date"`
	DueDate     time.Time      `json:"due_date"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationship
	Client Client `json:"client,omitempty" gorm:"foreignKey:ClientID"`
}

type CreateInvoiceRequest struct {
	ClientID    uint          `json:"client_id" binding:"required"`
	Amount      float64       `json:"amount" binding:"required,gt=0"`
	Status      InvoiceStatus `json:"status" binding:"omitempty,oneof=draft sent paid overdue cancelled"`
	IssueDate   time.Time     `json:"issue_date" binding:"required"`
	DueDate     time.Time     `json:"due_date" binding:"required"`
	Description string        `json:"description" binding:"max=500"`
}

type UpdateInvoiceRequest struct {
	Amount      float64       `json:"amount" binding:"omitempty,gt=0"`
	Status      InvoiceStatus `json:"status" binding:"omitempty,oneof=draft sent paid overdue cancelled"`
	IssueDate   time.Time     `json:"issue_date" binding:"omitempty"`
	DueDate     time.Time     `json:"due_date" binding:"omitempty"`
	Description string        `json:"description" binding:"omitempty,max=500"`
}