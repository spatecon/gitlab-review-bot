package ds

import "time"

type PolicyName string

type Team struct {
	ID            string               `bson:"_id"`
	Name          string               `bson:"name"`
	Members       []*User              `bson:"members"`
	Policy        PolicyName           `bson:"policy"`
	Notifications NotificationSettings `bson:"notifications"`
	CreatedAt     time.Time            `bson:"created_at"`
}

// Teammate checks if user is a member of a team
func (t *Team) Teammate(user *BasicUser) bool {
	for _, member := range t.Members {
		if member.BasicUser.GitLabID == user.GitLabID {
			return true
		}
	}

	return false
}

// Developers returns all developers of a team/list of users
func Developers(users []*User) []*User {
	devs := make([]*User, 0, len(users))

	for _, user := range users {
		if user.Labels.Has(DeveloperLabel) {
			devs = append(devs, user)
		}
	}

	return devs
}

// Lead returns first lead of a team/list of users
func Lead(users []*User) *User {
	for _, user := range users {
		if user.Labels.Has(LeadLabel) {
			return user
		}
	}

	return nil
}
