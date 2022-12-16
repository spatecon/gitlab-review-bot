package worker

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

const (
	gitlabPullerWorkerName = "gitlab_puller_worker"
)

type GitlabClient interface {
	MergeRequestsByProject(projectID int, createdAfter time.Time) ([]*ds.MergeRequest, error)
}

type MergeRequestHandler func(mr *ds.MergeRequest) error

type GitLabPuller struct {
	gitlab     GitlabClient
	handler    MergeRequestHandler
	projectID  int
	pullPeriod time.Duration
	close      chan struct{}
	after      time.Time
}

func NewGitLabPuller(pullPeriod time.Duration, after time.Time, gitlab GitlabClient, handler MergeRequestHandler, projectID int) (*GitLabPuller, error) {
	worker := &GitLabPuller{
		gitlab:     gitlab,
		handler:    handler,
		projectID:  projectID,
		pullPeriod: pullPeriod,
		after:      after,
		close:      make(chan struct{}),
	}

	return worker, nil
}

func (g *GitLabPuller) Run() {
	go func() {
		ticker := time.NewTicker(g.pullPeriod)
		startup := time.NewTimer(5 * time.Second)

		for {
			select {
			case <-startup.C:
				g.pullAndHandle()
			case <-ticker.C:
				g.pullAndHandle()
			case <-g.close:
				startup.Stop()
				ticker.Stop()
				return
			}
		}
	}()
}

func (g *GitLabPuller) pullAndHandle() {
	l := log.With().
		Str("worker", gitlabPullerWorkerName).
		Int("project_id", g.projectID).
		Logger()

	l.Info().Msg("pulling merge requests")

	mrs, err := g.gitlab.MergeRequestsByProject(g.projectID, g.after)
	if err != nil {
		l.Error().Err(err).Msg("failed to fetch merge requests")
	}

	l.Info().Int("project_id", g.projectID).
		Int("count", len(mrs)).
		Msg("pulled merge requests successfully")

	for _, mr := range mrs {
		err = g.handler(mr)
		if err != nil {
			l.Error().Err(err).Msg("failed to handle merge requests")
		}
	}

	log.Info().Int("project_id", g.projectID).Msg("merge requests handled")
}

func (g *GitLabPuller) Close() {
	g.close <- struct{}{}
}
