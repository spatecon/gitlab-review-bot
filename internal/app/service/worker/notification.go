package worker

import (
	"bytes"
	"math"
	"text/template"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

type Policy interface {
	IsApproved(team *ds.Team, mr *ds.MergeRequest) bool
}

type SlackClient interface {
	SendMessage(recipientID string, message string) error
}

type Repository interface {
	MergeRequestsByAuthor(authorID []int) ([]*ds.MergeRequest, error)
	MergeRequestsByReviewer(reviewerID []int) ([]*ds.MergeRequest, error)
}

type Notifications struct {
	team       *ds.Team
	policy     Policy
	repository Repository
	slack      SlackClient
}

func NewNotificationsWorker(team *ds.Team, policy Policy, repository Repository, slack SlackClient) *Notifications {
	return &Notifications{
		team:       team,
		policy:     policy,
		repository: repository,
		slack:      slack,
	}
}

type SlackMessage struct {
	RecipientID string
	Text        string
}

func (n *Notifications) Run() {
	l := log.With().Str("team", n.team.Name).Str("module", "notifications_worker").Logger()

	devs := ds.Developers(n.team.Members)

	authorToMR, reviewerToMR, err := n.getMergeRequests(devs)
	if err != nil {
		l.Error().Err(err).Msg("failed to fetch MRs in notifications")
		return
	}

	slackMessages, err := n.slackMessages(l, devs, authorToMR, reviewerToMR)
	if err != nil {
		l.Error().Err(err).Msg("failed to generate slack messages")
		return
	}

	for _, message := range slackMessages {
		err = n.slack.SendMessage(message.RecipientID, message.Text)
		if err != nil {
			l.Error().Err(err).Str("recipient_id", message.RecipientID).Msg("failed to send slack message")
			return
		}
	}
}

func (n *Notifications) getMergeRequests(devs []*ds.User) (
	authorToMR, reviewerToMR map[int][]*ds.MergeRequest,
	err error,
) {
	// TODO: optimize to not fetch all MRs (inc. closed, merged, etc.)
	authoredMRs, err := n.repository.MergeRequestsByAuthor(lo.Map(devs, func(d *ds.User, _ int) int {
		return d.GitLabID
	}))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch authored MRs")
	}

	authorToMR = make(map[int][]*ds.MergeRequest, len(authoredMRs))
	for _, mr := range authoredMRs {
		if n.policy.IsApproved(n.team, mr) {
			continue
		}

		authorToMR[mr.Author.GitLabID] = append(authorToMR[mr.Author.GitLabID], mr)
	}

	toReviewMRs, err := n.repository.MergeRequestsByReviewer(lo.Map(devs, func(d *ds.User, _ int) int {
		return d.GitLabID
	}))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch MRs by reviewer")
	}

	reviewerToMR = make(map[int][]*ds.MergeRequest, len(toReviewMRs))
	for _, mr := range toReviewMRs {
		if n.policy.IsApproved(n.team, mr) {
			continue
		}

		for _, reviewer := range mr.Reviewers {
			reviewerToMR[reviewer.GitLabID] = append(reviewerToMR[reviewer.GitLabID], mr)
		}
	}

	return nil, nil, err
}

func (n *Notifications) slackMessages(
	l zerolog.Logger,
	devs []*ds.User,
	authorToMR, reviewerToMR map[int][]*ds.MergeRequest,
) ([]SlackMessage, error) {
	userTemplate, err := template.New("user_notification").Parse(n.team.Notifications.UserTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse user template")
	}

	channelTemplate, err := template.New("channel_notification").Parse(n.team.Notifications.ChannelTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse channel template")
	}

	summary := ds.ChannelNotification{
		AverageCount:  0,
		TotalCount:    0,
		LastEditedMR:  nil,
		FirstEditedMR: nil,
	}
	uniqMR := make(map[int]bool, len(authorToMR)+len(reviewerToMR))

	slackMessages := make([]SlackMessage, 0, len(devs)+1)

	for _, dev := range devs {
		if dev.SlackID == "" {
			l.Warn().Str("user", dev.Name).Msg("skipping user without Slack ID")
			continue
		}

		msg := bytes.NewBufferString("")

		authorToMRs := authorToMR[dev.GitLabID]
		for _, mr := range authorToMRs {
			uniqMR[mr.ID] = true
			if summary.LastEditedMR == nil {
				summary.LastEditedMR = mr
				summary.FirstEditedMR = mr
			} else if mr.UpdatedAt.Before(*summary.FirstEditedMR.UpdatedAt) {
				summary.FirstEditedMR = mr
			} else if mr.UpdatedAt.After(*summary.LastEditedMR.UpdatedAt) {
				summary.LastEditedMR = mr
			}
		}

		err = userTemplate.Execute(msg, ds.UserNotification{
			AuthoredMR: authorToMRs,
			ReviewerMR: reviewerToMR[dev.GitLabID],
		})
		if err != nil {
			l.Error().Err(err).Str("user", dev.Name).Msg("failed to execute user notification template")
			continue
		}

		slackMessages = append(slackMessages, SlackMessage{
			RecipientID: dev.SlackID,
			Text:        msg.String(),
		})
	}

	summary.TotalCount = len(uniqMR)
	summary.AverageCount = int(math.Round(float64(summary.TotalCount) / float64(len(devs))))

	chanMsg := bytes.NewBufferString("")

	err = channelTemplate.Execute(chanMsg, summary)
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute channel notification template")
	}

	slackMessages = append(slackMessages, SlackMessage{
		RecipientID: n.team.Notifications.ChannelID,
		Text:        chanMsg.String(),
	})

	return slackMessages, nil
}
