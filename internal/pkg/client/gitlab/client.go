package gitlab

import (
	"time"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

type Client struct {
	// TODO: consider using rate limiter

	gitlab        *gitlab.Client
	pullMRsPeriod time.Duration
}

func New(token string, pullMRsPeriod time.Duration) (*Client, error) {
	glClient, err := gitlab.NewClient(token)
	if err != nil {
		return nil, errors.Wrap(err, "error creating gitlab client")
	}

	return &Client{
		gitlab:        glClient,
		pullMRsPeriod: pullMRsPeriod,
	}, nil
}
