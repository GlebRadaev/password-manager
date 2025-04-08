// Package config handles application configuration loading and validation.
// Supports YAML files and environment variables with required field validation.
package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"

	"github.com/GlebRadaev/password-manager/internal/common/app"
)

// LocalConfig contains service-specific configurations
type LocalConfig struct {
	DataSvc *app.GRPCClient `mapstructure:"dataSvc" validate:"required"` // gRPC client configuration for data service
}

// Config represents the complete application configuration
type Config struct {
	Env         string                `mapstructure:"env" envconfig:"ENVIRONMENT" validate:"required"` // Application environment (dev, prod, etc.)
	LogLevel    string                `mapstructure:"log_level" validate:"required"`                   // Logging level
	app.Config  `mapstructure:"sync"` // Embedded common app config
	LocalConfig `mapstructure:"sync"` // Embedded local service config
}

// New loads and validates application configuration from:
// 1. YAML config file (based on ENVIRONMENT variable)
// 2. Environment variables (overriding file values)
// 3. Validates all required fields
//
// Returns:
// - *Config: loaded configuration
// - error: if loading or validation fails
func New() (*Config, error) {
	environment, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		environment = "dev" // Default to dev environment
	}

	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(environment + ".config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed read env file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed unmarshal config: %w", err)
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed unmarshal env: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("missing required attributes %v", err)
	}

	return &cfg, nil
}
