package application

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/GlebRadaev/password-manager/internal/common/app"
	"github.com/GlebRadaev/password-manager/internal/gateway/config"
	"github.com/GlebRadaev/password-manager/pkg/auth"
	"github.com/GlebRadaev/password-manager/pkg/data"
	syncsrv "github.com/GlebRadaev/password-manager/pkg/sync"
)

type mockConfigProvider struct {
	mock.Mock
}

func (m *mockConfigProvider) New() (*config.Config, error) {
	args := m.Called()
	return args.Get(0).(*config.Config), args.Error(1)
}

type mockHTTPServer struct {
	mock.Mock
}

func (m *mockHTTPServer) ListenAndServe() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockHTTPServer) Shutdown(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestNew(t *testing.T) {
	app := New()
	assert.NotNil(t, app)
	assert.NotNil(t, app.errCh)
	assert.IsType(t, &defaultConfigProvider{}, app.ConfigProvider)
}

func TestApplication_Start_ConfigError(t *testing.T) {
	mockProvider := new(mockConfigProvider)
	mockProvider.On("New").Return((*config.Config)(nil), errors.New("config error"))

	app := New()
	app.ConfigProvider = mockProvider

	err := app.Start(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "can't init config")
	mockProvider.AssertExpectations(t)
}

func TestApplication_Start_Success(t *testing.T) {
	mockProvider := new(mockConfigProvider)
	mockProvider.On("New").Return(&config.Config{
		GatewayConfig: config.GatewayConfig{
			HTTPPort: ":8080",
			AuthSvc:  &app.GRPCClient{Endpoint: "auth"},
			DataSvc:  &app.GRPCClient{Endpoint: "data"},
			SyncSvc:  &app.GRPCClient{Endpoint: "sync"},
		},
	}, nil)

	app := New()
	app.ConfigProvider = mockProvider

	// Используем cancel для остановки сервера
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := app.Start(ctx)
	require.NoError(t, err)
	mockProvider.AssertExpectations(t)
}

func TestApplication_ListenAndServe_Error(t *testing.T) {
	mockServer := new(mockHTTPServer)
	mockServer.On("ListenAndServe").Return(errors.New("listen error"))

	app := &Application{
		errCh: make(chan error),
		wg:    sync.WaitGroup{},
	}

	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		if err := mockServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.errCh <- err
		}
	}()

	select {
	case err := <-app.errCh:
		require.Error(t, err)
		assert.Contains(t, err.Error(), "listen error")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected error but got none")
	}

	mockServer.AssertExpectations(t)
}

func TestApplication_Wait_WithError(t *testing.T) {
	app := &Application{
		errCh: make(chan error, 1),
		wg:    sync.WaitGroup{},
	}

	app.errCh <- errors.New("test error")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	var waitErr error
	go func() {
		defer wg.Done()
		waitErr = app.Wait(ctx, cancel)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	wg.Wait()

	require.Error(t, waitErr)
	assert.Equal(t, "test error", waitErr.Error())
}

func TestApplication_Wait_WithoutError(t *testing.T) {
	app := &Application{
		errCh: make(chan error),
		wg:    sync.WaitGroup{},
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	var waitErr error
	go func() {
		defer wg.Done()
		waitErr = app.Wait(ctx, cancel)
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()
	wg.Wait()

	require.NoError(t, waitErr)
}

var (
	grpcNewClient = grpc.NewClient
	authRegister  = auth.RegisterAuthServiceHandler
	dataRegister  = data.RegisterDataServiceHandler
	syncRegister  = syncsrv.RegisterSyncServiceHandler
)
