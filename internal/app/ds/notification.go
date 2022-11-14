package ds

type NotificationSettings struct {
	Cron            string `bson:"cron"`
	UserTemplate    string `bson:"user_template"`
	ChannelTemplate string `bson:"channel_template"`
}

// UserNotification is a set of variables that can be used in user notification templates.
type UserNotification struct {
	// AuthoredMR list of merge requests in review that were authored by the user.
	AuthoredMR []MergeRequest
	// ReviewerMR list of merge requests that should be reviewed.
	ReviewerMR []MergeRequest
}

// ChannelNotification is a set of variables that can be used in channel notification templates.
type ChannelNotification struct {
	// AverageCount of MRs per developer/member
	AverageCount int
	// TotalCount of uniq MRs in review state
	TotalCount int
	// LastEditedMR is the last edited MR in review state
	LastEditedMR MergeRequest
	// FirstEditedMR is the oldest MR in review state
	FirstEditedMR MergeRequest
}
