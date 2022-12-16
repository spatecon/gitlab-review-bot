package slack

import (
	"context"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// SendMessage sends message to slack channel or user
func (c *Client) SendMessage(recipientID string, message string) error {
	c.rl.Take()

	ctx, cancel := context.WithTimeout(c.ctx, defaultTimeout)
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
