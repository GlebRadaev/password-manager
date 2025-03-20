package app

import (
	"fmt"
	"time"
)

type Config struct {
	HttpPort string   `mapstructure:"httpPort" validate:"required"`
	GrpcPort string   `mapstructure:"grpcPort" validate:"required"`
	PgConfig PgConfig `mapstructure:"pgConfig" validate:"required"`
}

type PgConfig struct {
	User            string        `mapstructure:"user" envconfig:"PG_USER" validate:"required"`
	Password        string        `mapstructure:"password" envconfig:"PG_PASSWORD" validate:"required"`
	DbName          string        `mapstructure:"dbName" validate:"required"`
	MaxOpenConns    int32         `mapstructure:"maxOpenConns" validate:"required"`
	MaxIdleConns    int32         `mapstructure:"maxIdleConns" validate:"required"`
	MinConns        int32         `mapstructure:"minConns" validate:"required"`
	MaxConnLifetime time.Duration `mapstructure:"maxConnLifetime" validate:"required"`
	Timeout         time.Duration `mapstructure:"timeout" validate:"required"`
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
}

type GRPCClient struct {
	Endpoint string `mapstructure:"endpoint" envconfig:"ENDPOINT" validate:"required"`
}

func (c PgConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.DbName)
}
