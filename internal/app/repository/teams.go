package repository

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func (r *Repository) Teams() ([]*ds.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
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
