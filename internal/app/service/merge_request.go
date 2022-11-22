package service

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func (s *Service) mergeRequestsHandler(mr *ds.MergeRequest) error {
	// fetch MR from repository
	old, err := s.r.MergeRequestByID(mr.ID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch merge request from repository")
	}

	// if no changes, do nothing
	if old != nil && old.IsEqual(mr) {
		log.Info().
			Int("project_id", mr.ProjectID).
			Int("iid", mr.IID).
			Msg("mr skipped")
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

	log.Info().
		Int("project_id", mr.ProjectID).
		Int("iid", mr.IID).
		Str("url", mr.URL).
		Msg("mr updated or created")

	// process MR
	for _, team := range s.teams {

		if mr.CreatedAt != nil && mr.CreatedAt.Before(team.CreatedAt) {
			log.Info().Str("team_id", team.ID).Msg("skip team, mr created before team")
			continue
		}

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
