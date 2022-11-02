package repository

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func (r *Repository) Projects() ([]*ds.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	cursor, err := r.projects.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find projects")
	}

	projects := make([]*ds.Project, 0, 10)

	err = cursor.All(ctx, &projects)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode projects")
	}

	return projects, nil
}
