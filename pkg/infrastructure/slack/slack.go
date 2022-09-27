package slack

import (
	"context"

	"github.com/slack-go/slack"
	"golang.org/x/xerrors"
)

type ClientIface interface {
	PostMessage(ctx context.Context, channel string, msg slack.Msg) error
}

type Client struct {
	client    *slack.Client
	botUserId string
}

func NewClient(botToken string) (ClientIface, error) {
	client := slack.New(botToken)
	res, err := client.AuthTest()
	if err != nil {
		return nil, err
	}

	return &Client{client, res.UserID}, nil
}

func (s *Client) PostMessage(ctx context.Context, channel string, msg slack.Msg) error {
	_, _, err := s.client.PostMessageContext(ctx, channel,
		slack.MsgOptionText(msg.Text, false),
		slack.MsgOptionAttachments(msg.Attachments...),
		slack.MsgOptionBlocks(msg.Blocks.BlockSet...),
	)
	if err != nil {
		return xerrors.Errorf("message: %w", err)
	}
	return nil
}
