package notifier

import (
	"context"

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
	if err := slackClient.PostMessage(ctx, c.channelIds[m.TrackId()], ViewSession(m)); err != nil {
		return xerrors.Errorf("message: %w", err)
	}
	return nil
}
