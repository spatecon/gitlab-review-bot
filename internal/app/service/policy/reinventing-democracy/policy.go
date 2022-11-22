package reinventing_democracy

/**

Policy:					Reinventing Democracy
Reviewers Rotation:		random pick 2 developers from the team
Final Approve:			2 approves from the team

*/

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/zyedidia/generic/set"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

const (
	PolicyName ds.PolicyName = "reinventing_democracy"
	// DevelopersCount number of developers to be picked (also count of dev approves)
	DevelopersCount = 2
)

type Repository interface {
	// PolicyMetadata returns policy metadata for the given merge request
	PolicyMetadata(mr *ds.MergeRequest, team *ds.Team, policy ds.PolicyName) (bson.Raw, error)
	// UpdatePolicyMetadata updates policy metadata for the given merge request
	UpdatePolicyMetadata(mr *ds.MergeRequest, team *ds.Team, policy ds.PolicyName, d bson.Raw) error
}

type GitlabClient interface {
	// SetReviewers overwrites reviewers list for the merge request
	SetReviewers(mr *ds.MergeRequest, reviewers []int) error
}

type Policy struct {
	r Repository
	g GitlabClient
}

func New(r Repository, g GitlabClient) *Policy {
	return &Policy{
		r: r,
		g: g,
	}

}

type metadata struct {
	Approved     bool `bson:"approved"`
	ReviewersSet bool `bson:"reviewers_set"`
}

func (p *Policy) skip(mr *ds.MergeRequest, team *ds.Team) bool {
	// belongs to the team
	if !team.Teammate(mr.Author) {
		return true
	}

	// skip closed, merged, locked
	if mr.State != ds.StateOpened {
		return true
	}

	// not a draft
	if mr.Draft {
		return true
	}

	// not a release branch
	if strings.Contains(mr.SourceBranch, "release/") {
		return true
	}

	return false
}

func (p *Policy) ProcessChanges(team *ds.Team, mr *ds.MergeRequest) (err error) {
	if p.skip(mr, team) {
		return nil
	}

	// load metadata
	md := metadata{}

	raw, err := p.r.PolicyMetadata(mr, team, PolicyName)
	if err != nil {
		return errors.Wrap(err, "failed to get policy metadata")
	}

	if raw != nil {
		err = bson.Unmarshal(raw, &md)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal policy metadata")
		}
	}

	// save metadata
	defer func() {
		raw, err = bson.Marshal(md)
		if err != nil {
			err = errors.Wrap(err, "failed to marshal policy metadata")
		}

		err = p.r.UpdatePolicyMetadata(mr, team, PolicyName, raw)
		if err != nil {
			err = errors.Wrap(err, "failed to update policy metadata")
		}
	}()

	// wasn't set before
	if md.ReviewersSet {
		// check for final approved
		md.Approved = p.IsApproved(team, mr)

		return nil
	}

	// then set reviewers
	md.ReviewersSet, err = p.setReviewers(team, mr)
	if err != nil {
		return errors.Wrap(err, "failed to set reviewers")
	}

	return nil
}

func (p *Policy) setReviewers(team *ds.Team, mr *ds.MergeRequest) (bool, error) {
	reviewersSet := set.NewMapset[int]()

	for _, reviewer := range mr.Reviewers {
		reviewersSet.Put(reviewer.GitLabID)
	}

	devs := ds.Developers(team.Members)

	// without author
	devs = lo.Filter(devs, func(user *ds.User, _ int) bool {
		return user.GitLabID != mr.Author.GitLabID
	})

	// randomize
	devs = lo.Shuffle(devs)

	//TODO: not set developers that already set as reviewers and take it into account when counting efficientReviewersCount
	efficientReviewersCount := 0

	for i, dev := range devs {
		efficientReviewersCount++

		if i >= DevelopersCount {
			break
		}

		reviewersSet.Put(dev.GitLabID)
	}

	err := p.g.SetReviewers(mr, reviewersSet.Keys())
	if err != nil {
		return false, errors.Wrap(err, "failed to set reviewers")
	}

	if efficientReviewersCount < DevelopersCount {
		return false, nil
	}

	return true, nil
}

func (p *Policy) IsApproved(team *ds.Team, mr *ds.MergeRequest, byAll ...*ds.BasicUser) bool {
	if p.skip(mr, team) {
		// true means the MR was already approved
		// or didn't meet "need approve" state yet
		return true
	}

	if len(byAll) > 0 {
		allNeeded := set.NewMapset[int]()
		for _, user := range byAll {
			allNeeded.Put(user.GitLabID)
		}

		for _, user := range mr.Approves {
			allNeeded.Remove(user.GitLabID)
		}

		return allNeeded.Size() == 0 // all passed users approved the merge request
	}

	last := DevelopersCount

	for _, user := range mr.Approves {
		if user.GitLabID == mr.Author.GitLabID {
			continue
		}

		teammate, ok := lo.Find(team.Members, func(member *ds.User) bool {
			return member.GitLabID == user.GitLabID
		})

		if !ok {
			continue
		}

		if !teammate.Labels.Has(ds.DeveloperLabel) {
			continue
		}

		last--
	}

	// approved condition
	return last <= 0
}
