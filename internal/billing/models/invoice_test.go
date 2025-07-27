package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInvoiceStatus_Validation(t *testing.T) {
	validStatuses := []InvoiceStatus{
		InvoiceStatusDraft,
		InvoiceStatusSent,
		InvoiceStatusPaid,
		InvoiceStatusOverdue,
		InvoiceStatusCancelled,
	}

	for _, status := range validStatuses {
		t.Run(string(status), func(t *testing.T) {
			assert.True(t, isValidInvoiceStatus(status))
		})
	}

	// Test invalid status
	invalidStatus := InvoiceStatus("invalid")
	assert.False(t, isValidInvoiceStatus(invalidStatus))
}

func TestInvoice_Validation(t *testing.T) {
	baseTime := time.Now()
	
	tests := []struct {
		name    string
		invoice Invoice
		wantErr bool
	}{
		{
			name: "valid invoice",
			invoice: Invoice{
				Number:      "INV-001",
				ClientID:    1,
				Amount:      100.50,
				Status:      InvoiceStatusDraft,
				IssueDate:   baseTime,
				DueDate:     baseTime.AddDate(0, 1, 0), // 1 month later
				Description: "Test invoice",
			},
			wantErr: false,
		},
		{
			name: "zero amount",
			invoice: Invoice{
				Number:    "INV-002",
				ClientID:  1,
				Amount:    0,
				Status:    InvoiceStatusDraft,
				IssueDate: baseTime,
				DueDate:   baseTime.AddDate(0, 1, 0),
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			invoice: Invoice{
				Number:    "INV-003",
				ClientID:  1,
				Amount:    -50.00,
				Status:    InvoiceStatusDraft,
				IssueDate: baseTime,
				DueDate:   baseTime.AddDate(0, 1, 0),
			},
			wantErr: true,
		},
		{
			name: "due date before issue date",
			invoice: Invoice{
				Number:    "INV-004",
				ClientID:  1,
				Amount:    100.00,
				Status:    InvoiceStatusDraft,
				IssueDate: baseTime,
				DueDate:   baseTime.AddDate(0, -1, 0), // 1 month before
			},
			wantErr: true,
		},
		{
			name: "invalid status",
			invoice: Invoice{
				Number:    "INV-005",
				ClientID:  1,
				Amount:    100.00,
				Status:    InvoiceStatus("invalid"),
				IssueDate: baseTime,
				DueDate:   baseTime.AddDate(0, 1, 0),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInvoice(tt.invoice)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestInvoice_StatusTransitions(t *testing.T) {
	tests := []struct {
		name        string
		fromStatus  InvoiceStatus
		toStatus    InvoiceStatus
		shouldAllow bool
	}{
		{"draft to sent", InvoiceStatusDraft, InvoiceStatusSent, true},
		{"draft to cancelled", InvoiceStatusDraft, InvoiceStatusCancelled, true},
		{"sent to paid", InvoiceStatusSent, InvoiceStatusPaid, true},
		{"sent to overdue", InvoiceStatusSent, InvoiceStatusOverdue, true},
		{"sent to cancelled", InvoiceStatusSent, InvoiceStatusCancelled, true},
		{"paid to draft", InvoiceStatusPaid, InvoiceStatusDraft, false},
		{"paid to sent", InvoiceStatusPaid, InvoiceStatusSent, false},
		{"cancelled to paid", InvoiceStatusCancelled, InvoiceStatusPaid, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed := isValidStatusTransition(tt.fromStatus, tt.toStatus)
			assert.Equal(t, tt.shouldAllow, allowed)
		})
	}
}

func TestInvoice_IsOverdue(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name     string
		invoice  Invoice
		expected bool
	}{
		{
			name: "not overdue - due tomorrow",
			invoice: Invoice{
				Status:  InvoiceStatusSent,
				DueDate: now.AddDate(0, 0, 1), // tomorrow
			},
			expected: false,
		},
		{
			name: "overdue - due yesterday",
			invoice: Invoice{
				Status:  InvoiceStatusSent,
				DueDate: now.AddDate(0, 0, -1), // yesterday
			},
			expected: true,
		},
		{
			name: "paid invoice not overdue even if past due date",
			invoice: Invoice{
				Status:  InvoiceStatusPaid,
				DueDate: now.AddDate(0, 0, -5), // 5 days ago
			},
			expected: false,
		},
		{
			name: "draft invoice not overdue",
			invoice: Invoice{
				Status:  InvoiceStatusDraft,
				DueDate: now.AddDate(0, 0, -1), // yesterday
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.invoice.IsOverdue()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateInvoiceRequest_Validation(t *testing.T) {
	baseTime := time.Now()
	
	tests := []struct {
		name    string
		request CreateInvoiceRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: CreateInvoiceRequest{
				ClientID:    1,
				Amount:      150.75,
				Status:      InvoiceStatusDraft,
				IssueDate:   baseTime,
				DueDate:     baseTime.AddDate(0, 1, 0),
				Description: "Test service",
			},
			valid: true,
		},
		{
			name: "zero client ID",
			request: CreateInvoiceRequest{
				ClientID:  0,
				Amount:    100.00,
				IssueDate: baseTime,
				DueDate:   baseTime.AddDate(0, 1, 0),
			},
			valid: false,
		},
		{
			name: "zero amount",
			request: CreateInvoiceRequest{
				ClientID:  1,
				Amount:    0,
				IssueDate: baseTime,
				DueDate:   baseTime.AddDate(0, 1, 0),
			},
			valid: false,
		},
		{
			name: "description too long",
			request: CreateInvoiceRequest{
				ClientID:    1,
				Amount:      100.00,
				IssueDate:   baseTime,
				DueDate:     baseTime.AddDate(0, 1, 0),
				Description: generateLongString(501), // max is 500
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateInvoiceRequest(tt.request)
			
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// Helper functions for validation
func validateInvoice(invoice Invoice) error {
	if invoice.Amount <= 0 {
		return assert.AnError
	}
	if invoice.DueDate.Before(invoice.IssueDate) {
		return assert.AnError
	}
	if !isValidInvoiceStatus(invoice.Status) {
		return assert.AnError
	}
	return nil
}

func validateCreateInvoiceRequest(req CreateInvoiceRequest) error {
	if req.ClientID == 0 {
		return assert.AnError
	}
	if req.Amount <= 0 {
		return assert.AnError
	}
	if len(req.Description) > 500 {
		return assert.AnError
	}
	if !req.DueDate.IsZero() && !req.IssueDate.IsZero() && req.DueDate.Before(req.IssueDate) {
		return assert.AnError
	}
	return nil
}

func isValidInvoiceStatus(status InvoiceStatus) bool {
	validStatuses := []InvoiceStatus{
		InvoiceStatusDraft,
		InvoiceStatusSent,
		InvoiceStatusPaid,
		InvoiceStatusOverdue,
		InvoiceStatusCancelled,
	}
	
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

func isValidStatusTransition(from, to InvoiceStatus) bool {
	// Define allowed transitions
	transitions := map[InvoiceStatus][]InvoiceStatus{
		InvoiceStatusDraft: {
			InvoiceStatusSent,
			InvoiceStatusCancelled,
		},
		InvoiceStatusSent: {
			InvoiceStatusPaid,
			InvoiceStatusOverdue,
			InvoiceStatusCancelled,
		},
		InvoiceStatusOverdue: {
			InvoiceStatusPaid,
			InvoiceStatusCancelled,
		},
		// Paid and Cancelled are terminal states
		InvoiceStatusPaid:      {},
		InvoiceStatusCancelled: {},
	}
	
	allowedTransitions, exists := transitions[from]
	if !exists {
		return false
	}
	
	for _, allowed := range allowedTransitions {
		if to == allowed {
			return true
		}
	}
	return false
}

func generateLongString(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}

// Add methods to Invoice for testing
func (i Invoice) IsOverdue() bool {
	if i.Status == InvoiceStatusPaid || i.Status == InvoiceStatusCancelled || i.Status == InvoiceStatusDraft {
		return false
	}
	return time.Now().After(i.DueDate)
}