package gitlab

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/xanzy/go-gitlab"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func (c *Client) SetReviewers(mr *ds.MergeRequest, reviewers []int) error {
	log.Info().Int("project_id", mr.ProjectID).Int("mr_id", mr.IID).Ints("reviewers", reviewers).Msg("reviewers set")

	actual, _, err := c.gitlab.MergeRequests.UpdateMergeRequest(mr.ProjectID, mr.IID, &gitlab.UpdateMergeRequestOptions{
		ReviewerIDs: &reviewers,
	})
	if err != nil {
		return errors.Wrap(err, "error calling gitlab apid to set reviewers")
	}

	needed := make(map[int]bool, len(reviewers))

	for _, basicUser := range actual.Reviewers {
		needed[basicUser.ID] = true
	}

	for _, gitlabID := range reviewers {
		if _, ok := needed[gitlabID]; !ok {
			return errors.Errorf("reviewer not set: %d", gitlabID)
		}
	}

	return nil
}
