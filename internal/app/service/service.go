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
	UpsertMergeRequest(mr *ds.MergeRequest) error
}

type GitlabClient interface {
	MergeRequestsByProject(projectID int) ([]*ds.MergeRequest, error)
	MergeRequestApproves(projectID int, iid int) ([]*ds.BasicUser, error)
}

type Worker interface {
	Close()
}

type Service struct {
	repo   Repository
	gitlab GitlabClient

	workers []Worker
}

func New(repo Repository, gitlab GitlabClient) (*Service, error) {
	return &Service{
		repo:   repo,
		gitlab: gitlab,
	}, nil
}

func (s *Service) Close() error {
	for _, wrk := range s.workers {
		wrk.Close()
	}

	return nil
}

func (s *Service) SubscribeOnProjects() error {
	projects, err := s.repo.Projects()
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

func (s *Service) mergeRequestsHandler(actual *ds.MergeRequest) error {
	// fetch MR from repository
	old, err := s.repo.MergeRequestByID(actual.ID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch merge request from repository")
	}

	// if no changes, do nothing
	if old != nil && old.IsEqual(actual) {
		log.Info().Int("id", actual.ID).Msg("mr skipped")
		return nil
	}

	// enrich MR with approves
	approves, err := s.gitlab.MergeRequestApproves(actual.ProjectID, actual.IID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch merge request approves")
	}

	// TODO: consider m := m1.Merge(m2 MR) method
	actual.Approves = approves
	if old != nil {
		actual.ReviewersByBot = old.ReviewersByBot
	}

	// update (or create) it
	err = s.repo.UpsertMergeRequest(actual)
	if err != nil {
		return errors.Wrap(err, "failed to update merge request in repository")
	}

	log.Info().Int("id", actual.ID).Msg("mr updated or created")

	// TODO: launch pipeline

	return nil
}
