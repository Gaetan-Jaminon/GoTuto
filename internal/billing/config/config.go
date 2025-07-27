package config

import (
	"fmt"
	"gaetanjaminon/GoTuto/internal/shared/infrastructure"
)

// BillingConfig holds all configuration for the billing domain
type BillingConfig struct {
	Server    infrastructure.ServerConfig    `mapstructure:"server"`
	Database  infrastructure.DatabaseConfig  `mapstructure:"database"`
	Migration infrastructure.MigrationConfig `mapstructure:"migration"`
	Logging   infrastructure.LoggingConfig   `mapstructure:"logging"`
	CORS      infrastructure.CORSConfig      `mapstructure:"cors"`
	
	// Billing-specific configuration
	Pagination PaginationConfig `mapstructure:"pagination"`
	Invoice    InvoiceConfig    `mapstructure:"invoice"`
	Client     ClientConfig     `mapstructure:"client"`
}

// PaginationConfig holds pagination settings for billing domain
type PaginationConfig struct {
	DefaultLimit int `mapstructure:"default_limit"`
	MaxLimit     int `mapstructure:"max_limit"`
}

// InvoiceConfig holds invoice-specific settings
type InvoiceConfig struct {
	NumberPrefix      string `mapstructure:"number_prefix"`
	DefaultCurrency   string `mapstructure:"default_currency"`
	PaymentTermsDays  int    `mapstructure:"payment_terms_days"`
}

// ClientConfig holds client-specific settings
type ClientConfig struct {
	RequireEmailVerification bool `mapstructure:"require_email_verification"`
	MaxNameLength           int  `mapstructure:"max_name_length"`
}

// Validate checks if the configuration is valid
func (c *BillingConfig) Validate() error {
	// Server validation
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d (must be 1-65535)", c.Server.Port)
	}

	// Database validation
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d (must be 1-65535)", c.Database.Port)
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if c.Database.Username == "" {
		return fmt.Errorf("database username is required")
	}

	// Pagination validation
	if c.Pagination.DefaultLimit <= 0 {
		return fmt.Errorf("pagination default limit must be positive")
	}
	if c.Pagination.MaxLimit <= 0 {
		return fmt.Errorf("pagination max limit must be positive")
	}
	if c.Pagination.DefaultLimit > c.Pagination.MaxLimit {
		return fmt.Errorf("pagination default limit (%d) cannot exceed max limit (%d)", c.Pagination.DefaultLimit, c.Pagination.MaxLimit)
	}

	// Invoice validation
	if c.Invoice.PaymentTermsDays < 0 {
		return fmt.Errorf("payment terms days cannot be negative")
	}

	// Client validation
	if c.Client.MaxNameLength <= 0 {
		return fmt.Errorf("client max name length must be positive")
	}

	return nil
}

// Load reads billing configuration from files and environment
func Load() (*BillingConfig, error) {
	cfg, err := infrastructure.LoadDomainConfig[BillingConfig]("billing", "BILLING")
	if err != nil {
		return nil, err
	}
	
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	
	return cfg, nil
}

// MustLoad loads config and panics if it fails
func MustLoad() *BillingConfig {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}