package ds

type UserEventType int

const (
	UserEventTypeMRRequest UserEventType = iota
)

type UserEvent struct {
	Type   UserEventType
	UserID string
}
