package gitlab

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func TestMergeRequestConvert(t *testing.T) {
	fakeDate := time.Date(2022, 10, 28, 12, 0, 0, 0, time.UTC)

	tc := []struct {
		name string
		in   *gitlab.MergeRequest
		out  *ds.MergeRequest
	}{
		{
			name: "full fields",
			in: &gitlab.MergeRequest{
				ID:           155016530,
				IID:          133,
				TargetBranch: "main",
				SourceBranch: "release/1.0.0",
				ProjectID:    15513260,
				Title:        "Title of MR",
				State:        "merged",
				CreatedAt:    &fakeDate,
				UpdatedAt:    &fakeDate,
				Author:       &gitlab.BasicUser{ID: 1, Name: "Morning Dave"},
				Assignee:     nil,
				Assignees:    []*gitlab.BasicUser{{ID: 2, Name: "Elijah Oak"}},
				Reviewers: []*gitlab.BasicUser{
					{ID: 40, Name: "Best Reviewer 40"},
					{ID: 50, Name: "Best Reviewer 50"},
				},
				Description: "Full description of MR\nwith multilines",
				Draft:       true,
				SHA:         "ced4a15efbf5e9db64f871bf2fc363c8873e525c",
			},
			out: &ds.MergeRequest{
				ID:           155016530,
				IID:          133,
				TargetBranch: "main",
				SourceBranch: "release/1.0.0",
				ProjectID:    15513260,
				Title:        "Title of MR",
				Description:  "Full description of MR\nwith multilines",
				State:        ds.StateMerged,
				CreatedAt:    &fakeDate,
				UpdatedAt:    &fakeDate,
				Author: &ds.BasicUser{
					Name:     "Morning Dave",
					GitLabID: 1,
				},
				Assignees: []*ds.BasicUser{
					{Name: "Elijah Oak", GitLabID: 2},
				},
				Reviewers: []*ds.BasicUser{
					{Name: "Best Reviewer 40", GitLabID: 40},
					{Name: "Best Reviewer 50", GitLabID: 50},
				},
				Draft: true,
				SHA:   "ced4a15efbf5e9db64f871bf2fc363c8873e525c",
			},
		},
		{
			name: "some fields are empty",
			in: &gitlab.MergeRequest{
				ID:        12312312,
				Author:    nil,
				UpdatedAt: nil,
				CreatedAt: &fakeDate,
			},
			out: &ds.MergeRequest{
				ID:        12312312,
				Author:    nil,
				UpdatedAt: nil,
				CreatedAt: &fakeDate,
				Assignees: []*ds.BasicUser{},
				Reviewers: []*ds.BasicUser{},
			},
		},
	}

	for _, cs := range tc {
		t.Run(cs.name, func(t *testing.T) {
			actual := mergeRequestConvert(cs.in)
			require.Equal(t, cs.out, actual)
		})
	}
}
