package app

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
	"github.com/spatecon/gitlab-review-bot/internal/app/repository"
	"github.com/spatecon/gitlab-review-bot/internal/app/service"
	"github.com/spatecon/gitlab-review-bot/internal/pkg/client/gitlab"
)

type App struct {
	logger zerolog.Logger
	cfg    Config

	mongoClient *mongo.Client
	repository  *repository.Repository

	gitlabClient *gitlab.Client

	policies map[ds.PolicyName]service.Policy
	service  *service.Service

	// TODO: graceful shutdown
}

func New(configPath string) (*App, error) {
	app := &App{}

	err := app.initConfig(configPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init config")
	}

	app.initLogger()

	err = app.initRepository()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init repository")
	}

	err = app.initClients()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init clients")
	}

	err = app.initPolicies()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init policies")
	}

	err = app.initService()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init service")
	}

	return app, nil
}

func (a *App) Run() error {
	var err error

	a.logger.Info().Msg("app started")

	err = a.service.SubscribeOnProjects()
	if err != nil {
		return errors.Wrap(err, "failed to subscribe on projects")
	}

	ch := make(chan os.Signal)

	<-ch

	return nil
}
