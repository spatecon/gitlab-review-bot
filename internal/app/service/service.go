package service

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
	"github.com/spatecon/gitlab-review-bot/internal/app/service/worker"
)

type Repository interface {
	Teams() ([]*ds.Team, error)
	Projects() ([]*ds.Project, error)
	MergeRequestByID(id int) (*ds.MergeRequest, error)
	MergeRequestsByProject(projectID int) ([]*ds.MergeRequest, error)
	MergeRequestsByAuthor(authorID []int) ([]*ds.MergeRequest, error)
	MergeRequestsByReviewer(reviewerID []int) ([]*ds.MergeRequest, error)
	UpsertMergeRequest(mr *ds.MergeRequest) error
}

type GitlabClient interface {
	MergeRequestsByProject(projectID int) ([]*ds.MergeRequest, error)
	MergeRequestApproves(projectID int, iid int) ([]*ds.BasicUser, error)
}

type SlackClient interface {
	worker.SlackClient
}

type Worker interface {
	Close()
}

type Policy interface {
	ProcessChanges(team *ds.Team, mr *ds.MergeRequest) (err error)
	IsApproved(team *ds.Team, mr *ds.MergeRequest) bool
}

type Service struct {
	r        Repository
	gitlab   GitlabClient
	slack    SlackClient
	teams    []*ds.Team
	policies map[ds.PolicyName]Policy

	workers []Worker
}

func New(r Repository, g GitlabClient, p map[ds.PolicyName]Policy) (*Service, error) {
	svc := &Service{
		r:        r,
		gitlab:   g,
		policies: p,
	}

	// TODO: team hot reload
	err := svc.loadTeams()
	if err != nil {
		return nil, errors.Wrap(err, "failed to pre-cache teams")
	}

	err = svc.initNotifications()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init notifications")
	}

	return svc, nil
}

func (s *Service) Close() error {
	for _, wrk := range s.workers {
		wrk.Close()
	}

	return nil
}

// SubscribeOnProjects Creates workers for each project and subscribe on merge requests changes
func (s *Service) SubscribeOnProjects() error {
	projects, err := s.r.Projects()
	if err != nil {
		return err
	}

	for _, project := range projects {
		var wrk Worker

		wrk, err = worker.NewGitLabPuller(s.gitlab, s.mergeRequestsHandler, project.ID)
		if err != nil {
			return errors.Wrap(err, "failed to create gitlab puller")
		}

		s.workers = append(s.workers, wrk)
	}

	return nil
}

func (s *Service) mergeRequestsHandler(mr *ds.MergeRequest) error {
	// fetch MR from repository
	old, err := s.r.MergeRequestByID(mr.ID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch merge request from repository")
	}

	// if no changes, do nothing
	if old != nil && old.IsEqual(mr) {
		log.Info().Int("id", mr.ID).Msg("mr skipped")
		return nil
	}

	// enrich MR with approves
	approves, err := s.gitlab.MergeRequestApproves(mr.ProjectID, mr.IID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch merge request approves")
	}

	mr.Approves = approves

	// update (or create) it
	err = s.r.UpsertMergeRequest(mr)
	if err != nil {
		return errors.Wrap(err, "failed to update merge request in repository")
	}

	log.Info().Int("id", mr.ID).Msg("mr updated or created")

	// process MR

	for _, team := range s.teams {
		policy, ok := s.policies[team.Policy]
		if !ok {
			log.Error().
				Str("team", team.Name).
				Str("policy", string(team.Policy)).
				Msg("failed to process updates unknown policy")
			continue
		}

		err = policy.ProcessChanges(team, mr)
		if err != nil {
			log.Error().
				Err(err).
				Str("team", team.Name).
				Str("policy", string(team.Policy)).
				Msg("failed to process merge request")
		}

	}

	return nil
}
