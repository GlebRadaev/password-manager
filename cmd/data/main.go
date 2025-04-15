// Package main initializes and starts the password manager application.
// It handles graceful shutdown on system interrupt or termination signals.
package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/GlebRadaev/password-manager/internal/data/application"
)

var (
	// appCreator creates a new application instance. Can be replaced in tests.
	appCreator = func() AppInterface {
		return application.New()
	}
	// logFatal logs a fatal error and exits the program. Can be replaced in tests.
	logFatal = log.Fatalf
	// signalNotify sets up signal notification. Can be replaced in tests.
	signalNotify = signal.NotifyContext
)

// AppInterface defines the contract for the application instance.
type AppInterface interface {
	// Start initializes and starts the application components.
	Start(ctx context.Context) error
	// Wait blocks until application shutdown is complete.
	Wait(ctx context.Context, cancel context.CancelFunc) error
}

// main is the entry point of the application.
// It sets up signal handling for graceful shutdown and starts the application.

func main() {
	ctx, cancel := signalNotify(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app := appCreator()
	if err := app.Start(ctx); err != nil {
		logFatal("Can't start application: %s", err)
	}

	if err := app.Wait(ctx, cancel); err != nil {
		logFatal("All systems closed with errors. LastError: %s", err)
	}
}
