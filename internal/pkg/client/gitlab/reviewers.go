package gitlab

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/xanzy/go-gitlab"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func (c *Client) SetReviewers(mr *ds.MergeRequest, reviewers []int) error {
	log.Info().Int("project_id", mr.ProjectID).Int("mr_id", mr.IID).Ints("reviewers", reviewers).Msg("reviewers set")

	_, _, err := c.gitlab.MergeRequests.UpdateMergeRequest(mr.ProjectID, mr.IID, &gitlab.UpdateMergeRequestOptions{
		Description: gitlab.String(mr.Description),
		ReviewerIDs: &reviewers,
	})
	if err != nil {
		return errors.Wrap(err, "error calling gitlab apid to set reviewers")
	}

	return nil
}
