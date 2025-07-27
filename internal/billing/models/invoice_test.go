package models

import (
	"testing"

	"gaetanjaminon/GoTuto/internal/billing/models/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	invoiceData, err := testdata.LoadInvoices()
	require.NoError(t, err, "Failed to load invoice test data")

	t.Run("valid invoices", func(t *testing.T) {
		for i, testInvoice := range invoiceData.ValidInvoices {
			t.Run(testInvoice.Number, func(t *testing.T) {
				invoice := Invoice{
					Number:      testInvoice.Number,
					ClientID:    testInvoice.ClientID,
					Amount:      testInvoice.Amount,
					Status:      InvoiceStatus(testInvoice.Status),
					IssueDate:   testInvoice.IssueDate,
					DueDate:     testInvoice.DueDate,
					Description: testInvoice.Description,
				}
				err := validateInvoice(invoice)
				assert.NoError(t, err, "Valid invoice %d should pass validation", i)
			})
		}
	})

	t.Run("invalid invoices", func(t *testing.T) {
		for _, testCase := range invoiceData.InvalidInvoices {
			t.Run(testCase.ExpectedError, func(t *testing.T) {
				invoice := Invoice{
					Number:      testCase.Number,
					ClientID:    testCase.ClientID,
					Amount:      testCase.Amount,
					Status:      InvoiceStatus(testCase.Status),
					IssueDate:   testCase.IssueDate,
					DueDate:     testCase.DueDate,
					Description: testCase.Description,
				}
				err := validateInvoice(invoice)
				assert.Error(t, err, "Invalid invoice should fail validation: %s", testCase.ExpectedError)
			})
		}
	})
}

func TestInvoice_StatusTransitions(t *testing.T) {
	invoiceData, err := testdata.LoadInvoices()
	require.NoError(t, err, "Failed to load invoice test data")

	for _, scenario := range invoiceData.StatusScenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			fromStatus := InvoiceStatus(scenario.FromStatus)
			toStatus := InvoiceStatus(scenario.ToStatus)
			allowed := isValidStatusTransition(fromStatus, toStatus)
			assert.Equal(t, scenario.ShouldAllow, allowed, 
				"Status transition %s -> %s should be %v", scenario.FromStatus, scenario.ToStatus, scenario.ShouldAllow)
		})
	}
}

func TestInvoice_IsOverdue(t *testing.T) {
	invoiceData, err := testdata.LoadInvoices()
	require.NoError(t, err, "Failed to load invoice test data")

	for _, scenario := range invoiceData.OverdueScenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			testInvoice := testdata.GenerateOverdueInvoice(scenario)
			invoice := Invoice{
				ID:        testInvoice.ID,
				Number:    testInvoice.Number,
				ClientID:  testInvoice.ClientID,
				Amount:    testInvoice.Amount,
				Status:    InvoiceStatus(testInvoice.Status),
				IssueDate: testInvoice.IssueDate,
				DueDate:   testInvoice.DueDate,
			}
			result := invoice.IsOverdue()
			assert.Equal(t, scenario.ExpectedOverdue, result, 
				"Overdue check for scenario '%s' should be %v", scenario.Name, scenario.ExpectedOverdue)
		})
	}
}

func TestCreateInvoiceRequest_Validation(t *testing.T) {
	requestData, err := testdata.LoadCreateInvoiceRequests()
	require.NoError(t, err, "Failed to load create invoice request test data")

	t.Run("valid requests", func(t *testing.T) {
		for i, testRequest := range requestData.ValidRequests {
			t.Run(testRequest.Description, func(t *testing.T) {
				request := CreateInvoiceRequest{
					ClientID:    testRequest.ClientID,
					Amount:      testRequest.Amount,
					Status:      InvoiceStatus(testRequest.Status),
					IssueDate:   testRequest.IssueDate,
					DueDate:     testRequest.DueDate,
					Description: testRequest.Description,
				}
				err := validateCreateInvoiceRequest(request)
				assert.NoError(t, err, "Valid request %d should pass validation", i)
			})
		}
	})

	t.Run("invalid requests", func(t *testing.T) {
		for _, testCase := range requestData.InvalidRequests {
			t.Run(testCase.ExpectedError, func(t *testing.T) {
				request := CreateInvoiceRequest{
					ClientID:    testCase.ClientID,
					Amount:      testCase.Amount,
					Status:      InvoiceStatus(testCase.Status),
					IssueDate:   testCase.IssueDate,
					DueDate:     testCase.DueDate,
					Description: testCase.Description,
				}
				err := validateCreateInvoiceRequest(request)
				assert.Error(t, err, "Invalid request should fail validation: %s", testCase.ExpectedError)
			})
		}
	})
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