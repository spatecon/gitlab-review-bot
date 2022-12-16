package gitlab

import (
	"github.com/pkg/errors"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func (c *Client) MergeRequestApproves(projectID int, iid int) ([]*ds.BasicUser, error) {
	c.rl.Take()
	approvals, _, err := c.gitlab.MergeRequestApprovals.GetConfiguration(projectID, iid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get merge request approvals")
	}

	approvedUsers := make([]*ds.BasicUser, 0, len(approvals.ApprovedBy))

	for _, approvedBy := range approvals.ApprovedBy {
		approvedUsers = append(approvedUsers, &ds.BasicUser{
			Name:     approvedBy.User.Name,
			GitLabID: approvedBy.User.ID,
		})
	}

	return approvedUsers, nil
}
