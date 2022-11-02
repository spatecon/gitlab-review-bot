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
	ID           int          `bson:"id"`
	IID          int          `bson:"iid"`
	ProjectID    int          `bson:"project_id"`
	TargetBranch string       `bson:"target_branch"`
	SourceBranch string       `bson:"source_branch"`
	Title        string       `bson:"title"`
	Description  string       `bson:"description"`
	State        State        `bson:"state"`
	Author       *BasicUser   `bson:"author"`
	Assignees    []*BasicUser `bson:"assignees"`
	Reviewers    []*BasicUser `bson:"reviewers"`
	Draft        bool         `bson:"draft"`
	SHA          string       `bson:"sha"`
	UpdatedAt    *time.Time   `bson:"updated_at"`
	CreatedAt    *time.Time   `bson:"created_at"`

	// Additional information
	Approves []*BasicUser `bson:"approves"`
}

// IsEqual checks if two merge requests are equal (according to basic information)
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
