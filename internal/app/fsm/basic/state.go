package basic

type State struct {
	name     string
	final    bool
	skipped  bool
	approved bool // todo: should be calculated
}

func (s *State) Name() string {
	return policy + "." + s.name
}

func (s *State) IsFinal() bool {
	return s.final
}

func (s *State) IsSkipped() bool {
	return s.skipped
}

func (s *State) IsApproved() bool {
	return s.approved
}

var (
	DraftState    = &State{"DRAFT", false, true, false}
	InReviewState = &State{"IN_REVIEW", false, false, true}
	MergedState   = &State{"MERGED", true, false, false}
	ClosedState   = &State{"CLOSED", true, false, false}
	SkippedState  = &State{"SKIPPED", false, true, false}
)
