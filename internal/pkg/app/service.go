package app

import (
	"github.com/pkg/errors"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
	"github.com/spatecon/gitlab-review-bot/internal/app/service"

	prd "github.com/spatecon/gitlab-review-bot/internal/app/service/policy/reinventing-democracy"
)

func (a *App) initPolicies() error {
	a.policies = make(map[ds.PolicyName]service.Policy)

	a.policies[prd.PolicyName] = prd.New(a.repository, a.gitlabClient)

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
