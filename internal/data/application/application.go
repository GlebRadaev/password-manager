package application

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/GlebRadaev/password-manager/internal/common/app"
	"github.com/GlebRadaev/password-manager/internal/common/pg"
	"github.com/GlebRadaev/password-manager/internal/data/api"
	"github.com/GlebRadaev/password-manager/internal/data/config"
	"github.com/GlebRadaev/password-manager/internal/data/migrations"
	"github.com/GlebRadaev/password-manager/internal/data/repo"
	"github.com/GlebRadaev/password-manager/internal/data/service"
	"github.com/GlebRadaev/password-manager/pkg/data"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	AppName = "data"
)

type Application struct {
	cfg  *config.Config
	api  *api.Api
	srv  *service.Service
	repo *repo.Repo

	errCh chan error
	wg    sync.WaitGroup
}

func New() *Application {
	return &Application{
		errCh: make(chan error),
	}
}

func (a *Application) Start(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("can't init config: %w", err)
	}
	a.cfg = cfg
	app.NewLogger(cfg.LogLevel, AppName)

	pool, err := getPgxpool(ctx, cfg.PgConfig)
	if err != nil {
		return fmt.Errorf("can't build pgx pool: %v", err)
	}

	if err = migrations.Exec(pool); err != nil {
		return fmt.Errorf("can't executing migrations: %v", err)
	}

	txManager := pg.NewTXManager(pool)
	a.repo = repo.New(pool)
	a.srv = service.New(a.repo, txManager)
	a.api = api.New(a.srv)

	if err = a.startGrpcServer(ctx); err != nil {
		return fmt.Errorf("can't start grpc server: %w", err)
	}

	if err = a.startHttpServer(ctx); err != nil {
		return fmt.Errorf("can't start http server: %w", err)
	}

	log.Info().Msgf("all systems started successfully")
	return nil
}

func (a *Application) startGrpcServer(ctx context.Context) error {
	lis, err := net.Listen("tcp", a.cfg.GrpcPort)
	if err != nil {
		return fmt.Errorf("error listening on port '%s': %w", a.cfg.GrpcPort, err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpc_recovery.UnaryServerInterceptor()),
	)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		<-ctx.Done()
		s.GracefulStop()
		lis.Close()
	}()

	data.RegisterDataServiceServer(s, a.api)

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

func (a *Application) startHttpServer(ctx context.Context) error {
	conn, err := grpc.NewClient(a.cfg.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	server := http.Server{
		Addr:    a.cfg.HttpPort,
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

		if err := conn.Close(); err != nil {
			a.errCh <- fmt.Errorf("grpc client conn exited with error: %w", err)
		}
	}()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		log.Info().Msgf("starting http server on port %s", a.cfg.HttpPort)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.errCh <- fmt.Errorf("http server exited with error: %w", err)
		}
	}()

	return nil
}

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

func getPgxpool(ctx context.Context, cfg app.PgConfig) (*pgxpool.Pool, error) {
	cfgpool, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, err
	}
	dbpool, err := pgxpool.NewWithConfig(ctx, cfgpool)
	if err != nil {
		return nil, err
	}
	if err = dbpool.Ping(ctx); err != nil {
		return nil, err
	}
	return dbpool, nil
}
