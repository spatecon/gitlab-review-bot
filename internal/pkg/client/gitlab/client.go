package gitlab

import (
	"context"

	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
	"go.uber.org/ratelimit"
)

type Client struct {
	// TODO: consider using rate limiter
	ctx context.Context

	gitlab *gitlab.Client
	rl     ratelimit.Limiter
}

func New(rootCtx context.Context, token string) (*Client, error) {
	glClient, err := gitlab.NewClient(token)
	if err != nil {
		return nil, errors.Wrap(err, "error creating gitlab client")
	}

	return &Client{
		ctx:    rootCtx,
		gitlab: glClient,
		rl:     ratelimit.New(3),
	}, nil
}
