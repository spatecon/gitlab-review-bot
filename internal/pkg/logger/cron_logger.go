package logger

import (
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

type CronLogger struct {
	L zerolog.Logger
}

var _ cron.Logger = CronLogger{}

func (c CronLogger) Info(msg string, keysAndValues ...interface{}) {
	log := c.L.Info()

	for i := 0; i < len(keysAndValues); i += 2 {
		log = log.Interface(keysAndValues[i].(string), keysAndValues[i+1])
	}

	log.Msg(msg)
}

func (c CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	log := c.L.Error().Err(err)

	for i := 0; i < len(keysAndValues); i += 2 {
		log = log.Interface(keysAndValues[i].(string), keysAndValues[i+1])
	}

	log.Msg(msg)
}
