package ds

import "github.com/samber/lo"

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

func AreUsersEqual(left []*BasicUser, right []*BasicUser) bool {
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
	LeadLabel      UserLabel = "lead"
	DeveloperLabel UserLabel = "developer"
)

type UserLabels []UserLabel

func (u UserLabels) Has(label UserLabel) bool {
	return lo.Contains(u, label)
}

type User struct {
	*BasicUser `bson:"basic_user"`
	SlackID    string     `bson:"slack_id"`
	Labels     UserLabels `bson:"labels"`
}
