package notifier

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/slack"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

type Controller struct {
	logger       logr.Logger
	slackClients map[int32]slack.ClientIface
	channelIds   map[int32]string
}

func NewController(logger logr.Logger,
	slackClients map[int32]slack.ClientIface, channelIds map[int32]string,
) *Controller {
	return &Controller{logger, slackClients, channelIds}
}

func (c *Controller) Receive(talk model.Talk) error {
	ctx := logr.NewContext(context.Background(), c.logger)
	slackClient := c.slackClients[talk.TrackId]
	switch talk.Type {
	case model.TalkType_OnlineSession:
		if err := slackClient.PostMessage(ctx, c.channelIds[talk.TrackId], viewOnlineSession(talk)); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
	case model.TalkType_RecordingSession:
		if err := slackClient.PostMessage(ctx, c.channelIds[talk.TrackId], viewRecordingSession(talk)); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
	case model.TalkType_Commercial:
		if err := slackClient.PostMessage(ctx, c.channelIds[talk.TrackId], viewCommercial(talk)); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
	case model.TalkType_Opening:
		if err := slackClient.PostMessage(ctx, c.channelIds[talk.TrackId], viewOpening(talk)); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
	case model.TalkType_Ending:
		if err := slackClient.PostMessage(ctx, c.channelIds[talk.TrackId], viewEnding(talk)); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
	default:
		return fmt.Errorf("unknown talk type")
	}
	return nil
}
