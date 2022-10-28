package app

import (
	"github.com/pkg/errors"

	"github.com/spatecon/gitlab-review-bot/internal/pkg/client/gitlab"
)

func (a *App) initClients() error {
	var err error

	a.gitlabClient, err = gitlab.New(a.cfg.GitlabToken)
	if err != nil {
		return errors.Wrap(err, "failed to init gitlab client")
	}

	return nil
}
