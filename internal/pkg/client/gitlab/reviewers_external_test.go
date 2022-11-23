//go:build gitlab

package gitlab

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func TestClient_SetReviewers(t *testing.T) {
	type args struct {
		mr        *ds.MergeRequest
		reviewers []int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "set reviewers",
			args: args{
				mr: &ds.MergeRequest{
					IID:       1740,
					ProjectID: 14415894,
				},
				reviewers: []int{10592575, 11239710},
			},
			wantErr: false,
		},
	}

	token := os.Getenv("GITLAB_TOKEN")

	if token == "" {
		t.Skip("EXTERNAL TEST SKIPPED: env GITLAB_TOKEN is not set")
	}

	client, err := gitlab.NewClient(token)
	require.NoError(t, err, "error creating gitlab client")

	c := &Client{
		ctx:    context.TODO(),
		gitlab: client,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = c.SetReviewers(tt.args.mr, tt.args.reviewers)
			require.Equal(t, tt.wantErr, err != nil, "error expected or not")
		})
	}
}
