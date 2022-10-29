package ds

type Team struct {
	Name           string  `bson:"name"`
	Members        []*User `bson:"members"`
	SlackChannelID string  `bson:"slack_channel_id"`
	Policy         string  `bson:"policy"`
}
