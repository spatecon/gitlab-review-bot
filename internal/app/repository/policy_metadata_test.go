//go:build mongodb

package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

func TestRepository_PolicyMetadata(t *testing.T) {
	rep := repositoryHelper(t)

	team := &ds.Team{
		ID: "12345",
	}
	mr := &ds.MergeRequest{
		ID: 1,
	}
	policyName := ds.PolicyName("policyName")

	t.Run("should return empty", func(t *testing.T) {
		md, err := rep.PolicyMetadata(mr, team, policyName)
		require.NoError(t, err, "should not return error")
		require.Empty(t, md, "should return empty list")
	})

	md1, bsonErr := bson.Marshal(map[string]string{"f1": "v1"})
	require.NoError(t, bsonErr, "bson marshal failed")

	t.Run("create policy metadata", func(t *testing.T) {
		err := rep.UpdatePolicyMetadata(mr, team, policyName, md1)
		require.NoError(t, err, "failed to create policy metadata")
	})

	t.Run("should return created policy metadata", func(t *testing.T) {
		md, err := rep.PolicyMetadata(mr, team, policyName)
		require.NoError(t, err, "failed to get policy metadata")
		require.EqualValues(t, md1, md, "metadata should be the same")
	})

	t.Run("another policy", func(t *testing.T) {
		md, err := rep.PolicyMetadata(mr, team, "anotherPolicy")
		require.NoError(t, err, "failed to get policy metadata")
		require.Nil(t, md, "metadata should be not found")
	})

	md1, bsonErr = bson.Marshal(map[string]string{"f2": "v2"})
	require.NoError(t, bsonErr, "bson marshal failed")

	t.Run("update policy metadata", func(t *testing.T) {
		err := rep.UpdatePolicyMetadata(mr, team, policyName, md1)
		require.NoError(t, err, "failed to update policy metadata")
	})

	t.Run("should return updated policy metadata", func(t *testing.T) {
		md, err := rep.PolicyMetadata(mr, team, policyName)
		require.NoError(t, err, "failed to get policy metadata")
		require.EqualValues(t, md1, md, "metadata should be the same")
	})
}
