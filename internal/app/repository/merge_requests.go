package repository

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func (r *Repository) MergeRequestByID(id int) (*ds.MergeRequest, error) {
	ctx, cancel := context.WithTimeout(r.ctx, defaultTimeout)
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
	ctx, cancel := context.WithTimeout(r.ctx, defaultTimeout)
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

func (r *Repository) MergeRequestsByAuthor(authorID []int) ([]*ds.MergeRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	cursor, err := r.mergeRequests.Find(ctx, bson.D{{"author.gitlab_id", bson.M{"$in": authorID}}})
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

func (r *Repository) MergeRequestsByReviewer(reviewerID []int) ([]*ds.MergeRequest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	cursor, err := r.mergeRequests.Find(ctx,
		bson.M{"reviewers": bson.M{"$elemMatch": bson.M{"gitlab_id": bson.M{"$in": reviewerID}}}},
	)
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
	ctx, cancel := context.WithTimeout(r.ctx, defaultTimeout)
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
