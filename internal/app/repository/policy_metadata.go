package repository

import (
	"context"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

type policyMetadataDocument struct {
	MergeRequestID int      `bson:"mr_id"`
	PolicyName     string   `bson:"policy_name"`
	TeamID         string   `bson:"team_id"`
	Data           bson.Raw `bson:"data"`
}

func (r *Repository) PolicyMetadata(mr *ds.MergeRequest, team *ds.Team, policy ds.PolicyName) (bson.Raw, error) {
	ctx, cancel := context.WithTimeout(r.ctx, defaultTimeout)
	defer cancel()

	result := policyMetadataDocument{}

	err := r.policyMetadata.FindOne(ctx, bson.D{
		{"mr_id", mr.ID},
		{"policy_name", policy},
		{"team_id", team.ID},
	}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to find policy metadata")
	}

	return result.Data, nil
}

func (r *Repository) UpdatePolicyMetadata(mr *ds.MergeRequest, team *ds.Team, policy ds.PolicyName, d bson.Raw) error {
	ctx, cancel := context.WithTimeout(r.ctx, defaultTimeout)
	defer cancel()

	doc := policyMetadataDocument{
		MergeRequestID: mr.ID,
		PolicyName:     string(policy),
		TeamID:         team.ID,
		Data:           d,
	}

	_, err := r.policyMetadata.UpdateOne(ctx, bson.D{
		{"mr_id", mr.ID},
		{"policy_name", policy},
		{"team_id", team.ID},
	}, bson.D{{"$set", doc}}, &options.UpdateOptions{Upsert: lo.ToPtr(true)})
	if err != nil {
		return errors.Wrap(err, "failed to update policy metadata")
	}

	return nil
}
