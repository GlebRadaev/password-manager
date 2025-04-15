// Package config provides configuration management for the application.
// It handles loading and validating configuration from both YAML files and environment variables.
package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"

	"github.com/GlebRadaev/password-manager/internal/common/app"
)

// GatewayConfig contains configuration for the API gateway service,
// including ports and gRPC client configurations for dependent services.
type GatewayConfig struct {
	HTTPPort string          `mapstructure:"httpPort" validate:"required"`
	GrpcPort string          `mapstructure:"grpcPort" validate:"required"`
	AuthSvc  *app.GRPCClient `mapstructure:"authSvc" validate:"required"`
	DataSvc  *app.GRPCClient `mapstructure:"dataSvc" validate:"required"`
	SyncSvc  *app.GRPCClient `mapstructure:"syncSvc" validate:"required"`
}

// Config represents the overall application configuration,
// including environment settings, logging, and gateway configuration.
type Config struct {
	Env           string `mapstructure:"env" envconfig:"ENVIRONMENT" validate:"required"`
	LogLevel      string `mapstructure:"log_level" validate:"required"`
	GatewayConfig `mapstructure:"gateway"`
}

// New creates and validates a new Config instance
func New() (*Config, error) {
	environment, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		environment = "dev"
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
