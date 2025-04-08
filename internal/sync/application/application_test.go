package application_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"github.com/GlebRadaev/password-manager/internal/common/app"
	"github.com/GlebRadaev/password-manager/internal/sync/application"
	"github.com/GlebRadaev/password-manager/internal/sync/config"
	"github.com/GlebRadaev/password-manager/pkg/sync"
)

type ApplicationSuite struct {
	suite.Suite
	app *application.Application

	configProvider  *mockConfigProvider
	pgxPoolProvider *mockPgxPoolProvider
	migrationsExec  *mockMigrationsExecutor
}

func TestApplication(t *testing.T) {
	suite.Run(t, &ApplicationSuite{})
}

type mockConfigProvider struct {
	config *config.Config
	err    error
}

func (m *mockConfigProvider) New() (*config.Config, error) {
	return m.config, m.err
}

type mockPgxPoolProvider struct {
	pool *pgxpool.Pool
	err  error
}

func (m *mockPgxPoolProvider) GetPgxpool(ctx context.Context, cfg app.PgConfig) (*pgxpool.Pool, error) {
	return m.pool, m.err
}

type mockMigrationsExecutor struct {
	err error
}

func (m *mockMigrationsExecutor) Exec(pool *pgxpool.Pool) error {
	return m.err
}

type mockListener struct {
	net.Listener
	acceptFn func() (net.Conn, error)
}

func (m *mockListener) Accept() (net.Conn, error) {
	if m.acceptFn != nil {
		return m.acceptFn()
	}
	return nil, nil
}

func (m *mockListener) Close() error {
	return nil
}

func (m *mockListener) Addr() net.Addr {
	return &net.TCPAddr{}
}

type mockGrpcServer struct {
	serveFn        func(net.Listener) error
	gracefulStopFn func()
	registerFn     func(srv sync.SyncServiceServer)
}

func (m *mockGrpcServer) Serve(lis net.Listener) error {
	if m.serveFn != nil {
		return m.serveFn(lis)
	}
	return nil
}

func (m *mockGrpcServer) GracefulStop() {
	if m.gracefulStopFn != nil {
		m.gracefulStopFn()
	}
}

func (m *mockGrpcServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {}

func (m *mockGrpcServer) RegisterSyncService(srv sync.SyncServiceServer) {
	if m.registerFn != nil {
		m.registerFn(srv)
	}
}

type mockHttpServer struct {
	listenAndServeFn func() error
	shutdownFn       func(ctx context.Context) error
}

func (m *mockHttpServer) ListenAndServe() error {
	if m.listenAndServeFn != nil {
		return m.listenAndServeFn()
	}
	return nil
}

func (m *mockHttpServer) Shutdown(ctx context.Context) error {
	if m.shutdownFn != nil {
		return m.shutdownFn(ctx)
	}
	return nil
}

func (s *ApplicationSuite) SetupTest() {
	grpcPort := fmt.Sprintf(":%d", 50051+rand.Intn(1000))
	httpPort := fmt.Sprintf(":%d", 8080+rand.Intn(1000))
	s.app = application.New()
	s.configProvider = &mockConfigProvider{
		config: &config.Config{
			Env:      "test",
			LogLevel: "debug",
			Config: app.Config{
				HTTPPort: httpPort,
				GrpcPort: grpcPort,
				PgConfig: app.PgConfig{
					Host:            "localhost",
					Port:            5432,
					User:            "user",
					Password:        "password",
					DbName:          "dbname",
					MaxOpenConns:    10,
					MaxIdleConns:    5,
					MinConns:        2,
					MaxConnLifetime: time.Minute * 30,
					Timeout:         time.Second * 5,
				},
			},
			LocalConfig: config.LocalConfig{
				DataSvc: &app.GRPCClient{
					Endpoint: "localhost:50000",
				},
			},
		},
	}

	s.pgxPoolProvider = &mockPgxPoolProvider{}
	s.migrationsExec = &mockMigrationsExecutor{}

	s.app.ConfigProvider = s.configProvider
	s.app.PgxPoolProvider = s.pgxPoolProvider
	s.app.MigrationsExecutor = s.migrationsExec
}

func (s *ApplicationSuite) TestStart_Success() {
	ctx := context.Background()

	s.pgxPoolProvider.pool = &pgxpool.Pool{}
	s.migrationsExec.err = nil

	err := s.app.Start(ctx)

	if errors.Is(err, grpc.ErrServerStopped) {
		return
	}

	s.NoError(err)
}

func (s *ApplicationSuite) TestStart_ConfigError() {
	s.configProvider.err = errors.New("config error")

	err := s.app.Start(context.Background())

	s.Error(err)
	s.Contains(err.Error(), "can't init config")
}

func (s *ApplicationSuite) TestStart_DBConnectionError() {
	s.pgxPoolProvider.err = errors.New("connection error")

	err := s.app.Start(context.Background())

	s.Error(err)
	s.Contains(err.Error(), "can't build pgx pool")
}

func (s *ApplicationSuite) TestStart_MigrationsError() {
	s.pgxPoolProvider.pool = &pgxpool.Pool{}
	s.migrationsExec.err = errors.New("migrations error")

	err := s.app.Start(context.Background())

	s.Error(err)
	s.Contains(err.Error(), "can't executing migrations")
}

func (s *ApplicationSuite) TestStart_GrpcServerListenError() {
	s.pgxPoolProvider.pool = &pgxpool.Pool{}
	s.migrationsExec.err = nil

	originalListener := s.app.NewListener
	s.app.NewListener = func(network, address string) (net.Listener, error) {
		return nil, errors.New("listen error")
	}
	defer func() { s.app.NewListener = originalListener }()

	err := s.app.Start(context.Background())
	s.Error(err)
	s.Contains(err.Error(), "can't start grpc server")
	s.Contains(err.Error(), "listen error")
}

func (s *ApplicationSuite) TestStart_GrpcServerServeError() {
	s.pgxPoolProvider.pool = &pgxpool.Pool{}
	s.migrationsExec.err = nil

	serveCalled := make(chan struct{}, 1)

	originalGrpcServer := s.app.NewGrpcServer
	s.app.NewGrpcServer = func(opts ...grpc.ServerOption) application.GrpcServer {
		return &mockGrpcServer{
			serveFn: func(net.Listener) error {
				serveCalled <- struct{}{}
				return errors.New("serve error")
			},
			registerFn: func(srv sync.SyncServiceServer) {},
		}
	}
	defer func() { s.app.NewGrpcServer = originalGrpcServer }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.app.Start(ctx)
	}()

	select {
	case <-serveCalled:
		select {
		case err := <-s.app.GetErrCh():
			s.Error(err)
			s.Contains(err.Error(), "grpc server exited with error")
			s.Contains(err.Error(), "serve error")
		case <-time.After(500 * time.Millisecond):
			s.Fail("timeout waiting for error in app.errCh")
		}
	case <-time.After(500 * time.Millisecond):
		s.Fail("timeout waiting for Serve to be called")
	}
}

func (s *ApplicationSuite) TestStart_HttpServerConnError() {
	s.pgxPoolProvider.pool = &pgxpool.Pool{}
	s.migrationsExec.err = nil

	originalGrpcClient := s.app.NewGrpcClient
	s.app.NewGrpcClient = func(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
		return nil, errors.New("connection error")
	}
	defer func() { s.app.NewGrpcClient = originalGrpcClient }()

	err := s.app.Start(context.Background())
	s.Error(err)
	s.Contains(err.Error(), "can't start http server")
	s.Contains(err.Error(), "connection error")
}

func (s *ApplicationSuite) TestGrpcServerErrorChannel() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.pgxPoolProvider.pool = &pgxpool.Pool{}
	s.migrationsExec.err = nil

	originalGrpcServer := s.app.NewGrpcServer
	s.app.NewGrpcServer = func(opts ...grpc.ServerOption) application.GrpcServer {
		return &mockGrpcServer{
			serveFn: func(net.Listener) error {
				return errors.New("server error")
			},
		}
	}
	defer func() { s.app.NewGrpcServer = originalGrpcServer }()

	go func() {
		err := s.app.Start(ctx)
		s.NoError(err)
	}()

	select {
	case err := <-s.app.GetErrCh():
		s.Error(err)
		s.Contains(err.Error(), "grpc server exited with error")
	case <-time.After(100 * time.Millisecond):
		s.Fail("timeout waiting for error")
	}
}

func (s *ApplicationSuite) TestHttpServerErrorChannel() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.pgxPoolProvider.pool = &pgxpool.Pool{}
	s.migrationsExec.err = nil

	originalHttpServer := s.app.NewHTTPServer
	s.app.NewHTTPServer = func(addr string, handler http.Handler) application.NewHTTPServer {
		return &mockHttpServer{
			listenAndServeFn: func() error {
				return errors.New("http server error")
			},
		}
	}
	defer func() { s.app.NewHTTPServer = originalHttpServer }()

	go func() {
		err := s.app.Start(ctx)
		s.NoError(err)
	}()

	select {
	case err := <-s.app.GetErrCh():
		s.Error(err)
		s.Contains(err.Error(), "http server exited with error")
	case <-time.After(100 * time.Millisecond):
		s.Fail("timeout waiting for error")
	}
}

func (s *ApplicationSuite) TestWait_WithError() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		s.app.GetErrCh() <- errors.New("test error")
	}()

	err := s.app.Wait(ctx, cancel)

	s.Error(err)
	s.Contains(err.Error(), "test error")
}

func (s *ApplicationSuite) TestWait_ContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err := s.app.Wait(ctx, cancel)

	s.NoError(err)
}

func (s *ApplicationSuite) TestGetPgxpool_Success() {
	s.T().Skip("This test requires a real database connection")
}

func (s *ApplicationSuite) TestGetPgxpool_InvalidConfig() {
	ctx := context.Background()

	invalidCfg := app.PgConfig{
		Host: "",
		Port: 0,
	}

	_, err := application.GetPgxpool(ctx, invalidCfg)
	s.Error(err)
	s.Contains(err.Error(), "failed to parse pg config")
	s.Contains(err.Error(), "invalid port")
}

func (s *ApplicationSuite) TestGetPgxpool_ContextCancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg := app.PgConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "testuser",
		Password:        "testpass",
		DbName:          "testdb",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		MinConns:        2,
		MaxConnLifetime: time.Minute * 30,
		Timeout:         time.Second * 5,
	}

	_, err := application.GetPgxpool(ctx, cfg)
	s.Error(err)
	s.Contains(err.Error(), "context canceled")
}
