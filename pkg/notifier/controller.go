package notifier

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/slack"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

type Controller struct {
	logger       logr.Logger
	slackClients map[int32]slack.Client

	channelIds map[int32]string
}

func NewController(logger logr.Logger,
	slackClients map[int32]slack.Client, channelIds map[int32]string,
) *Controller {
	return &Controller{logger, slackClients, channelIds}
}

func (c *Controller) Receive(m model.CurrentAndNextTalk) error {
	ctx := logr.NewContext(context.Background(), c.logger)
	slackClient := c.slackClients[m.TrackId()]
	slackChannelId, ok := c.channelIds[m.TrackId()]
	if !ok {
		c.logger.Info(fmt.Sprintf("notifier is disabled on trackId %d", m.TrackId()))
		return nil
	}
	if err := slackClient.PostMessage(ctx, slackChannelId, ViewSession(m)); err != nil {
		return xerrors.Errorf("message: %w", err)
	}
	return nil
}
