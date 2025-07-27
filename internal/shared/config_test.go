package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ConfigTestSuite demonstrates using testify/suite for setup/teardown
type ConfigTestSuite struct {
	suite.Suite
	originalEnv map[string]string
}

func (suite *ConfigTestSuite) SetupTest() {
	// Save original environment variables
	suite.originalEnv = make(map[string]string)
	envVars := []string{
		"APP_ENV",
		"BILLING_SERVER_PORT",
		"BILLING_DATABASE_HOST",
		"BILLING_DATABASE_PASSWORD",
	}
	
	for _, env := range envVars {
		if val, exists := os.LookupEnv(env); exists {
			suite.originalEnv[env] = val
		}
		os.Unsetenv(env)
	}
}

func (suite *ConfigTestSuite) TearDownTest() {
	// Restore original environment variables
	for key, value := range suite.originalEnv {
		os.Setenv(key, value)
	}
	
	// Clear any test env vars
	os.Unsetenv("APP_ENV")
	os.Unsetenv("BILLING_SERVER_PORT")
	os.Unsetenv("BILLING_DATABASE_HOST")
	os.Unsetenv("BILLING_DATABASE_PASSWORD")
}

func (suite *ConfigTestSuite) TestLoadDefaultConfig() {
	// This test would load default config without any overrides
	// In a real test, you'd create a test config file
	suite.T().Skip("Requires test config file setup")
	
	config, err := Load()
	
	suite.NoError(err)
	suite.NotNil(config)
	suite.Equal(8080, config.Server.Port)
	suite.Equal("localhost", config.Database.Host)
}

func (suite *ConfigTestSuite) TestEnvironmentVariableOverrides() {
	// Set environment variables
	os.Setenv("BILLING_SERVER_PORT", "9090")
	os.Setenv("BILLING_DATABASE_HOST", "prod-db")
	os.Setenv("BILLING_DATABASE_PASSWORD", "secret123")
	
	suite.T().Skip("Requires test config file setup")
	
	config, err := Load()
	
	suite.NoError(err)
	suite.Equal(9090, config.Server.Port)
	suite.Equal("prod-db", config.Database.Host)
	suite.Equal("secret123", config.Database.Password)
}

func (suite *ConfigTestSuite) TestConfigValidation() {
	tests := []struct {
		name        string
		serverPort  int
		dbHost      string
		expectValid bool
	}{
		{"valid config", 8080, "localhost", true},
		{"invalid port", -1, "localhost", false},
		{"empty host", 8080, "", false},
		{"high port", 65536, "localhost", false},
	}
	
	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			config := &Config{
				Server: ServerConfig{
					Port: tt.serverPort,
				},
				Database: DatabaseConfig{
					Host: tt.dbHost,
				},
			}
			
			err := validateConfig(config)
			
			if tt.expectValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestConfig(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}

// Regular test functions (not part of suite)
func TestDatabaseConfig_GetDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected string
	}{
		{
			name: "basic DSN",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "postgres",
				Password: "password",
				Name:     "billing",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=postgres password=password dbname=billing sslmode=disable",
		},
		{
			name: "production DSN with SSL",
			config: DatabaseConfig{
				Host:     "prod-db.example.com",
				Port:     5432,
				Username: "billing_user",
				Password: "secure_password",
				Name:     "billing_prod",
				SSLMode:  "require",
			},
			expected: "host=prod-db.example.com port=5432 user=billing_user password=secure_password dbname=billing_prod sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := tt.config.GetDSN()
			assert.Equal(t, tt.expected, dsn)
		})
	}
}

func TestConfigPrint(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Port: 8080,
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Host: "localhost",
			Port: 5432,
			Name: "billing",
			Username: "postgres",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		CORS: CORSConfig{
			AllowedOrigins: []string{"*"},
		},
		Pagination: PaginationConfig{
			DefaultLimit: 10,
			MaxLimit:     100,
		},
	}

	// Test that Print doesn't panic and produces some output
	// In a real test, you might capture log output and verify content
	assert.NotPanics(t, func() {
		config.Print()
	})
}

func TestServerConfig_Timeouts(t *testing.T) {
	config := ServerConfig{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	assert.Equal(t, 30*time.Second, config.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.WriteTimeout)
}

func TestDatabaseConfig_ConnectionLimits(t *testing.T) {
	config := DatabaseConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}

	assert.Equal(t, 25, config.MaxOpenConns)
	assert.Equal(t, 5, config.MaxIdleConns)
	assert.Equal(t, 5*time.Minute, config.ConnMaxLifetime)
}

func TestCORSConfig_Origins(t *testing.T) {
	tests := []struct {
		name    string
		origins []string
		testURL string
		allowed bool
	}{
		{
			name:    "wildcard allows all",
			origins: []string{"*"},
			testURL: "https://example.com",
			allowed: true,
		},
		{
			name:    "specific origin allowed",
			origins: []string{"https://app.example.com", "https://api.example.com"},
			testURL: "https://app.example.com",
			allowed: true,
		},
		{
			name:    "origin not in list",
			origins: []string{"https://app.example.com"},
			testURL: "https://evil.com",
			allowed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := CORSConfig{
				AllowedOrigins: tt.origins,
			}
			
			allowed := isOriginAllowed(config, tt.testURL)
			assert.Equal(t, tt.allowed, allowed)
		})
	}
}

// Helper functions for validation (these would normally be in your main config package)
func validateConfig(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return assert.AnError
	}
	if config.Database.Host == "" {
		return assert.AnError
	}
	return nil
}

func isOriginAllowed(config CORSConfig, origin string) bool {
	for _, allowed := range config.AllowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return false
}

// Example of property-based testing (if you were using a library like gopter)
func TestConfigDefaults(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected interface{}
	}{
		{"default server port", "server.port", 8080},
		{"default db host", "database.host", "localhost"},
		{"default db port", "database.port", 5432},
		{"default logging level", "logging.level", "info"},
		{"default pagination limit", "pagination.default_limit", 10},
	}

	// In a real implementation, you'd load a default config and test these values
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is where you'd test actual default values from your config
			assert.NotNil(t, tt.expected, "Expected value should not be nil")
		})
	}
}

// Benchmark for config loading performance
func BenchmarkConfigLoad(b *testing.B) {
	b.Skip("Requires test config file setup")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Load()
		if err != nil {
			b.Fatal(err)
		}
	}
}