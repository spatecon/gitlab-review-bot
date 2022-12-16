package repository

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func (r *Repository) Teams() ([]*ds.Team, error) {
	ctx, cancel := context.WithTimeout(r.ctx, defaultTimeout)
	defer cancel()

	cursor, err := r.teams.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find teams")
	}

	teams := make([]*ds.Team, 0, 10)

	err = cursor.All(ctx, &teams)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode teams")
	}

	return teams, nil
}

func (r *Repository) UserBySlackID(slackID string) (*ds.User, *ds.Team, error) {
	ctx, cancel := context.WithTimeout(r.ctx, defaultTimeout)
	defer cancel()

	team := &ds.Team{}

	err := r.teams.FindOne(ctx, bson.M{"members": bson.M{"$elemMatch": bson.M{"slack_id": slackID}}}).Decode(team)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil, nil
	}

	var user *ds.User
	for _, member := range team.Members {
		if member.SlackID == slackID {
			user = member
			break
		}
	}

	return user, team, nil
}
