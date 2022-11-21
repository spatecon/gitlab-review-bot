package worker

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

type SlackWorkerRepository interface {
	UserBySlackID(slackID string) (*ds.User, *ds.Team, error)
}

type SlackWorker struct {
	svc    NotificationService
	r      SlackWorkerRepository
	slack  SlackClient
	events chan ds.UserEvent
	close  chan struct{}
}

func (w *SlackWorker) Close() {
	w.close <- struct{}{}
}

func NewSlackWorker(svc NotificationService, r SlackWorkerRepository, slack SlackClient, events chan ds.UserEvent) *SlackWorker {
	return &SlackWorker{
		svc:    svc,
		r:      r,
		slack:  slack,
		events: events,
		close:  make(chan struct{}),
	}
}

func (w *SlackWorker) Run() {
	for {
		select {
		case <-w.close:
			return
		case event := <-w.events:
			err := w.processEvent(event)
			if err != nil {
				log.Error().Err(err).Msg("failed to process slack event")
			}
		}
	}
}

func (w *SlackWorker) processEvent(event ds.UserEvent) error {
	if event.Type != ds.UserEventTypeMRRequest {
		return nil
	}

	user, team, err := w.r.UserBySlackID(event.UserID)
	if err != nil {
		return errors.Wrap(err, "failed to get user by slack id")
	}

	if user == nil {
		return nil
	}

	authorToMR, reviewerToMR, err := w.svc.GetAuthoredReviewedMRs(team, []*ds.User{user})
	if err != nil {
		return errors.Wrap(err, "failed to get authored and reviewed mrs")
	}

	msg, err := w.svc.UserNotification(user, team, authorToMR, reviewerToMR)
	if err != nil {
		return errors.Wrap(err, "failed to get user notification")
	}

	err = w.slack.SendMessage(event.UserID, msg)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

	return nil
}
