package config

import (
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

// Load reads billing configuration from files and environment
func Load() (*BillingConfig, error) {
	return infrastructure.LoadDomainConfig[BillingConfig]("billing", "BILLING")
}

// MustLoad loads config and panics if it fails
func MustLoad() *BillingConfig {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}