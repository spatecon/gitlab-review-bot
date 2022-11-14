package slack

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"go.uber.org/ratelimit"
)

const defaultTimeout = 5 * time.Second

type Client struct {
	slack *slack.Client
	rl    ratelimit.Limiter
}

func New(token string) (*Client, error) {
	c := &Client{
		slack: slack.New(token),
		rl:    ratelimit.New(1),
	}

	return c, nil
}

// SendMessage sends message to slack channel or user
func (c *Client) SendMessage(recipientID string, message string) error {
	c.rl.Take()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	_, _, err := c.slack.PostMessageContext(ctx,
		recipientID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		return errors.Wrapf(err, "failed to send slack message to (%s)", recipientID)
	}

	return nil
}
