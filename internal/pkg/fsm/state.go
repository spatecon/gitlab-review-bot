package fsm

import "github.com/spatecon/gitlab-review-bot/internal/app/ds"

type StateMachine interface {
	Policy() string
	Team() *ds.Team
	State(request *ds.MergeRequest) State
}

type State interface {
	Name() string
	IsSkipped() bool
	IsFinal() bool
	IsApproved() bool
}
