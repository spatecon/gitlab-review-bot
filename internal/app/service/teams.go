package service

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"

	"github.com/spatecon/gitlab-review-bot/internal/app/service/worker"
	"github.com/spatecon/gitlab-review-bot/internal/pkg/logger"
)

func (s *Service) loadTeams() (err error) {
	s.teams, err = s.r.Teams()
	if err != nil {
		return errors.Wrap(err, "failed to load teams")
	}

	return nil
}

func (s *Service) initNotifications() error {
	l := logger.CronLogger{L: log.Logger}

	s.cron = cron.New(cron.WithChain(
		cron.Recover(l),
		cron.SkipIfStillRunning(l),
	))

	for _, team := range s.teams {
		if !team.Notifications.Enabled {
			continue
		}

		if team.Notifications.IsEmpty() {
			continue
		}

		_, ok := s.policies[team.Policy]
		if !ok {
			log.Warn().
				Str("policy", string(team.Policy)).
				Str("team_id", team.ID).
				Msg("unknown policy while initializing notifications")
			continue
		}

		wrk := worker.NewNotificationsWorker(team, s, s.slack)

		_, err := s.cron.AddJob(team.Notifications.Cron, wrk)
		if err != nil {
			return errors.Wrapf(err,
				"failed to add notification worker for %s team (%s)", team.ID, team.Notifications.Cron)
		}
	}

	s.cron.Start()

	return nil
}
