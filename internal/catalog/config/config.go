package config

import (
	"gaetanjaminon/GoTuto/internal/shared/infrastructure"
)

// CatalogConfig holds all configuration for the catalog domain
type CatalogConfig struct {
	Server    infrastructure.ServerConfig    `mapstructure:"server"`
	Database  infrastructure.DatabaseConfig  `mapstructure:"database"`
	Migration infrastructure.MigrationConfig `mapstructure:"migration"`
	Logging   infrastructure.LoggingConfig   `mapstructure:"logging"`
	CORS      infrastructure.CORSConfig      `mapstructure:"cors"`
	
	// Catalog-specific configuration
	Pagination PaginationConfig `mapstructure:"pagination"`
	Product    ProductConfig    `mapstructure:"product"`
	Category   CategoryConfig   `mapstructure:"category"`
}

// PaginationConfig holds pagination settings for catalog domain
type PaginationConfig struct {
	DefaultLimit int `mapstructure:"default_limit"`
	MaxLimit     int `mapstructure:"max_limit"`
}

// ProductConfig holds product-specific settings
type ProductConfig struct {
	SKUPrefix       string `mapstructure:"sku_prefix"`
	DefaultCurrency string `mapstructure:"default_currency"`
	AllowZeroPrice  bool   `mapstructure:"allow_zero_price"`
}

// CategoryConfig holds category-specific settings
type CategoryConfig struct {
	MaxDepth         int  `mapstructure:"max_depth"`
	AllowCircularRefs bool `mapstructure:"allow_circular_refs"`
}

// Load reads catalog configuration from files and environment
func Load() (*CatalogConfig, error) {
	return infrastructure.LoadDomainConfig[CatalogConfig]("catalog", "CATALOG")
}

// MustLoad loads config and panics if it fails
func MustLoad() *CatalogConfig {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}