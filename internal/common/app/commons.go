// Package app contains shared application configuration types
package app

import (
	"fmt"
	"time"
)

// Config holds main application configuration
type Config struct {
	HTTPPort string   `mapstructure:"httpPort" validate:"required"` // HTTP server port
	GrpcPort string   `mapstructure:"grpcPort" validate:"required"` // gRPC server port
	PgConfig PgConfig `mapstructure:"pgConfig" validate:"required"` // Postgres config
}

// PgConfig contains PostgreSQL connection settings
type PgConfig struct {
	User            string        `mapstructure:"user" envconfig:"PG_USER" validate:"required"`
	Password        string        `mapstructure:"password" envconfig:"PG_PASSWORD" validate:"required"`
	DbName          string        `mapstructure:"dbName" validate:"required"`
	MaxOpenConns    int32         `mapstructure:"maxOpenConns" validate:"required"` // Max open connections
	MaxIdleConns    int32         `mapstructure:"maxIdleConns" validate:"required"` // Max idle connections
	MinConns        int32         `mapstructure:"minConns" validate:"required"`     // Min connections
	MaxConnLifetime time.Duration `mapstructure:"maxConnLifetime" validate:"required"`
	Timeout         time.Duration `mapstructure:"timeout" validate:"required"`
	Host            string        `mapstructure:"host"` // DB host (default: localhost)
	Port            int           `mapstructure:"port"` // DB port (default: 5432)
}

// GRPCClient contains gRPC client endpoint config
type GRPCClient struct {
	Endpoint string `mapstructure:"endpoint" envconfig:"ENDPOINT" validate:"required"`
}

// DSN generates PostgreSQL connection string
func (c PgConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.DbName)
}
