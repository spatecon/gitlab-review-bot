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
	PolicyName ds.PolicyName = "rd"
	// RequiredDevelopersCount number of developers to be picked (also count of dev approves)
	RequiredDevelopersCount = 2
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
	ApprovedByPolicy  bool  `bson:"approved_by_policy"`
	ReviewersSet      bool  `bson:"reviewers_set"`
	ReviewersByPolicy []int `bson:"reviewers_by_policy"`
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

	// weren't set before
	if md.ReviewersSet {
		// check if approved by policy
		md.ApprovedByPolicy = p.ApprovedByPolicy(team, mr)

		return nil
	}

	// then set reviewers
	err = p.setReviewers(team, mr, &md)
	if err != nil {
		return errors.Wrap(err, "failed to set reviewers")
	}

	return nil
}

func (p *Policy) setReviewers(team *ds.Team, mr *ds.MergeRequest, md *metadata) error {
	md.ReviewersSet = true

	reviewersSet := set.NewMapset[int]()

	for _, reviewer := range mr.Reviewers {
		reviewersSet.Put(reviewer.GitLabID)
	}

	devs := ds.Developers(team.Members)

	// without author
	devs = lo.Filter(devs, func(user *ds.User, _ int) bool {
		return user.GitLabID != mr.Author.GitLabID
	})

	developersSet := set.NewMapset[int]()
	for _, dev := range devs {
		developersSet.Put(dev.GitLabID)
	}

	// developers who are reviewers
	inner := reviewersSet.Intersection(developersSet)

	// count of developers which needed to set as reviewers
	needDevsCount := RequiredDevelopersCount - inner.Size()

	// if we have enough reviewers from developers
	if needDevsCount <= 0 {
		return nil
	}

	// developers who are not reviewers
	notPickedDevsSet := developersSet.Difference(reviewersSet)
	notPickedDevs := lo.Shuffle(notPickedDevsSet.Keys())

	for _, dev := range notPickedDevs {
		if needDevsCount <= 0 {
			break
		}

		md.ReviewersByPolicy = append(md.ReviewersByPolicy, dev)

		reviewersSet.Put(dev)
		needDevsCount--
	}

	err := p.g.SetReviewers(mr, reviewersSet.Keys())
	if err != nil {
		md.ReviewersSet = false
		return errors.Wrap(err, "failed to set reviewers")
	}

	return nil
}

func (p *Policy) ApprovedByUser(team *ds.Team, mr *ds.MergeRequest, byAll ...*ds.BasicUser) bool {
	if p.skip(mr, team) {
		// true means the MR meet "need approve" state yet
		// or closed, merged, locked
		return true
	}

	if len(byAll) == 0 {
		return false
	}

	allNeeded := set.NewMapset[int]()
	for _, user := range byAll {
		allNeeded.Put(user.GitLabID)
	}

	approvesSet := set.NewMapset[int]()
	for _, approve := range mr.Approves {
		approvesSet.Put(approve.GitLabID)
	}

	return allNeeded.Difference(approvesSet).Size() == 0 // all passed users approved the merge request
}

func (p *Policy) ApprovedByPolicy(team *ds.Team, mr *ds.MergeRequest) bool {
	if p.skip(mr, team) {
		// true means the MR meet "need approve" state yet
		// or closed, merged, locked
		return true
	}

	left := RequiredDevelopersCount

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

		left--
	}

	// approved condition
	return left <= 0
}
