package config

import (
	"os"
	"testing"

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
gateway:
  httpPort: ":8080"
  grpcPort: ":9090"
  authSvc:
    endpoint: "localhost:50051"
  dataSvc:
    endpoint: "localhost:50052"
  syncSvc:
    endpoint: "localhost:50053"
`,
			want: &Config{
				Env:      "dev",
				LogLevel: "debug",
				GatewayConfig: GatewayConfig{
					HTTPPort: ":8080",
					GrpcPort: ":9090",
					AuthSvc: &app.GRPCClient{
						Endpoint: "localhost:50051",
					},
					DataSvc: &app.GRPCClient{
						Endpoint: "localhost:50052",
					},
					SyncSvc: &app.GRPCClient{
						Endpoint: "localhost:50053",
					},
				},
			},
			wantErr: false,
			setup: func() {
				os.Setenv("ENVIRONMENT", "dev")
				err := os.WriteFile("dev.config.yaml", []byte(`env: dev
log_level: debug
gateway:
  httpPort: ":8080"
  grpcPort: ":9090"
  authSvc:
    endpoint: "localhost:50051"
  dataSvc:
    endpoint: "localhost:50052"
  syncSvc:
    endpoint: "localhost:50053"
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
gateway:
  httpPort: ":8080"
`,
			wantErr:     true,
			errContains: "missing required attributes",
			setup: func() {
				os.Setenv("ENVIRONMENT", "validation")
				err := os.WriteFile("validation.config.yaml", []byte(`env: dev
log_level: debug
gateway:
  httpPort: ":8080"
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

				// Check GatewayConfig
				assert.Equal(t, tt.want.HTTPPort, got.HTTPPort)
				assert.Equal(t, tt.want.GrpcPort, got.GrpcPort)

				// Check GRPC clients
				if tt.want.AuthSvc != nil {
					assert.Equal(t, tt.want.AuthSvc.Endpoint, got.AuthSvc.Endpoint)
				}
				if tt.want.DataSvc != nil {
					assert.Equal(t, tt.want.DataSvc.Endpoint, got.DataSvc.Endpoint)
				}
				if tt.want.SyncSvc != nil {
					assert.Equal(t, tt.want.SyncSvc.Endpoint, got.SyncSvc.Endpoint)
				}
			}
		})
	}
}

func TestDefaultEnvironment(t *testing.T) {
	originalEnv := os.Getenv("ENVIRONMENT")
	defer func() {
		os.Setenv("ENVIRONMENT", originalEnv)
	}()

	os.Unsetenv("ENVIRONMENT")

	err := os.WriteFile("dev.config.yaml", []byte(`env: dev
log_level: debug
gateway:
  httpPort: ":8080"
  grpcPort: ":9090"
  authSvc:
    endpoint: "localhost:50051"
  dataSvc:
    endpoint: "localhost:50052"
  syncSvc:
    endpoint: "localhost:50053"
`), 0644)
	require.NoError(t, err)
	defer os.Remove("dev.config.yaml")

	cfg, err := New()
	require.NoError(t, err)
	assert.Equal(t, "dev", cfg.Env)
}

func TestRequiredFields(t *testing.T) {
	originalEnv := os.Getenv("ENVIRONMENT")
	defer func() {
		os.Setenv("ENVIRONMENT", originalEnv)
	}()

	t.Run("missing required fields", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "testreq")
		err := os.WriteFile("testreq.config.yaml", []byte(`env: testreq
log_level: debug
`), 0644)
		require.NoError(t, err)
		defer os.Remove("testreq.config.yaml")

		_, err = New()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing required attributes")
	})
}
