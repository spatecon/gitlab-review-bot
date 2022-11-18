package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func TestRepository_MergeRequests(t *testing.T) {
	rep := repositoryHelper(t)

	t.Run("should return empty list", func(t *testing.T) {
		mr, err := rep.MergeRequestByID(1)
		require.NoError(t, err, "failed to get merge requests")
		require.Nil(t, mr, "merge requests should be not found")
	})

	ts := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	mr1 := &ds.MergeRequest{
		ID:           1,
		IID:          2,
		ProjectID:    3,
		TargetBranch: "4",
		SourceBranch: "5",
		Title:        "6",
		Description:  "7",
		State:        "8",
		Author:       &ds.BasicUser{GitLabID: 9, Name: "10"},
		Assignees:    []*ds.BasicUser{{GitLabID: 11, Name: "12"}},
		Reviewers:    []*ds.BasicUser{{GitLabID: 13, Name: "14"}},
		Draft:        true,
		SHA:          "15",
		URL:          "16",
		UpdatedAt:    &ts,
		CreatedAt:    &ts,
		Approves:     []*ds.BasicUser{{GitLabID: 17, Name: "18"}},
	}

	t.Run("creates a merge request", func(t *testing.T) {
		err := rep.UpsertMergeRequest(mr1)
		require.NoError(t, err, "failed to create merge request")
	})

	t.Run("should return created merge request", func(t *testing.T) {
		mr, err := rep.MergeRequestByID(1)
		require.NoError(t, err, "failed to get merge requests")
		require.EqualValues(t, mr1, mr, "merge requests should be equal")
	})
}
