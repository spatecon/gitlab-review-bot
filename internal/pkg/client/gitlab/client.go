package gitlab

import (
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type Client struct {
	// TODO: consider using rate limiter

	gitlab *gitlab.Client
}

func New(token string) (*Client, error) {
	glClient, err := gitlab.NewClient(token)
	if err != nil {
		return nil, errors.Wrap(err, "error creating gitlab client")
	}

	return &Client{
		gitlab: glClient,
	}, nil
}
