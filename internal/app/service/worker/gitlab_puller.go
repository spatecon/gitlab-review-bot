package worker

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

const pullPeriod = 5 * time.Minute

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
		once := time.NewTimer(5 * time.Second)

		for {
			select {
			case <-once.C:
				g.pullAndHandle()
			case <-ticker.C:
				g.pullAndHandle()
			case <-g.close:
				once.Stop()
				ticker.Stop()
				return
			}
		}
	}()
}

func (g *GitLabPuller) pullAndHandle() {
	mrs, err := g.gitlab.MergeRequestsByProject(g.projectID)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch merge requests")
	}

	log.Info().Int("project_id", g.projectID).
		Int("count", len(mrs)).
		Msg("fetched merge requests")

	for _, mr := range mrs {
		err = g.handler(mr)
		if err != nil {
			log.Error().Err(err).Msg("failed to handle merge requests")
		}
	}
}

func (g *GitLabPuller) Close() {
	g.close <- struct{}{}
}
