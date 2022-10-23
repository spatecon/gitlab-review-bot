package app

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type App struct {
	logger zerolog.Logger
	cfg    Config
}

func New(configPath string) (*App, error) {
	app := &App{}

	err := app.initConfig(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init config")
	}

	app.initLogger()

	return app, nil
}

func (a *App) Run() error {
	a.logger.Info().Str("v", a.cfg.GitlabToken).Msg("app started")
	return nil
}
