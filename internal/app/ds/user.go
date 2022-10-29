package ds

type BasicUser struct {
	Name     string `bson:"name"`
	GitLabID int    `bson:"gitlab_id"`
}

func EqualUser(user *BasicUser, other *BasicUser) bool {
	A := user == nil
	B := other == nil

	// both nil
	if A && B {
		return true
	}

	// one is not nil
	if A || B {
		return false
	}

	return user.GitLabID == other.GitLabID
}

func EqualUsers(left []*BasicUser, right []*BasicUser) bool {
	if len(left) != len(right) {
		return false
	}

	var found int
	seen := map[BasicUser]struct{}{}

	for _, elem := range left {
		seen[*elem] = struct{}{}
	}

	for _, elem := range right {
		if _, ok := seen[*elem]; ok {
			found++
		}
	}

	return found == len(right)
}

type UserLabel string

const (
	Lead UserLabel = "lead"
)

type User struct {
	BasicUser *BasicUser  `bson:"basic_user"`
	SlackID   string      `bson:"slack_id"`
	Labels    []UserLabel `bson:"labels"`
}
