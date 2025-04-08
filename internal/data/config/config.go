// Package config handles app configuration from YAML and env vars
package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"

	"github.com/GlebRadaev/password-manager/internal/common/app"
)

// Config contains app settings
type Config struct {
	Env        string `mapstructure:"env" envconfig:"ENVIRONMENT" validate:"required"`
	LogLevel   string `mapstructure:"log_level" validate:"required"`
	app.Config `mapstructure:"data"`
}

// New loads config from YAML + env vars and validates it
// ENVIRONMENT var determines which YAML file to load (dev.config.yaml etc)
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
