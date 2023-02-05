package service

import (
	"bytes"
	"math"
	"text/template"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
	"github.com/spatecon/gitlab-review-bot/pkg/templating"
)

func (s *Service) GetAuthoredReviewedMRs(team *ds.Team, users []*ds.User) (authorToMR, reviewerToMR map[int][]*ds.MergeRequest, err error) {
	policy, ok := s.policies[team.Policy]
	if !ok {
		return nil, nil, errors.Errorf("policy %s not found", team.Policy)
	}

	authoredMRs, err := s.r.MergeRequestsByAuthor(lo.Map(users, func(d *ds.User, _ int) int {
		return d.GitLabID
	}))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch authored MRs")
	}

	authorToMR = make(map[int][]*ds.MergeRequest, len(authoredMRs))
	for _, mr := range authoredMRs {
		if mr.Author == nil {
			continue
		}

		if policy.ApprovedByPolicy(team, mr) {
			continue
		}

		authorToMR[mr.Author.GitLabID] = append(authorToMR[mr.Author.GitLabID], mr)
	}

	toReviewMRs, err := s.r.MergeRequestsByReviewer(lo.Map(users, func(d *ds.User, _ int) int {
		return d.GitLabID
	}))
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to fetch MRs by reviewer")
	}

	reviewerToMR = make(map[int][]*ds.MergeRequest, len(toReviewMRs))
	for _, mr := range toReviewMRs {
		for _, reviewer := range mr.Reviewers {
			if policy.ApprovedByUser(team, mr, reviewer) {
				continue
			}

			reviewerToMR[reviewer.GitLabID] = append(reviewerToMR[reviewer.GitLabID], mr)
		}
	}

	return authorToMR, reviewerToMR, nil
}

func (s *Service) UserNotification(user *ds.User, team *ds.Team, authorToMR, reviewerToMR map[int][]*ds.MergeRequest) (message string, err error) {
	// TODO: optimize initializations for performance
	userTemplate := template.New("user_notification").Funcs(s.templateFuncMap(team.Notifications.Locale))

	userTemplate, err = userTemplate.Parse(team.Notifications.UserTemplate)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse user template")
	}

	msg := bytes.NewBufferString("")

	err = userTemplate.Execute(msg, ds.UserNotification{
		User:       user,
		AuthoredMR: authorToMR[user.GitLabID],
		ReviewerMR: reviewerToMR[user.GitLabID],
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to execute user notification template")
	}

	return msg.String(), nil
}

func (s *Service) TeamNotification(team *ds.Team, authorToMR, reviewerToMR map[int][]*ds.MergeRequest) (message string, err error) {
	channelTemplate := template.New("team_notification").Funcs(s.templateFuncMap(team.Notifications.Locale))

	channelTemplate, err = channelTemplate.Parse(team.Notifications.ChannelTemplate)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse team template")
	}

	summary := ds.ChannelNotification{
		Team:          team,
		AverageCount:  0,
		TotalCount:    0,
		LastEditedMR:  nil,
		FirstEditedMR: nil,
	}
	uniqMR := make(map[int]bool, len(authorToMR)+len(reviewerToMR))

	for _, member := range team.Members {
		if member.SlackID == "" {
			log.Warn().
				Str("user", member.Name).
				Msg("skipping user without Slack ID")
			continue
		}

		authorToMRs := authorToMR[member.GitLabID]
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
	}

	summary.TotalCount = len(uniqMR)
	summary.AverageCount = int(math.Round(float64(summary.TotalCount) / float64(len(team.Members))))

	chanMsg := bytes.NewBufferString("")

	err = channelTemplate.Execute(chanMsg, summary)
	if err != nil {
		return "", errors.Wrap(err, "failed to execute channel notification template")
	}

	return chanMsg.String(), nil
}

func (s *Service) templateFuncMap(locale string) template.FuncMap {
	loc, ok := templating.ParseLocale(locale)

	if !ok {
		log.Warn().Str("locale", locale).Msg("failed to parse locale, using default (en_EN)")
	}

	tools := templating.NewTools(loc)

	return template.FuncMap{
		"since":  tools.Since,
		"plural": tools.Plural,
	}
}
