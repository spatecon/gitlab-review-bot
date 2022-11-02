package app

import (
	"github.com/pkg/errors"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
	"github.com/spatecon/gitlab-review-bot/internal/app/service"

	pscoring "github.com/spatecon/gitlab-review-bot/internal/app/service/policy/scoring"
)

func (a *App) initPolicies() error {
	a.policies = make(map[ds.PolicyName]service.Policy)

	a.policies[pscoring.PolicyName] = pscoring.New(a.repository, a.gitlabClient)

	return nil
}

func (a *App) initService() error {
	var err error

	a.service, err = service.New(a.repository, a.gitlabClient, a.policies)
	if err != nil {
		return errors.Wrap(err, "failed to init service")
	}

	return nil
}
