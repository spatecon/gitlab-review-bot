//go:build mongodb

package repository

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func TestRepository_UserBySlackID(t *testing.T) {
	rep := repositoryHelper(t)

	t.Run("empty collection", func(t *testing.T) {
		teams, err := rep.Teams()
		require.NoError(t, err, "failed to get teams")
		require.Len(t, teams, 0, "teams should be not found")
	})

	excepted := &ds.Team{
		ID:   "100",
		Name: "test team",
		Members: []*ds.User{
			{
				BasicUser: &ds.BasicUser{
					Name:     "test name",
					GitLabID: 333,
				},
				SlackID: "BBXXXXXBB",
				Labels:  ds.UserLabels{ds.DeveloperLabel},
			},
		},
		Policy: "test policy",
		Notifications: ds.NotificationSettings{
			Enabled:         true,
			Cron:            "* * * * *",
			UserTemplate:    "test user template {{.User.Name}}",
			ChannelID:       "DDXXXXXDD",
			ChannelTemplate: "test channel template}",
		},
	}

	t.Run("add team", func(t *testing.T) {
		res, err := rep.teams.InsertOne(rep.ctx, excepted)

		require.NoError(t, err, "failed to insert team")
		require.Equal(t, res.InsertedID, "100")
	})

	t.Run("should return created merge request", func(t *testing.T) {
		user, team, err := rep.UserBySlackID("BBXXXXXBB")
		require.NoError(t, err, "failed to get teammate")
		require.EqualValues(t, excepted, team, "teams should be equal")
		require.EqualValues(t, excepted.Members[0], user, "users should be equal")
	})
}
