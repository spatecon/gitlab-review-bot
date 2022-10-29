package basic

import (
	"strings"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
	"github.com/spatecon/gitlab-review-bot/internal/pkg/fsm"
)

const policy = "BASIC"

type StateMachine struct {
	team *ds.Team
}

func New(team *ds.Team) (*StateMachine, error) {
	return &StateMachine{
		team: team,
	}, nil
}

func (s *StateMachine) Policy() string {
	return policy
}

func (s *StateMachine) Team() *ds.Team {
	return s.team
}

func (s *StateMachine) isTeamAuthored(mr *ds.MergeRequest) bool {
	for _, member := range s.team.Members {
		if member.BasicUser.GitLabID == mr.Author.GitLabID {
			return true
		}
	}

	return false
}

func (s *StateMachine) State(mr *ds.MergeRequest) fsm.State {
	if !s.isTeamAuthored(mr) {
		return SkippedState
	}

	if mr.SourceBranch == "master" && strings.HasPrefix(mr.SourceBranch, "release/") {
		return SkippedState
	}

	if mr.Draft {
		return DraftState
	}

	if mr.State == ds.StateMerged {
		return MergedState
	}

	if mr.State == ds.StateClosed {
		return ClosedState
	}

	if mr.State == ds.StateLocked {
		return SkippedState
	}

	return InReviewState
}
