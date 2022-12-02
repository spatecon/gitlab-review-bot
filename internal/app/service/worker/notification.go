package worker

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

const notificationWorkerName = "notification_worker"

type Policy interface {
	IsApproved(team *ds.Team, mr *ds.MergeRequest) bool
}

type SlackClient interface {
	SendMessage(recipientID string, message string) error
}

type NotificationService interface {
	GetAuthoredReviewedMRs(team *ds.Team, users []*ds.User) (authorToMR, reviewerToMR map[int][]*ds.MergeRequest, err error)
	UserNotification(user *ds.User, team *ds.Team, authorToMR, reviewerToMR map[int][]*ds.MergeRequest) (message string, err error)
	TeamNotification(team *ds.Team, authorToMR, reviewerToMR map[int][]*ds.MergeRequest) (message string, err error)
}

type Notifications struct {
	team  *ds.Team
	svc   NotificationService
	slack SlackClient
}

func NewNotificationsWorker(team *ds.Team, svc NotificationService, slack SlackClient) *Notifications {
	return &Notifications{
		team:  team,
		svc:   svc,
		slack: slack,
	}
}

type SlackMessage struct {
	RecipientID string
	Text        string
}

func (n *Notifications) Run() {
	l := log.With().
		Str("worker", notificationWorkerName).
		Str("team_id", n.team.ID).
		Logger()

	members := make([]*ds.User, 0, len(n.team.Members))
	for _, member := range n.team.Members {
		// only leads and developers are notified
		if !member.Labels.Has(ds.DeveloperLabel) && !member.Labels.Has(ds.LeadLabel) {
			continue
		}

		members = append(members, member)
	}

	authorToMR, reviewerToMR, err := n.svc.GetAuthoredReviewedMRs(n.team, members)
	if err != nil {
		l.Error().Err(err).Msg("failed to fetch MRs in notifications")
		return
	}

	slackMessages, err := n.slackMessages(members, authorToMR, reviewerToMR)
	if err != nil {
		l.Error().Err(err).Msg("failed to generate slack messages")
		return
	}

	for _, message := range slackMessages {
		if strings.TrimSpace(message.Text) == "" {
			continue
		}

		err = n.slack.SendMessage(message.RecipientID, message.Text)
		if err != nil {
			l.Error().Err(err).Str("recipient_id", message.RecipientID).Msg("failed to send slack message")
			return
		}
	}
}

func (n *Notifications) slackMessages(
	devs []*ds.User,
	authorToMR, reviewerToMR map[int][]*ds.MergeRequest,
) ([]SlackMessage, error) {

	slackMessages := make([]SlackMessage, 0, len(devs)+1)

	for _, dev := range devs {
		message, err := n.svc.UserNotification(dev, n.team, authorToMR, reviewerToMR)
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate user notification")
		}

		slackMessages = append(slackMessages, SlackMessage{
			RecipientID: dev.SlackID,
			Text:        message,
		})
	}

	teamMessage, err := n.svc.TeamNotification(n.team, authorToMR, reviewerToMR)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate team notification")
	}

	slackMessages = append(slackMessages, SlackMessage{
		RecipientID: n.team.Notifications.ChannelID,
		Text:        teamMessage,
	})

	return slackMessages, nil
}
