package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultTimeout = 5 * time.Second

type Repository struct {
	conn *mongo.Client

	teams          *mongo.Collection
	projects       *mongo.Collection
	mergeRequests  *mongo.Collection
	policyMetadata *mongo.Collection
}

func New(conn *mongo.Client, databaseName string) (*Repository, error) {
	database := conn.Database(databaseName)

	r := &Repository{
		conn:           conn,
		teams:          database.Collection("teams"),
		projects:       database.Collection("projects"),
		mergeRequests:  database.Collection("merge_requests"),
		policyMetadata: database.Collection("policy_metadata"),
	}

	err := r.createIndexes()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create indexes")
	}

	return r, nil
}

func (r *Repository) createIndexes() error {
	_, err := r.mergeRequests.Indexes().CreateMany(context.Background(),
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{"iid", 1}, {"project", 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{"id", 1}},
				Options: options.Index().SetUnique(true),
			},
		})
	if err != nil {
		return errors.Wrap(err, "failed to create merge_requests indexes")
	}

	_, err = r.policyMetadata.Indexes().CreateMany(context.Background(),
		[]mongo.IndexModel{
			{
				Keys:    bson.D{{"mr_id", 1}, {"team_id", 1}, {"policy_name", 1}},
				Options: options.Index().SetUnique(true),
			},
		})
	if err != nil {
		return errors.Wrap(err, "failed to create policy_metadata indexes")
	}

	return nil
}
