package app

import (
	"github.com/pkg/errors"

	"github.com/spatecon/gitlab-review-bot/internal/app/service"
)

func (a *App) initService() error {
	var err error

	a.service, err = service.New(a.repository, a.gitlabClient)
	if err != nil {
		return errors.Wrap(err, "failed to init service")
	}

	return nil
}
