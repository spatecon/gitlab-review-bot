package ds

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
	Author       BasicUser
	Assignees    []BasicUser
	Reviewers    []BasicUser
	Draft        bool
	SHA          string
	//Discussion   Discussion // todo: fetch discussion
}
