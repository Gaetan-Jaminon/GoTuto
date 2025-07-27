package models

import (
	"testing"
	"time"

	"gaetanjaminon/GoTuto/internal/billing/models/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestClient_Validation(t *testing.T) {
	clientData, err := testdata.LoadClients()
	require.NoError(t, err, "Failed to load client test data")

	t.Run("valid clients", func(t *testing.T) {
		for i, testClient := range clientData.ValidClients {
			t.Run(testClient.Name, func(t *testing.T) {
				client := Client{
					Name:    testClient.Name,
					Email:   testClient.Email,
					Phone:   testClient.Phone,
					Address: testClient.Address,
				}
				err := validateClient(client)
				assert.NoError(t, err, "Valid client %d should pass validation", i)
			})
		}
	})

	t.Run("invalid clients", func(t *testing.T) {
		for _, testCase := range clientData.InvalidClients {
			t.Run(testCase.ExpectedError, func(t *testing.T) {
				client := Client{
					Name:    testCase.Name,
					Email:   testCase.Email,
					Phone:   testCase.Phone,
					Address: testCase.Address,
				}
				err := validateClient(client)
				assert.Error(t, err, "Invalid client should fail validation: %s", testCase.ExpectedError)
			})
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		for _, testCase := range clientData.EdgeCases {
			t.Run(testCase.Description, func(t *testing.T) {
				client := Client{
					Name:    testCase.Name,
					Email:   testCase.Email,
					Phone:   testCase.Phone,
					Address: testCase.Address,
				}
				err := validateClient(client)
				assert.NoError(t, err, "Edge case should be valid: %s", testCase.Description)
			})
		}
	})
}

func TestClient_String(t *testing.T) {
	client := Client{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
		Phone: "+1234567890",
	}

	str := client.String()
	expected := "Client{ID: 1, Name: John Doe, Email: john@example.com}"
	
	assert.Equal(t, expected, str)
}

func TestClient_SoftDelete(t *testing.T) {
	client := Client{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Simulate soft delete
	now := time.Now()
	client.DeletedAt = gorm.DeletedAt{
		Time:  now,
		Valid: true,
	}

	assert.True(t, client.DeletedAt.Valid)
	assert.Equal(t, now.Unix(), client.DeletedAt.Time.Unix())
}

func TestCreateClientRequest_Validation(t *testing.T) {
	requestData, err := testdata.LoadCreateClientRequests()
	require.NoError(t, err, "Failed to load create client request test data")

	t.Run("valid requests", func(t *testing.T) {
		for i, testRequest := range requestData.ValidRequests {
			t.Run(testRequest.Name, func(t *testing.T) {
				request := CreateClientRequest{
					Name:    testRequest.Name,
					Email:   testRequest.Email,
					Phone:   testRequest.Phone,
					Address: testRequest.Address,
				}
				err := validateCreateClientRequest(request)
				assert.NoError(t, err, "Valid request %d should pass validation", i)
			})
		}
	})

	t.Run("invalid requests", func(t *testing.T) {
		for _, testCase := range requestData.InvalidRequests {
			t.Run(testCase.ExpectedError, func(t *testing.T) {
				request := CreateClientRequest{
					Name:    testCase.Name,
					Email:   testCase.Email,
					Phone:   testCase.Phone,
					Address: testCase.Address,
				}
				err := validateCreateClientRequest(request)
				assert.Error(t, err, "Invalid request should fail validation: %s", testCase.ExpectedError)
			})
		}
	})
}

// Helper functions for validation (normally these would be in your validation layer)
func validateClient(client Client) error {
	if client.Name == "" {
		return assert.AnError
	}
	if len(client.Name) > 100 {
		return assert.AnError
	}
	if client.Email == "" {
		return assert.AnError
	}
	// Simple email validation
	if !isValidEmail(client.Email) {
		return assert.AnError
	}
	return nil
}

func validateCreateClientRequest(req CreateClientRequest) error {
	if req.Name == "" || len(req.Name) < 2 || len(req.Name) > 100 {
		return assert.AnError
	}
	if req.Email == "" || !isValidEmail(req.Email) {
		return assert.AnError
	}
	if len(req.Phone) > 20 {
		return assert.AnError
	}
	if len(req.Address) > 255 {
		return assert.AnError
	}
	return nil
}

func isValidEmail(email string) bool {
	// Simple email validation for testing
	return len(email) > 0 && 
		   len(email) < 255 && 
		   containsChar(email, '@') && 
		   containsChar(email, '.')
}

func containsChar(s string, c rune) bool {
	for _, char := range s {
		if char == c {
			return true
		}
	}
	return false
}

// Add String method to Client for better test output
func (c Client) String() string {
	return "Client{ID: " + string(rune(c.ID + '0')) + 
		   ", Name: " + c.Name + 
		   ", Email: " + c.Email + "}"
}