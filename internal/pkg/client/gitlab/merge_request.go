package gitlab

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/xanzy/go-gitlab"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

const (
	perPage  = 100 // 100 is api max
	maxPages = 10
)

// MergeRequestsByProject only last 1000 merge requests are processed
func (c *Client) MergeRequestsByProject(projectID int, createdAfter time.Time) ([]*ds.MergeRequest, error) {
	// TODO: consider using webhooks

	allMergeRequests := make([]*ds.MergeRequest, 0, perPage)

	for i := 1; i <= maxPages; i++ {
		log.Trace().Msg("fetching merge requests")
		// docs: https://docs.gitlab.com/ee/api/merge_requests.html#list-project-merge-requests
		mergeRequests, resp, err := c.gitlab.MergeRequests.ListProjectMergeRequests(
			projectID,
			&gitlab.ListProjectMergeRequestsOptions{
				CreatedAfter: &createdAfter,
				ListOptions: gitlab.ListOptions{
					Page:    i,
					PerPage: perPage,
				},
			},
			gitlab.WithContext(c.ctx))
		if err != nil {
			return nil, errors.Wrap(err, "error getting merge requests")
		}

		for _, mergeRequest := range mergeRequests {
			allMergeRequests = append(allMergeRequests, mergeRequestConvert(mergeRequest))
		}

		if resp.TotalPages <= i {
			break
		}
	}

	return allMergeRequests, nil
}

func mergeRequestConvert(req *gitlab.MergeRequest) *ds.MergeRequest {
	var author *ds.BasicUser
	if req.Author != nil {
		author = &ds.BasicUser{
			Name:     req.Author.Name,
			GitLabID: req.Author.ID,
		}
	}

	assignees := make([]*ds.BasicUser, 0, len(req.Assignees))
	for _, assignee := range req.Assignees {
		assignees = append(assignees, &ds.BasicUser{
			Name:     assignee.Name,
			GitLabID: assignee.ID,
		})
	}

	reviewers := make([]*ds.BasicUser, 0, len(req.Reviewers))
	for _, reviewer := range req.Reviewers {
		reviewers = append(reviewers, &ds.BasicUser{
			Name:     reviewer.Name,
			GitLabID: reviewer.ID,
		})
	}

	return &ds.MergeRequest{
		ID:           req.ID,
		IID:          req.IID,
		ProjectID:    req.ProjectID,
		TargetBranch: req.TargetBranch,
		SourceBranch: req.SourceBranch,
		Title:        req.Title,
		Description:  req.Description,
		State:        ds.State(req.State),
		Author:       author,
		Assignees:    assignees,
		Reviewers:    reviewers,
		Draft:        req.Draft,
		SHA:          req.SHA,
		URL:          req.WebURL,
		UpdatedAt:    req.UpdatedAt,
		CreatedAt:    req.CreatedAt,
	}
}
