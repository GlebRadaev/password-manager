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

// main is the entry point of the application.
// It sets up signal handling for graceful shutdown and starts the application.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app := application.New()
	err := app.Start(ctx)
	if err != nil {
		log.Fatalf("Can't start application: %s", err)
	}

	err = app.Wait(ctx, cancel)
	if err != nil {
		log.Fatalf("All systems closed with errors. LastError: %s", err)
	}
}
