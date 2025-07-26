package config

import (
	"fmt"
	"log"
	"strings"
	"time"
	
	"github.com/spf13/viper"
)

// Config holds all configuration for our application
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Logging    LoggingConfig    `mapstructure:"logging"`
	CORS       CORSConfig       `mapstructure:"cors"`
	Pagination PaginationConfig `mapstructure:"pagination"`
}

type ServerConfig struct {
	Port         int           `mapstructure:"port"`
	Mode         string        `mapstructure:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

type PaginationConfig struct {
	DefaultLimit int `mapstructure:"default_limit"`
	MaxLimit     int `mapstructure:"max_limit"`
}

// Global config instance
var Cfg *Config

// Load reads configuration from files and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	
	// Read base configuration
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config: %w", err)
	}
	
	// Get environment from APP_ENV or default to "dev"
	env := viper.GetString("APP_ENV")
	if env == "" {
		env = "dev"
	}
	
	// Load environment-specific config
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	if err := viper.MergeInConfig(); err != nil {
		log.Printf("No environment-specific config found for %s: %v", env, err)
	}
	
	// Enable environment variable overrides
	// Example: BILLING_DATABASE_HOST will override database.host
	viper.SetEnvPrefix("BILLING")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	
	// Unmarshal config into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	// Store globally for easy access
	Cfg = &config
	
	return &config, nil
}

// GetDSN returns PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Name, c.SSLMode)
}

// Print logs the current configuration (with sensitive data masked)
func (c *Config) Print() {
	log.Println("=== Configuration ===")
	log.Printf("Server: Port=%d, Mode=%s", c.Server.Port, c.Server.Mode)
	log.Printf("Database: Host=%s:%d, Name=%s, User=%s", 
		c.Database.Host, c.Database.Port, c.Database.Name, c.Database.Username)
	log.Printf("Logging: Level=%s, Format=%s", c.Logging.Level, c.Logging.Format)
	log.Printf("CORS Origins: %v", c.CORS.AllowedOrigins)
	log.Printf("Pagination: Default=%d, Max=%d", c.Pagination.DefaultLimit, c.Pagination.MaxLimit)
}

// MustLoad loads config and panics if it fails
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	return cfg
}