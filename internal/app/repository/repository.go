package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

const defaultTimeout = 5 * time.Second

type Repository struct {
	conn *mongo.Client

	teams         *mongo.Collection
	projects      *mongo.Collection
	mergeRequests *mongo.Collection
}

func New(conn *mongo.Client, databaseName string) (*Repository, error) {
	database := conn.Database(databaseName)

	return &Repository{
		conn:          conn,
		teams:         database.Collection("teams"),
		projects:      database.Collection("projects"),
		mergeRequests: database.Collection("merge_requests"),
	}, nil
}

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

func (r *Repository) MergeRequestByID(id int) (*ds.MergeRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	mergeRequest := &ds.MergeRequest{}

	err := r.mergeRequests.FindOne(ctx, bson.D{{"id", id}}).Decode(&mergeRequest)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to find merge request")
	}

	return mergeRequest, nil
}

func (r *Repository) MergeRequestsByProject(projectID int) ([]*ds.MergeRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	cursor, err := r.mergeRequests.Find(ctx, bson.D{{"project_id", projectID}})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find merge requests")
	}

	mrs := make([]*ds.MergeRequest, 0, 100)

	err = cursor.All(ctx, &mrs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode merge requests")
	}

	return mrs, nil
}

func (r *Repository) UpsertMergeRequest(mr *ds.MergeRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	opts := &options.UpdateOptions{}
	opts.SetUpsert(true)

	_, err := r.mergeRequests.UpdateOne(ctx,
		bson.D{{"id", mr.ID}},
		bson.D{{"$set", mr}},
		opts)
	if err != nil {
		return errors.Wrap(err, "failed to upsert merge request")
	}

	return nil
}
