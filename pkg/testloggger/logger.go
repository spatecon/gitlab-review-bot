package testloggger

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(t *testing.T) zerolog.Logger {
	return zerolog.New(zerolog.NewTestWriter(t))
}

// Set new test logger as global logger
func Set(t *testing.T) {
	log.Logger = New(t)
}

// Unset replace global logger with Nop logger
// usually used with defer after Set call
func Unset() {
	log.Logger = zerolog.Nop()
}
