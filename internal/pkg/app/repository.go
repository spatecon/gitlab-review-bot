package app

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spatecon/gitlab-review-bot/internal/app/repository"
)

const appName = "gitlab-review-bot"

func (a *App) initRepository() error {
	var err error

	URI := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		a.cfg.Mongo.User,
		a.cfg.Mongo.Pass,
		a.cfg.Mongo.Host,
		a.cfg.Mongo.Port)

	a.mongoClient, err = mongo.NewClient(options.Client().SetAppName(appName).ApplyURI(URI))
	if err != nil {
		return errors.Wrap(err, "failed to create mongo client")
	}

	log.Trace().Msg("trying to connect to mongo")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = a.mongoClient.Connect(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to connect to mongo")
	}

	err = a.mongoClient.Ping(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to ping mongo")
	}

	log.Trace().Msg("connected to mongo successfully")

	a.repository, err = repository.New(a.mongoClient, a.cfg.Mongo.DB)
	if err != nil {
		return errors.Wrap(err, "failed to create repository")
	}

	return nil
}
