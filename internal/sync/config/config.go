package config

import (
	"fmt"
	"os"

	"github.com/GlebRadaev/password-manager/internal/common/app"
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type LocalConfig struct {
	DataSvc *app.GRPCClient `mapstructure:"dataSvc" validate:"required"`
}
type Config struct {
	Env         string `mapstructure:"env" envconfig:"ENVIRONMENT" validate:"required"`
	LogLevel    string `mapstructure:"log_level" validate:"required"`
	app.Config  `mapstructure:"sync"`
	LocalConfig `mapstructure:"sync"`
}

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
	fmt.Printf("%v\n", cfg)
	return &cfg, nil
}
