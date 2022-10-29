package service

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
	"github.com/spatecon/gitlab-review-bot/internal/app/service/worker"
	"github.com/spatecon/gitlab-review-bot/internal/pkg/fsm"
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
	repo          Repository
	gitlab        GitlabClient
	teams         []*ds.Team
	stateMachines []fsm.StateMachine

	workers []Worker
}

func New(repo Repository, gitlab GitlabClient) (*Service, error) {
	svc := &Service{
		repo:   repo,
		gitlab: gitlab,
	}

	err := svc.loadTeams()
	if err != nil {
		return nil, errors.Wrap(err, "failed to pre-cache teams")
	}

	err = svc.initStateMachines()
	if err != nil {
		return nil, errors.Wrap(err, "failed to init state machines")
	}

	return svc, nil
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

func (s *Service) mergeRequestsHandler(mr *ds.MergeRequest) error {
	// fetch MR from repository
	old, err := s.repo.MergeRequestByID(mr.ID)
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

	// TODO: consider m := m1.Merge(m2 MR) method
	mr.Approves = approves
	if old != nil {
		mr.ReviewersByBot = old.ReviewersByBot
	}

	// update (or create) it
	err = s.repo.UpsertMergeRequest(mr)
	if err != nil {
		return errors.Wrap(err, "failed to update merge request in repository")
	}

	log.Info().Int("id", mr.ID).Msg("mr updated or created")

	// launch pipeline on each state machine
	for _, sm := range s.stateMachines {
		state := sm.State(mr)
		log.Info().
			Int("id", mr.IID).
			Str("state", state.Name()).
			Str("policy", sm.Policy()).
			Bool("is_final", state.IsFinal()).
			Bool("is_approved", state.IsApproved()).
			Msg("state machine")
	}

	return nil
}
