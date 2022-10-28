package ds

import "time"

type State string

const (
	StateOpened State = "opened"
	StateClosed State = "closed"
	StateLocked State = "locked"
	StateMerged State = "merged"
)

type MergeRequest struct {
	ID           int
	IID          int
	ProjectID    int
	TargetBranch string
	SourceBranch string
	Title        string
	Description  string
	State        State
	Author       *BasicUser
	Assignees    []*BasicUser
	Reviewers    []*BasicUser
	Draft        bool
	SHA          string
	UpdatedAt    *time.Time
	CreatedAt    *time.Time

	//Discussion   Discussion // TODO: fetch discussion
}

func (a *MergeRequest) IsEqual(b *MergeRequest) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if a.ID != b.ID {
		return false
	}

	if a.IID != b.IID {
		return false
	}

	if a.ProjectID != b.ProjectID {
		return false
	}

	if a.TargetBranch != b.TargetBranch {
		return false
	}

	if a.SourceBranch != b.SourceBranch {
		return false
	}

	if a.Title != b.Title {
		return false
	}

	if a.Description != b.Description {
		return false
	}

	if a.State != b.State {
		return false
	}

	if !EqualUsers(a.Assignees, b.Assignees) {
		return false
	}

	if !EqualUsers(a.Reviewers, b.Reviewers) {
		return false
	}

	if a.Draft != b.Draft {
		return false
	}

	if a.SHA != b.SHA {
		return false
	}

	if (a.UpdatedAt != nil && b.UpdatedAt != nil && !a.UpdatedAt.Equal(*b.UpdatedAt)) || (a.UpdatedAt == nil && b.UpdatedAt != nil) || (a.UpdatedAt != nil && b.UpdatedAt == nil) {
		return false
	}

	if (a.CreatedAt != nil && b.CreatedAt != nil && !a.CreatedAt.Equal(*b.CreatedAt)) || (a.CreatedAt == nil && b.CreatedAt != nil) || (a.CreatedAt != nil && b.CreatedAt == nil) {
		return false
	}

	return true
}
