// Package application provides configuration management for the gateway service.
package application

import "github.com/GlebRadaev/password-manager/internal/gateway/config"

// ConfigProvider defines the interface for retrieving gateway service configuration.
type ConfigProvider interface {
	// New creates and returns a new gateway configuration instance.
	New() (*config.Config, error)
}

// defaultConfigProvider is the default implementation of ConfigProvider
// that loads configuration from the environment variables.
type defaultConfigProvider struct{}

// New loads and returns the gateway configuration from environment variables.
func (d *defaultConfigProvider) New() (*config.Config, error) {
	return config.New()
}
