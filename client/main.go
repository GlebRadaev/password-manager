// Package main is the entry point for the password manager client application.
// It initializes and executes the root command of the CLI application.
package main

import (
	"github.com/GlebRadaev/password-manager/client/cmd"
)

// main is the application entry point that delegates execution to the CLI root command.
func main() {
	cmd.Execute()
}
