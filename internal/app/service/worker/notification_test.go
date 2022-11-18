package worker_test

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
	"github.com/spatecon/gitlab-review-bot/internal/app/service/worker"
	"github.com/spatecon/gitlab-review-bot/internal/app/service/worker/mocks"
	"github.com/spatecon/gitlab-review-bot/pkg/testloggger"
)

var (
	// John is a developer from the team
	John = &ds.BasicUser{
		Name:     "John Snow",
		GitLabID: 12345,
	}
	// Gordon is a lead and a developer from the team
	Gordon = &ds.BasicUser{
		Name:     "Gordon Freeman",
		GitLabID: 99991,
	}
	// Jane is a developer from another team
	Jane = &ds.BasicUser{
		Name:     "Jane Doe",
		GitLabID: 54321,
	}
)

func TestNotifications_Run(t *testing.T) {
	ctrl := gomock.NewController(t)

	testloggger.Set(t)
	defer testloggger.Unset()

	ts := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	// Define message templates
	userTemplate := `{{.User.Name}} has {{.ReviewerMR | len}} MRs waiting for review. Also, {{.User.Name}} has {{.AuthoredMR | len}} authored MRs in review.`

	chanTemplate := `There are {{.TotalCount}} MRs waiting for review. Members of {{.Team.Name}} have an average of {{.AverageCount}} MRs waiting for review.`

	/*
		+--------+-----+-----+-----+-----+
		|        | MR1 | MR2 | MR3 | MR4 |	A – Authored
		+--------+-----+-----+-----+-----+	R – Reviewed
		| John   | A   | A   | R   | R   |
		+--------+-----+-----+-----+-----+
		| Gordon |     | R   | A   | R   |
		+--------+-----+-----+-----+-----+
		| Jane   | R   | R   |     | A   |
		+--------+-----+-----+-----+-----+
	*/

	MR1 := &ds.MergeRequest{
		ID:        77770,
		IID:       1,
		Title:     "Test MR 1",
		URL:       "https://gitlab.com/test/test/-/merge_requests/1",
		Author:    John,
		Reviewers: []*ds.BasicUser{Jane},
		UpdatedAt: &ts,
		CreatedAt: &ts,
	}
	MR2 := &ds.MergeRequest{
		ID:        77771,
		IID:       2,
		Title:     "Test MR 2",
		URL:       "https://gitlab.com/test/test/-/merge_requests/2",
		Author:    John,
		Reviewers: []*ds.BasicUser{Jane, Gordon},
		UpdatedAt: &ts,
		CreatedAt: &ts,
	}
	MR3 := &ds.MergeRequest{
		ID:        77772,
		IID:       3,
		Title:     "Test MR 3",
		URL:       "https://gitlab.com/test/test/-/merge_requests/3",
		Author:    Gordon,
		Reviewers: []*ds.BasicUser{John},
		UpdatedAt: &ts,
		CreatedAt: &ts,
	}
	MR4 := &ds.MergeRequest{
		ID:        77774,
		IID:       4,
		Title:     "Test MR 4",
		URL:       "https://gitlab.com/test/test/-/merge_requests/4",
		Author:    Jane,
		Reviewers: []*ds.BasicUser{John, Gordon},
		UpdatedAt: &ts,
		CreatedAt: &ts,
	}

	team := &ds.Team{
		ID:   "123fd00000000000000",
		Name: "Test Team",
		Members: []*ds.User{
			{
				BasicUser: John,
				SlackID:   "XJFAAAAAAAAA",
				Labels:    ds.UserLabels{ds.DeveloperLabel},
			},
			{
				BasicUser: Gordon,
				SlackID:   "XJFBBBBBBBBB",
				Labels:    ds.UserLabels{ds.DeveloperLabel, ds.LeadLabel},
			},
		},
		Policy: "test_policy",
		Notifications: ds.NotificationSettings{
			Enabled:         true,
			Cron:            "0 0 0 * * *",
			UserTemplate:    userTemplate,
			ChannelID:       "XCCCCCCCCCCC",
			ChannelTemplate: chanTemplate,
		},
	}

	policy := mocks.NewNotificationPolicy(ctrl)
	policy.EXPECT().
		IsApproved(team, gomock.Any()).
		Return(false).
		MinTimes(1)

	repository := mocks.NewNotificationRepository(ctrl)
	repository.EXPECT().
		MergeRequestsByAuthor(gomock.InAnyOrder([]int{John.GitLabID, Gordon.GitLabID})).
		Return([]*ds.MergeRequest{MR1, MR2, MR3}, nil).
		Times(1)
	repository.EXPECT().
		MergeRequestsByReviewer(gomock.InAnyOrder([]int{John.GitLabID, Gordon.GitLabID})).
		Return([]*ds.MergeRequest{MR2, MR3, MR4}, nil).
		Times(1)

	slackClient := mocks.NewNotificationSlackClient(ctrl)

	exceptedJohnMessage := `John Snow has 2 MRs waiting for review. Also, John Snow has 2 authored MRs in review.`
	slackClient.EXPECT().
		SendMessage(gomock.Eq("XJFAAAAAAAAA"), gomock.Eq(exceptedJohnMessage)).
		Return(nil)

	exceptedGordonMessage := `Gordon Freeman has 2 MRs waiting for review. Also, Gordon Freeman has 1 authored MRs in review.`
	slackClient.EXPECT().
		SendMessage(gomock.Eq("XJFBBBBBBBBB"), gomock.Eq(exceptedGordonMessage)).
		Return(nil)

	exceptedChanMessage := `There are 3 MRs waiting for review. Members of Test Team have an average of 2 MRs waiting for review.`
	slackClient.EXPECT().
		SendMessage(gomock.Eq("XCCCCCCCCCCC"), gomock.Eq(exceptedChanMessage)).
		Return(nil)

	notificationsWorker := worker.NewNotificationsWorker(team, policy, repository, slackClient)
	notificationsWorker.Run()
}
