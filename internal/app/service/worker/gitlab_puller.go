package worker

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

const pullPeriod = 2 * time.Minute

type GitlabClient interface {
	MergeRequestsByProject(projectID int) ([]*ds.MergeRequest, error)
}

type MergeRequestHandler func(mr *ds.MergeRequest) error

type GitLabPuller struct {
	gitlab     GitlabClient
	handler    MergeRequestHandler
	projectID  int
	pullPeriod time.Duration
	close      chan struct{}
}

func NewGitLabPuller(gitlab GitlabClient, handler MergeRequestHandler, projectID int) (*GitLabPuller, error) {
	worker := &GitLabPuller{
		gitlab:     gitlab,
		handler:    handler,
		projectID:  projectID,
		pullPeriod: pullPeriod,
		close:      make(chan struct{}),
	}

	worker.Start()

	return worker, nil
}

func (g *GitLabPuller) Start() {
	go func() {
		ticker := time.NewTicker(g.pullPeriod)

		select {
		case <-ticker.C:
			mrs, err := g.gitlab.MergeRequestsByProject(g.projectID)
			if err != nil {
				log.Error().Err(err).Msg("failed to fetch merge requests")
			}

			for _, mr := range mrs {
				err = g.handler(mr)
				if err != nil {
					log.Error().Err(err).Msg("failed to handle merge requests")
				}
			}
		case <-g.close:
			ticker.Stop()
			return
		}
	}()
}

func (g *GitLabPuller) Close() {
	g.close <- struct{}{}
}
