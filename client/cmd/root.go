// Package cmd implements the command-line interface for the password manager
package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

// Build-time variables (set via ldflags)
var (
	version = "1.0.0"   // Application version
	commit  = "none"    // Git commit hash
	date    = "unknown" // Build date
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pm",
	Short: "Password Manager CLI",
	Long:  "GophKeeper - client-server password manager with local storage and sync capabilities",
	Version: fmt.Sprintf("%s\nCommit: %s\nBuild date: %s\nGo version: %s",
		version,
		commit,
		date,
		runtime.Version()),
}

// Execute runs the root command and handles errors
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
