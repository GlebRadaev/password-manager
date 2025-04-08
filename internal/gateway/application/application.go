// Package application provides the gateway server implementation for the password manager.
// It handles HTTP requests and proxies them to corresponding gRPC services.
package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/GlebRadaev/password-manager/internal/gateway/config"
	"github.com/GlebRadaev/password-manager/internal/gateway/interceptor"
	"github.com/GlebRadaev/password-manager/pkg/auth"
	"github.com/GlebRadaev/password-manager/pkg/data"
	syncsrv "github.com/GlebRadaev/password-manager/pkg/sync"
)

const (
	// AppName defines the application name used in logging and configuration
	AppName = "gateway"
)

// Application represents the gateway server application.
type Application struct {
	cfg   *config.Config
	errCh chan error
	wg    sync.WaitGroup

	ConfigProvider ConfigProvider
}

// New creates and returns a new Application instance with default configuration.
func New() *Application {
	return &Application{
		errCh:          make(chan error),
		ConfigProvider: &defaultConfigProvider{},
	}
}

// Start initializes and starts the HTTP gateway server.
func (a *Application) Start(ctx context.Context) error {
	cfg, err := a.ConfigProvider.New()
	if err != nil {
		return fmt.Errorf("can't init config: %w", err)
	}
	a.cfg = cfg

	if err := a.startHTTPServer(ctx); err != nil {
		return fmt.Errorf("can't start http server: %w", err)
	}

	log.Info().Msgf("gateway server started successfully")
	return nil
}

// startHTTPServer configures and starts the HTTP server with gRPC gateway handlers.
func (a *Application) startHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	authConn, err := grpc.NewClient(
		a.cfg.AuthSvc.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial auth service: %w", err)
	}
	// defer authConn.Close()

	authClient := auth.NewAuthServiceClient(authConn)
	authInterceptor := interceptor.AuthInterceptor(authClient)

	services := map[string]struct {
		endpoint string
		register func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
	}{
		"auth": {
			endpoint: a.cfg.AuthSvc.Endpoint,
			register: func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
				return auth.RegisterAuthServiceHandler(ctx, mux, authConn)
			},
		},
		"data": {
			endpoint: a.cfg.DataSvc.Endpoint,
			register: data.RegisterDataServiceHandler,
		},
		"sync": {
			endpoint: a.cfg.SyncSvc.Endpoint,
			register: syncsrv.RegisterSyncServiceHandler,
		},
	}

	for name, svc := range services {
		conn, err := grpc.NewClient(
			svc.endpoint,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(authInterceptor),
		)
		if err != nil {
			return fmt.Errorf("failed to dial %s service: %w", name, err)
		}

		if err := svc.register(ctx, mux, conn); err != nil {
			return fmt.Errorf("failed to register %s service: %w", name, err)
		}
	}

	server := http.Server{
		Addr:    a.cfg.HTTPPort,
		Handler: mux,
	}

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		<-ctx.Done()

		sCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(sCtx); err != nil {
			a.errCh <- fmt.Errorf("http server exited with error: %w", err)
		}
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		log.Info().Msgf("starting http server on port %s", a.cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errCh <- fmt.Errorf("http server exited with error: %w", err)
		}
	}()

	return nil
}

// Wait blocks until the application context is cancelled and all goroutines finish.
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
