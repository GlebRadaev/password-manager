// Package application provides the core application logic for the authentication service.
// It manages database connections, migrations, and serves both gRPC and HTTP endpoints
// for authentication-related operations.
package application

import (
	"context"
	"net"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	"github.com/GlebRadaev/password-manager/internal/auth/config"
	"github.com/GlebRadaev/password-manager/internal/auth/migrations"
	"github.com/GlebRadaev/password-manager/internal/common/app"
	"github.com/GlebRadaev/password-manager/pkg/auth"
)

// ConfigProvider defines the interface for providing service configuration.
type ConfigProvider interface {
	// New creates and returns a new configuration instance.
	New() (*config.Config, error)
}

// PgxPoolProvider defines the interface for creating PostgreSQL connection pools.
type PgxPoolProvider interface {
	// GetPgxpool creates a new PostgreSQL connection pool with the given configuration.
	GetPgxpool(ctx context.Context, cfg app.PgConfig) (*pgxpool.Pool, error)
}

// MigrationsExecutor defines the interface for executing database migrations.
type MigrationsExecutor interface {
	// Exec runs all pending database migrations using the provided connection pool.
	Exec(pool *pgxpool.Pool) error
}

// GrpcServer defines the interface for gRPC server operations.
type GrpcServer interface {
	// Serve starts accepting incoming connections on the listener.
	Serve(net.Listener) error
	// GracefulStop gracefully stops the server.
	GracefulStop()
	// RegisterService registers a service implementation with the gRPC server.
	RegisterService(desc *grpc.ServiceDesc, impl interface{})
	// RegisterAuthService specifically registers the AuthService implementation.
	RegisterAuthService(srv auth.AuthServiceServer)
}

// HTTPServer defines the interface for HTTP server operations.
type HTTPServer interface {
	// ListenAndServe starts the HTTP server.
	ListenAndServe() error
	// Shutdown gracefully shuts down the server.
	Shutdown(ctx context.Context) error
}

// grpcServerWrapper wraps a gRPC server to implement the GrpcServer interface.
type grpcServerWrapper struct {
	*grpc.Server
}

// Serve starts the gRPC server on the provided listener.
func (w *grpcServerWrapper) Serve(lis net.Listener) error {
	return w.Server.Serve(lis)
}

// GracefulStop performs a graceful shutdown of the gRPC server.
func (w *grpcServerWrapper) GracefulStop() {
	w.Server.GracefulStop()
}

// RegisterService registers a gRPC service implementation.
func (w *grpcServerWrapper) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	w.Server.RegisterService(desc, impl)
}

// RegisterAuthService registers the AuthService implementation.
func (w *grpcServerWrapper) RegisterAuthService(srv auth.AuthServiceServer) {
	auth.RegisterAuthServiceServer(w.Server, srv)
}

// httpServerWrapper wraps an HTTP server to implement the HTTPServer interface.
type httpServerWrapper struct {
	*http.Server
}

// ListenAndServe starts the HTTP server.
func (w *httpServerWrapper) ListenAndServe() error {
	return w.Server.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server.
func (w *httpServerWrapper) Shutdown(ctx context.Context) error {
	return w.Server.Shutdown(ctx)
}

// defaultConfigProvider provides the default configuration implementation.
type defaultConfigProvider struct{}

// New creates a new configuration instance using the default provider.
func (d *defaultConfigProvider) New() (*config.Config, error) {
	return config.New()
}

// defaultPgxPoolProvider provides the default PostgreSQL connection pool implementation.
type defaultPgxPoolProvider struct{}

// GetPgxpool creates a new PostgreSQL connection pool using default settings.
func (d *defaultPgxPoolProvider) GetPgxpool(ctx context.Context, cfg app.PgConfig) (*pgxpool.Pool, error) {
	return GetPgxpool(ctx, cfg)
}

// defaultMigrationsExecutor provides the default database migrations implementation.
type defaultMigrationsExecutor struct{}

// Exec runs all pending database migrations using the default executor.
func (d *defaultMigrationsExecutor) Exec(pool *pgxpool.Pool) error {
	return migrations.Exec(pool)
}
