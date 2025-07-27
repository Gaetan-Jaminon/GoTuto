package testdata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Client represents a client for testing (matches models.Client)
type Client struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Invoice represents an invoice for testing (matches models.Invoice)
type Invoice struct {
	ID          uint      `json:"id"`
	Number      string    `json:"number"`
	ClientID    uint      `json:"client_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	IssueDate   time.Time `json:"issue_date"`
	DueDate     time.Time `json:"due_date"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateClientRequest represents a create client request for testing
type CreateClientRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// CreateInvoiceRequest represents a create invoice request for testing
type CreateInvoiceRequest struct {
	ClientID    uint      `json:"client_id"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	IssueDate   time.Time `json:"issue_date"`
	DueDate     time.Time `json:"due_date"`
	Description string    `json:"description"`
}

// ClientTestData represents the structure of clients.json
type ClientTestData struct {
	ValidClients   []Client `json:"valid_clients"`
	InvalidClients []struct {
		Client
		ExpectedError string `json:"expected_error"`
	} `json:"invalid_clients"`
	EdgeCases []struct {
		Client
		Description string `json:"description"`
	} `json:"edge_cases"`
}

// InvoiceTestData represents the structure of invoices.json
type InvoiceTestData struct {
	ValidInvoices   []Invoice `json:"valid_invoices"`
	InvalidInvoices []struct {
		Invoice
		ExpectedError string `json:"expected_error"`
	} `json:"invalid_invoices"`
	StatusScenarios []struct {
		Name        string `json:"name"`
		FromStatus  string `json:"from_status"`
		ToStatus    string `json:"to_status"`
		ShouldAllow bool   `json:"should_allow"`
	} `json:"status_scenarios"`
	OverdueScenarios []struct {
		Name               string `json:"name"`
		Status             string `json:"status"`
		DueDateOffsetDays  int    `json:"due_date_offset_days"`
		ExpectedOverdue    bool   `json:"expected_overdue"`
	} `json:"overdue_scenarios"`
}

// CreateClientRequestTestData represents the structure of requests/create_client.json
type CreateClientRequestTestData struct {
	ValidRequests   []CreateClientRequest `json:"valid_requests"`
	InvalidRequests []struct {
		CreateClientRequest
		ExpectedError string `json:"expected_error"`
	} `json:"invalid_requests"`
}

// CreateInvoiceRequestTestData represents the structure of requests/create_invoice.json
type CreateInvoiceRequestTestData struct {
	ValidRequests   []CreateInvoiceRequest `json:"valid_requests"`
	InvalidRequests []struct {
		CreateInvoiceRequest
		ExpectedError string `json:"expected_error"`
	} `json:"invalid_requests"`
}

// LoadClients loads client test data from clients.json
func LoadClients() (*ClientTestData, error) {
	data, err := loadJSONFile("clients.json")
	if err != nil {
		return nil, err
	}

	var clientData ClientTestData
	if err := json.Unmarshal(data, &clientData); err != nil {
		return nil, err
	}

	return &clientData, nil
}

// LoadInvoices loads invoice test data from invoices.json
func LoadInvoices() (*InvoiceTestData, error) {
	data, err := loadJSONFile("invoices.json")
	if err != nil {
		return nil, err
	}

	var invoiceData InvoiceTestData
	if err := json.Unmarshal(data, &invoiceData); err != nil {
		return nil, err
	}

	// Convert date strings to time.Time and adjust relative dates
	now := time.Now()
	for i := range invoiceData.ValidInvoices {
		// Parse issue and due dates
		if invoiceData.ValidInvoices[i].IssueDate.IsZero() {
			invoiceData.ValidInvoices[i].IssueDate = now
		}
		if invoiceData.ValidInvoices[i].DueDate.IsZero() {
			invoiceData.ValidInvoices[i].DueDate = now.AddDate(0, 1, 0)
		}
	}

	for i := range invoiceData.InvalidInvoices {
		if invoiceData.InvalidInvoices[i].IssueDate.IsZero() {
			invoiceData.InvalidInvoices[i].IssueDate = now
		}
		if invoiceData.InvalidInvoices[i].DueDate.IsZero() {
			invoiceData.InvalidInvoices[i].DueDate = now.AddDate(0, 1, 0)
		}
	}

	return &invoiceData, nil
}

// LoadCreateClientRequests loads create client request test data
func LoadCreateClientRequests() (*CreateClientRequestTestData, error) {
	data, err := loadJSONFile("requests/create_client.json")
	if err != nil {
		return nil, err
	}

	var requestData CreateClientRequestTestData
	if err := json.Unmarshal(data, &requestData); err != nil {
		return nil, err
	}

	return &requestData, nil
}

// LoadCreateInvoiceRequests loads create invoice request test data
func LoadCreateInvoiceRequests() (*CreateInvoiceRequestTestData, error) {
	data, err := loadJSONFile("requests/create_invoice.json")
	if err != nil {
		return nil, err
	}

	var requestData CreateInvoiceRequestTestData
	if err := json.Unmarshal(data, &requestData); err != nil {
		return nil, err
	}

	// Convert date strings to time.Time for valid requests
	now := time.Now()
	for i := range requestData.ValidRequests {
		if requestData.ValidRequests[i].IssueDate.IsZero() {
			requestData.ValidRequests[i].IssueDate = now
		}
		if requestData.ValidRequests[i].DueDate.IsZero() {
			requestData.ValidRequests[i].DueDate = now.AddDate(0, 1, 0)
		}
	}

	for i := range requestData.InvalidRequests {
		if requestData.InvalidRequests[i].IssueDate.IsZero() {
			requestData.InvalidRequests[i].IssueDate = now
		}
		if requestData.InvalidRequests[i].DueDate.IsZero() {
			requestData.InvalidRequests[i].DueDate = now.AddDate(0, 1, 0)
		}
	}

	return &requestData, nil
}

// GenerateOverdueInvoice creates an invoice for overdue testing based on scenario
func GenerateOverdueInvoice(scenario struct {
	Name               string `json:"name"`
	Status             string `json:"status"`
	DueDateOffsetDays  int    `json:"due_date_offset_days"`
	ExpectedOverdue    bool   `json:"expected_overdue"`
}) Invoice {
	now := time.Now()
	return Invoice{
		ID:        1,
		Number:    "INV-TEST",
		ClientID:  1,
		Amount:    100.00,
		Status:    scenario.Status,
		IssueDate: now.AddDate(0, 0, -30), // 30 days ago
		DueDate:   now.AddDate(0, 0, scenario.DueDateOffsetDays),
	}
}

// loadJSONFile loads a JSON file from the testdata directory
func loadJSONFile(filename string) ([]byte, error) {
	// Use runtime.Caller to get the current file's directory
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to get current file path")
	}
	
	// Get the directory of this loader.go file (which is the testdata directory)
	testdataDir := filepath.Dir(currentFile)
	filePath := filepath.Join(testdataDir, filename)
	
	return os.ReadFile(filePath)
}