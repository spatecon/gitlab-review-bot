package ds

type UserLabel string

const (
	Lead UserLabel = "lead"
)

type BasicUser struct {
	Name     string
	GitLabID string
}

type User struct {
	BasicUser
	SlackID string
	Label   []UserLabel
}
