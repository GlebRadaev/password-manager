// Package app provides shared application utilities
package app

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const timeLayout = "15:04:05 02-01-2006" // Custom timestamp format

// NewLogger creates configured zerolog logger:
// - Sets log level from envLvl (defaults to Info)
// - Adds app name and timestamp to each log
// - Returns global logger instance
func NewLogger(envLvl string, appName string) *zerolog.Logger {
	lvl, err := zerolog.ParseLevel(envLvl)
	if err != nil {
		lvl = zerolog.InfoLevel // Default level if parsing fails
	}
	zerolog.SetGlobalLevel(lvl)
	zerolog.TimeFieldFormat = timeLayout

	log.Logger = zerolog.New(os.Stdout).With().
		Str("app", appName).
		Timestamp().
		Logger()

	return &log.Logger
}
