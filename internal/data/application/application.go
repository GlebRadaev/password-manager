// Package application provides the main application logic for the data service.
// It handles initialization, configuration, and lifecycle management of gRPC and HTTP servers.
package application

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/GlebRadaev/password-manager/internal/common/app"
	"github.com/GlebRadaev/password-manager/internal/common/pg"
	"github.com/GlebRadaev/password-manager/internal/common/swagger"
	"github.com/GlebRadaev/password-manager/internal/data/api"
	"github.com/GlebRadaev/password-manager/internal/data/config"
	"github.com/GlebRadaev/password-manager/internal/data/repo"
	"github.com/GlebRadaev/password-manager/internal/data/service"
	"github.com/GlebRadaev/password-manager/pkg/data"
)

const (
	// AppName defines the application name used in logging and configuration
	AppName = "data"
)

// Application represents the main application struct that orchestrates all components
type Application struct {
	cfg  *config.Config
	api  *api.API
	srv  *service.Service
	repo *repo.Repo

	errCh chan error
	wg    sync.WaitGroup

	// Dependency injection interfaces for testing
	ConfigProvider     ConfigProvider
	PgxPoolProvider    PgxPoolProvider
	MigrationsExecutor MigrationsExecutor
	NewListener        func(network, address string) (net.Listener, error)
	NewGrpcServer      func(...grpc.ServerOption) GrpcServer
	NewHTTPServer      func(addr string, handler http.Handler) HTTPServer
	NewGrpcClient      func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
}

// New creates a new Application instance with default providers
func New() *Application {
	return &Application{
		errCh:              make(chan error),
		ConfigProvider:     &defaultConfigProvider{},
		PgxPoolProvider:    &defaultPgxPoolProvider{},
		MigrationsExecutor: &defaultMigrationsExecutor{},
		NewListener:        net.Listen,
		NewGrpcServer: func(opts ...grpc.ServerOption) GrpcServer {
			return &grpcServerWrapper{
				Server: grpc.NewServer(opts...),
			}
		},
		NewHTTPServer: func(addr string, handler http.Handler) HTTPServer {
			return &httpServerWrapper{
				Server: &http.Server{
					Addr:    addr,
					Handler: handler,
				},
			}
		},
		NewGrpcClient: grpc.NewClient,
	}
}

// Start initializes and starts the application components
// Returns error if any initialization step fails
func (a *Application) Start(ctx context.Context) error {
	cfg, err := a.ConfigProvider.New()
	if err != nil {
		return fmt.Errorf("can't init config: %w", err)
	}
	a.cfg = cfg
	app.NewLogger(cfg.LogLevel, AppName)

	pool, err := a.PgxPoolProvider.GetPgxpool(ctx, cfg.PgConfig)
	if err != nil {
		return fmt.Errorf("can't build pgx pool: %v", err)
	}

	if err = a.MigrationsExecutor.Exec(pool); err != nil {
		return fmt.Errorf("can't executing migrations: %v", err)
	}

	txManager := pg.NewTXManager(pool)
	a.repo = repo.New(pool)
	a.srv = service.New(a.repo, txManager)
	a.api = api.New(a.srv)

	if err = a.startGrpcServer(ctx); err != nil {
		return fmt.Errorf("can't start grpc server: %w", err)
	}

	if err = a.startHTTPServer(ctx); err != nil {
		return fmt.Errorf("can't start http server: %w", err)
	}

	log.Info().Msgf("all systems started successfully")
	return nil
}

// startGrpcServer initializes and starts the gRPC server
func (a *Application) startGrpcServer(ctx context.Context) error {
	lis, err := a.NewListener("tcp", a.cfg.GrpcPort)
	if err != nil {
		return fmt.Errorf("error listening on port '%s': %w", a.cfg.GrpcPort, err)
	}

	s := a.NewGrpcServer(
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor()),
	)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		<-ctx.Done()
		s.GracefulStop()
		lis.Close()
	}()

	s.RegisterDataService(a.api)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		log.Info().Msgf("starting gRPC server on port %s", a.cfg.GrpcPort)
		if err = s.Serve(lis); err != nil {
			a.errCh <- fmt.Errorf("grpc server exited with error: %w", err)
		}
	}()

	return nil
}

// startHttpServer initializes and starts the HTTP gateway server
func (a *Application) startHTTPServer(ctx context.Context) error {
	conn, err := a.NewGrpcClient(a.cfg.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error conn to gRPC server: %w", err)
	}

	mux := grpc_runtime.NewServeMux(
		grpc_runtime.WithIncomingHeaderMatcher(grpc_runtime.DefaultHeaderMatcher),
	)

	err = data.RegisterDataServiceHandler(ctx, mux, conn)
	if err != nil {
		return fmt.Errorf("error register auth service: %w", err)
	}

	handler := swagger.RegisterSwaggerUI(mux, AppName)

	server := a.NewHTTPServer(a.cfg.HTTPPort, handler)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		<-ctx.Done()

		sCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(sCtx); err != nil {
			a.errCh <- fmt.Errorf("http server exited with error: %w", err)
		}

		if err := conn.Close(); err != nil {
			a.errCh <- fmt.Errorf("grpc client conn exited with error: %w", err)
		}
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		log.Info().Msgf("starting http server on port %s", a.cfg.HTTPPort)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errCh <- fmt.Errorf("http server exited with error: %w", err)
		}
	}()

	return nil
}

// Wait blocks until application shutdown is complete
// Returns the first error that caused shutdown or nil if graceful shutdown
func (a *Application) Wait(ctx context.Context, cancel context.CancelFunc) error {
	var appErr error

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for err := range a.errCh {
			cancel()
			log.Error().Err(err).Send()
			appErr = err
		}
	}()

	<-ctx.Done()
	a.wg.Wait()
	close(a.errCh)
	wg.Wait()

	return appErr
}

// GetPgxpool creates a new pgx connection pool with the given configuration
func GetPgxpool(ctx context.Context, cfg app.PgConfig) (*pgxpool.Pool, error) {
	cfgpool, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse pg config: %w", err)
	}
	dbpool, err := pgxpool.NewWithConfig(ctx, cfgpool)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}
	if err = dbpool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return dbpool, nil
}

// GetErrCh returns the application error channel
func (a *Application) GetErrCh() chan error {
	return a.errCh
}
