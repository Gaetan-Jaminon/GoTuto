package infrastructure

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// LoadDomainConfig loads configuration for a specific domain with environment overrides
func LoadDomainConfig[T any](domainName string, envPrefix string) (*T, error) {
	// Get environment (default to "dev")
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// Create new viper instance for isolated config loading
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath("./config/base")
	v.AddConfigPath(fmt.Sprintf("./config/%s", domainName))

	// 1. Load base defaults
	v.SetConfigName("base")
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read base config: %w", err)
	}

	// 2. Load base environment overrides (optional)
	v.SetConfigName(env)
	v.SetConfigFile(fmt.Sprintf("./config/base/%s.yaml", env))
	if err := v.MergeInConfig(); err != nil {
		// Ignore error - environment override is optional
	}

	// 3. Load domain defaults
	v.SetConfigName(domainName)
	v.SetConfigFile(fmt.Sprintf("./config/%s/%s.yaml", domainName, domainName))
	if err := v.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read %s domain config: %w", domainName, err)
	}

	// 4. Load domain environment overrides (optional)
	v.SetConfigName(env)
	v.SetConfigFile(fmt.Sprintf("./config/%s/%s.yaml", domainName, env))
	if err := v.MergeInConfig(); err != nil {
		// Ignore error - environment override is optional
	}

	// 5. Environment variables with domain-specific prefix
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 6. Unmarshal to struct
	var config T
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}