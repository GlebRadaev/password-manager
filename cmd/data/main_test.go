package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockApp struct {
	mock.Mock
}

func (m *MockApp) Start(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *MockApp) Wait(ctx context.Context, cancel context.CancelFunc) error {
	return m.Called(ctx, cancel).Error(0)
}

func TestMain_Success(t *testing.T) {
	origDeps := struct {
		appCreator   func() AppInterface
		logFatal     func(string, ...interface{})
		signalNotify func(ctx context.Context, sig ...os.Signal) (context.Context, context.CancelFunc)
	}{appCreator, logFatal, signalNotify}
	defer func() {
		appCreator = origDeps.appCreator
		logFatal = origDeps.logFatal
		signalNotify = origDeps.signalNotify
	}()

	mockApp := new(MockApp)
	mockApp.On("Start", mock.Anything).Return(nil)
	mockApp.On("Wait", mock.Anything, mock.Anything).Return(nil)
	appCreator = func() AppInterface { return mockApp }

	ctx, cancel := context.WithCancel(context.Background())
	signalNotify = func(_ context.Context, _ ...os.Signal) (context.Context, context.CancelFunc) {
		return ctx, cancel
	}

	var fatalCalled bool
	logFatal = func(string, ...interface{}) {
		fatalCalled = true
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		main()
	}()

	time.Sleep(100 * time.Millisecond)

	cancel()

	wg.Wait()

	assert.False(t, fatalCalled, "log.Fatal был вызван")
	mockApp.AssertExpectations(t)
}

func TestMain_StartError(t *testing.T) {
	origDeps := struct {
		appCreator   func() AppInterface
		logFatal     func(string, ...interface{})
		signalNotify func(ctx context.Context, sig ...os.Signal) (context.Context, context.CancelFunc)
	}{appCreator, logFatal, signalNotify}
	defer func() {
		appCreator = origDeps.appCreator
		logFatal = origDeps.logFatal
		signalNotify = origDeps.signalNotify
	}()

	mockApp := new(MockApp)
	expectedErr := errors.New("start error")
	mockApp.On("Start", mock.Anything).Return(expectedErr)
	appCreator = func() AppInterface { return mockApp }

	signalNotify = func(ctx context.Context, _ ...os.Signal) (context.Context, context.CancelFunc) {
		return ctx, func() {}
	}

	var fatalMsg string
	logFatal = func(format string, args ...interface{}) {
		fatalMsg = fmt.Sprintf(format, args...)
		panic("log.Fatal called")
	}

	func() {
		defer func() {
			if r := recover(); r != "log.Fatal called" {
				t.Errorf("Unexpected panic: %v", r)
			}
		}()
		main()
	}()

	assert.Equal(t, "Can't start application: start error", fatalMsg)
	mockApp.AssertExpectations(t)
	mockApp.AssertNotCalled(t, "Wait")
}

func TestMain_WaitError(t *testing.T) {
	origDeps := struct {
		appCreator   func() AppInterface
		logFatal     func(string, ...interface{})
		signalNotify func(ctx context.Context, sig ...os.Signal) (context.Context, context.CancelFunc)
	}{appCreator, logFatal, signalNotify}
	defer func() {
		appCreator = origDeps.appCreator
		logFatal = origDeps.logFatal
		signalNotify = origDeps.signalNotify
	}()

	mockApp := new(MockApp)
	expectedErr := errors.New("wait error")
	mockApp.On("Start", mock.Anything).Return(nil)
	mockApp.On("Wait", mock.Anything, mock.Anything).Return(expectedErr)
	appCreator = func() AppInterface { return mockApp }

	ctx, cancel := context.WithCancel(context.Background())
	signalNotify = func(_ context.Context, _ ...os.Signal) (context.Context, context.CancelFunc) {
		return ctx, cancel
	}

	var fatalMsg string
	logFatal = func(format string, args ...interface{}) {
		fatalMsg = fmt.Sprintf(format, args...)
		panic("log.Fatal called")
	}

	go func() {
		defer func() {
			if r := recover(); r != "log.Fatal called" {
				t.Errorf("Unexpected panic: %v", r)
			}
		}()
		main()
	}()

	time.Sleep(100 * time.Millisecond)

	cancel()

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, "All systems closed with errors. LastError: wait error", fatalMsg)
	mockApp.AssertExpectations(t)
}
