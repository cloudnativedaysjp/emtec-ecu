package slackwh

import (
	"context"

	"github.com/slack-go/slack"
)

type WebhookAPI interface {
	Send(context.Context, slack.WebhookMessage) error
}

type Client struct {
	WebhookUrl string
}

func NewClient(whUrl string) WebhookAPI {
	return &Client{whUrl}
}

func (c *Client) Send(ctx context.Context, msg slack.WebhookMessage) error {
	err := slack.PostWebhookContext(ctx, c.WebhookUrl, &msg)
	if err != nil {
		return err
	}
	return nil
}
