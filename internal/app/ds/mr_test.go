package ds

import (
	"testing"
	"time"
)

func TestMergeRequest_IsEqual(t *testing.T) {
	fakeTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	fakeTimeInUTC4 := fakeTime.In(time.FixedZone("UTC+4", 4*60*60))
	anotherFakeTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		a    *MergeRequest
		b    *MergeRequest
		want bool
	}{
		{
			name: "nil and nil",
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			name: "nil and not nil",
			a:    nil,
			b:    &MergeRequest{},
			want: false,
		},
		{
			name: "not nil and nil",
			a:    &MergeRequest{},
			b:    nil,
			want: false,
		},
		{
			name: "not nil and not nil",
			a:    &MergeRequest{},
			b:    &MergeRequest{},
			want: true,
		},
		{
			name: "not nil and not nil",
			a:    &MergeRequest{Title: "title"},
			b:    &MergeRequest{Title: "title"},
			want: true,
		},
		{
			name: "different titles",
			a:    &MergeRequest{Title: "title"},
			b:    &MergeRequest{Title: "title2"},
			want: false,
		},
		{
			name: "different assignees",
			a:    &MergeRequest{Assignees: []*BasicUser{{GitLabID: 123}}},
			b:    &MergeRequest{Assignees: []*BasicUser{{GitLabID: 676}}},
			want: false,
		},
		{
			name: "fuzz all values",
			a: &MergeRequest{
				ID:           123,
				IID:          321,
				ProjectID:    444,
				TargetBranch: "target",
				SourceBranch: "source",
				Title:        "title",
				Description:  "description",
				State:        "state",
				Assignees:    []*BasicUser{{GitLabID: 123}},
				Reviewers:    []*BasicUser{{GitLabID: 888}, {GitLabID: 555}},
				Draft:        true,
				SHA:          "sha",
				UpdatedAt:    &fakeTime,
				CreatedAt:    &fakeTime,
			},
			b: &MergeRequest{
				ID:           123,
				IID:          321,
				ProjectID:    444,
				TargetBranch: "target",
				SourceBranch: "source",
				Title:        "title",
				Description:  "description",
				State:        "state",
				Assignees:    []*BasicUser{{GitLabID: 123}},
				Reviewers:    []*BasicUser{{GitLabID: 555}, {GitLabID: 888}},
				Draft:        true,
				SHA:          "sha",
				UpdatedAt:    &fakeTime,
				CreatedAt:    &fakeTime,
			},
			want: true,
		},
		{
			name: "different updated at",
			a: &MergeRequest{
				ID:        123,
				UpdatedAt: &fakeTime,
			},
			b: &MergeRequest{
				ID:        123,
				UpdatedAt: &anotherFakeTime,
			},
			want: false,
		},
		{
			name: "same updated at, but different timezone",
			a: &MergeRequest{
				ID:        123,
				UpdatedAt: &fakeTime,
			},
			b: &MergeRequest{
				ID:        123,
				UpdatedAt: &fakeTimeInUTC4,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.IsEqual(tt.b); got != tt.want {
				t.Errorf("MergeRequest.IsEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
