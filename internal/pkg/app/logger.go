package app

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func (a *App) initLogger() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if a.cfg.HumanReadableLog {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.Stamp,
		}
		log.Logger = log.Output(consoleWriter)
		a.logger = zerolog.New(consoleWriter).
			Level(zerolog.InfoLevel).
			With().Timestamp().
			Logger()
		return
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	a.logger = zerolog.New(os.Stderr).
		Level(zerolog.ErrorLevel).
		With().Timestamp().
		Logger()
}
