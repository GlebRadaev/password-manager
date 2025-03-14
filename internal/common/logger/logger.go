package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"honnef.co/go/tools/config"
)

const timeLayout = "15:04:05 02-01-2006"

func NewLogger(cfg *config.Config, appName string) *zerolog.Logger {
	lvl, err := zerolog.ParseLevel(cfg.LogLvl)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
	zerolog.TimeFieldFormat = timeLayout

	log.Logger = zerolog.New(os.Stdout).With().
		Str("app", appName).
		Timestamp().
		Logger()

	return &log.Logger
}
