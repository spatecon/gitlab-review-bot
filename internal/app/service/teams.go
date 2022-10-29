package service

import (
	"github.com/pkg/errors"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
	"github.com/spatecon/gitlab-review-bot/internal/app/fsm/basic"
	"github.com/spatecon/gitlab-review-bot/internal/pkg/fsm"
)

func (s *Service) loadTeams() (err error) {
	s.teams, err = s.repo.Teams()
	if err != nil {
		return errors.Wrap(err, "failed to load teams")
	}

	return nil
}

func (s *Service) initStateMachines() error {
	s.stateMachines = make([]fsm.StateMachine, 0, len(s.teams))

	for _, team := range s.teams {
		sm, err := s.stateMachine(team)
		if err != nil {
			return errors.Wrap(err, "failed to init state machine")
		}

		s.stateMachines = append(s.stateMachines, sm)
	}

	return nil
}

func (s *Service) stateMachine(team *ds.Team) (fsm.StateMachine, error) {
	switch team.Policy {
	case "basic":
		return basic.New(team)
	}

	return nil, errors.Errorf("unknown policy: %s", team.Policy)
}
