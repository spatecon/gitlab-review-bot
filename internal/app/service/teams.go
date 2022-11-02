package service

import (
	"github.com/pkg/errors"
)

func (s *Service) loadTeams() (err error) {
	s.teams, err = s.r.Teams()
	if err != nil {
		return errors.Wrap(err, "failed to load teams")
	}

	return nil
}
