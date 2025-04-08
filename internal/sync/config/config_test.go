package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GlebRadaev/password-manager/internal/common/app"
)

func TestNewConfig(t *testing.T) {
	originalEnv := os.Getenv("ENVIRONMENT")
	defer func() {
		os.Setenv("ENVIRONMENT", originalEnv)
	}()

	tests := []struct {
		name        string
		envValue    string
		configFile  string
		want        *Config
		wantErr     bool
		errContains string
		setup       func()
		cleanup     func()
	}{
		{
			name:     "successful config load",
			envValue: "dev",
			configFile: `env: dev
log_level: debug
sync:
  httpPort: :8080
  grpcPort: :9090
  dataSvc:
    endpoint: "localhost:50051"
  pgConfig:
    user: postgres
    password: postgres
    dbName: auth
    maxOpenConns: 10
    maxIdleConns: 5
    minConns: 1
    maxConnLifetime: 30s
    timeout: 5s
    host: localhost
    port: 5555
`,
			want: &Config{
				Env:      "dev",
				LogLevel: "debug",
				Config: app.Config{
					HTTPPort: ":8080",
					GrpcPort: ":9090",
					PgConfig: app.PgConfig{
						Host:            "localhost",
						Port:            5555,
						User:            "postgres",
						Password:        "postgres",
						DbName:          "auth",
						MaxOpenConns:    10,
						MaxIdleConns:    5,
						MinConns:        1,
						MaxConnLifetime: 30 * time.Second,
						Timeout:         5 * time.Second,
					},
				},
				LocalConfig: LocalConfig{
					DataSvc: &app.GRPCClient{
						Endpoint: "localhost:50051",
					},
				},
			},
			wantErr: false,
			setup: func() {
				os.Setenv("ENVIRONMENT", "dev")
				err := os.WriteFile("dev.config.yaml", []byte(`env: dev
log_level: debug
sync:
  httpPort: :8080
  grpcPort: :9090
  dataSvc:
    endpoint: "localhost:50051"
  pgConfig:
    user: postgres
    password: postgres
    dbName: auth
    maxOpenConns: 10
    maxIdleConns: 5
    minConns: 1
    maxConnLifetime: 30s
    timeout: 5s
    host: localhost
    port: 5555
`), 0644)
				if err != nil {
					t.Fatal(err)
				}
			},
			cleanup: func() {
				os.Remove("dev.config.yaml")
			},
		},
		{
			name:        "missing config file",
			envValue:    "missing",
			wantErr:     true,
			errContains: "failed read env file",
			setup: func() {
				os.Setenv("ENVIRONMENT", "missing")
			},
		},
		{
			name:     "invalid yaml syntax",
			envValue: "invalid",
			configFile: `env: dev
log_level debug  # missing colon
`,
			wantErr:     true,
			errContains: "could not find expected ':'",
			setup: func() {
				os.Setenv("ENVIRONMENT", "invalid")
				err := os.WriteFile("invalid.config.yaml", []byte(`env: dev
log_level debug
`), 0644)
				if err != nil {
					t.Fatal(err)
				}
			},
			cleanup: func() {
				os.Remove("invalid.config.yaml")
			},
		},
		{
			name:     "validation error - missing required fields",
			envValue: "validation",
			configFile: `env: dev
log_level: debug
sync:
  httpPort: :8080
`,
			wantErr:     true,
			errContains: "missing required attributes",
			setup: func() {
				os.Setenv("ENVIRONMENT", "validation")
				err := os.WriteFile("validation.config.yaml", []byte(`env: dev
log_level: debug
sync:
  httpPort: :8080
`), 0644)
				if err != nil {
					t.Fatal(err)
				}
			},
			cleanup: func() {
				os.Remove("validation.config.yaml")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}
			if tt.cleanup != nil {
				defer tt.cleanup()
			}

			got, err := New()

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want.Env, got.Env)
				assert.Equal(t, tt.want.LogLevel, got.LogLevel)

				// Проверка Config
				assert.Equal(t, tt.want.HTTPPort, got.HTTPPort)
				assert.Equal(t, tt.want.GrpcPort, got.GrpcPort)

				// Проверка PgConfig
				if tt.want.PgConfig != (app.PgConfig{}) {
					assert.Equal(t, tt.want.PgConfig.Host, got.PgConfig.Host)
					assert.Equal(t, tt.want.PgConfig.Port, got.PgConfig.Port)
					assert.Equal(t, tt.want.PgConfig.User, got.PgConfig.User)
					assert.Equal(t, tt.want.PgConfig.Password, got.PgConfig.Password)
					assert.Equal(t, tt.want.PgConfig.DbName, got.PgConfig.DbName)
					assert.Equal(t, tt.want.PgConfig.MaxOpenConns, got.PgConfig.MaxOpenConns)
					assert.Equal(t, tt.want.PgConfig.MaxIdleConns, got.PgConfig.MaxIdleConns)
					assert.Equal(t, tt.want.PgConfig.MinConns, got.PgConfig.MinConns)
					assert.Equal(t, tt.want.PgConfig.MaxConnLifetime, got.PgConfig.MaxConnLifetime)
					assert.Equal(t, tt.want.PgConfig.Timeout, got.PgConfig.Timeout)
				}

				// Проверка LocalConfig
				if tt.want.DataSvc != nil {
					assert.Equal(t, tt.want.DataSvc.Endpoint, got.DataSvc.Endpoint)
				}
			}
		})
	}
}

func TestDefaultEnvironment(t *testing.T) {
	os.Unsetenv("ENVIRONMENT")

	err := os.WriteFile("dev.config.yaml", []byte(`env: dev
log_level: debug
sync:
  httpPort: :8080
  grpcPort: :9090
  dataSvc:
    endpoint: "localhost:50051"
  pgConfig:
    user: postgres
    password: postgres
    dbName: auth
    maxOpenConns: 10
    maxIdleConns: 5
    minConns: 1
    maxConnLifetime: 30s
    timeout: 5s
    host: localhost
    port: 5555
`), 0644)
	require.NoError(t, err)
	defer os.Remove("dev.config.yaml")

	cfg, err := New()
	require.NoError(t, err)
	assert.Equal(t, "dev", cfg.Env)
}
