package app

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

// closer implements graceful shutdown
func (a *App) closer() {
	var err error

	err = a.service.Close()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to close service")
	}

	// close context
	a.closeCtx()

	// close db connection
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err = a.mongoClient.Disconnect(ctx)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to disconnect from mongoDB")
	}

	time.Sleep(2 * time.Second)

	if err == nil {
		log.Info().Msg("gracefully stopped")
	}
}
