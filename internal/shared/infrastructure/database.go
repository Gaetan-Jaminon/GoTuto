package infrastructure

import (
	"fmt"
	"time"
)

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host              string        `mapstructure:"host"`
	Port              int           `mapstructure:"port"`
	Username          string        `mapstructure:"username"`
	Password          string        `mapstructure:"password"`
	Name              string        `mapstructure:"name"`
	Schema            string        `mapstructure:"schema"`
	SSLMode           string        `mapstructure:"ssl_mode"`
	MaxOpenConns      int           `mapstructure:"max_open_conns"`
	MaxIdleConns      int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime   time.Duration `mapstructure:"conn_max_lifetime"`
	ConnectionTimeout time.Duration `mapstructure:"connection_timeout"`
}

// GetDSN returns PostgreSQL connection string with schema support
func (c *DatabaseConfig) GetDSN() string {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Name, c.SSLMode)
	
	// Add search_path if schema is specified
	if c.Schema != "" {
		dsn += fmt.Sprintf(" search_path=%s", c.Schema)
	}
	
	return dsn
}

// MigrationConfig holds migration-specific database configuration
type MigrationConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Schema   string `mapstructure:"schema"`
}