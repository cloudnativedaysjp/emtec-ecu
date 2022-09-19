package notifier

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/slackwh"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

type Controller struct {
	logger    logr.Logger
	whClients map[int32]slackwh.WebhookAPI
}

func NewController(logger logr.Logger, whClients map[int32]slackwh.WebhookAPI) *Controller {
	return &Controller{logger, whClients}
}

func (c *Controller) Receive(talk model.Talk) error {
	ctx := logr.NewContext(context.Background(), c.logger)
	whClient := c.whClients[talk.TrackId]
	switch talk.Type {
	case model.TalkType_OnlineSession:
		if err := whClient.Send(ctx, viewOnlineSession(talk)); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
	case model.TalkType_RecordingSession:
		if err := whClient.Send(ctx, viewRecordingSession(talk)); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
	case model.TalkType_Commercial:
		if err := whClient.Send(ctx, viewCommercial(talk)); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
	case model.TalkType_Opening:
		if err := whClient.Send(ctx, viewOpening(talk)); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
	case model.TalkType_Ending:
		if err := whClient.Send(ctx, viewEnding(talk)); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
	default:
		return fmt.Errorf("unknown talk type")
	}
	return nil
}
