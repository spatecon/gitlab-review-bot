package slack

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"go.uber.org/ratelimit"

	"github.com/spatecon/gitlab-review-bot/internal/app/ds"
)

const defaultTimeout = 5 * time.Second

type Client struct {
	ctx         context.Context
	slack       *slack.Client
	slackSocket *socketmode.Client
	rl          ratelimit.Limiter
	closer      chan struct{}
}

func (c *Client) Close() {
	c.closer <- struct{}{}
}

func New(rootCtx context.Context, botToken, appToken string) (*Client, error) {
	c := &Client{
		ctx:    rootCtx,
		slack:  slack.New(botToken, slack.OptionAppLevelToken(appToken)),
		rl:     ratelimit.New(1),
		closer: make(chan struct{}),
	}

	return c, nil
}

func (c *Client) Subscribe() (chan ds.UserEvent, error) {
	eventsChan := make(chan ds.UserEvent, 20)

	c.slackSocket = socketmode.New(c.slack)

	handler := socketmode.NewSocketmodeHandler(c.slackSocket)
	handler.HandleSlashCommand("/mr", mrRequestHandler(eventsChan)) // TODO: should be configurable
	handler.HandleDefault(func(evt *socketmode.Event, client *socketmode.Client) {
	})

	go func() {
		err := handler.RunEventLoop()
		if err != nil {
			log.Error().Err(err).Msg("error running socket mode event loop")
		}
	}()

	return eventsChan, nil
}

func mrRequestHandler(eventsChan chan ds.UserEvent) func(*socketmode.Event, *socketmode.Client) {
	return func(evt *socketmode.Event, client *socketmode.Client) {
		cmd, ok := evt.Data.(slack.SlashCommand)
		if !ok {
			return
		}

		eventsChan <- ds.UserEvent{
			Type:   ds.UserEventTypeMRRequest,
			UserID: cmd.UserID,
		}

		client.Ack(*evt.Request, nil)
	}
}
