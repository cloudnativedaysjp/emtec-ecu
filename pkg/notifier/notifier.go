package notifier

import (
	"context"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/slackwh"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

const componentName = "notifier"

type Config struct {
	Targets                      []Target
	NotificationEventReceiveChan <-chan model.Talk
}

type Target struct {
	TrackId    int32
	WebhookUrl string
}

func Run(ctx context.Context, conf Config) error {
	// setup logger
	zapConf := zap.NewProductionConfig()
	zapConf.DisableStacktrace = true // due to output wrapped error in errorVerbose
	zapLogger, err := zapConf.Build()
	if err != nil {
		return err
	}
	logger := zapr.NewLogger(zapLogger).WithName(componentName)
	ctx = logr.NewContext(ctx, logger)

	whClients := make(map[int32]slackwh.WebhookAPI)
	for _, target := range conf.Targets {
		whClients[target.TrackId] = slackwh.NewClient(target.WebhookUrl)
	}
	c := NewController(logger, whClients)

	for {
		select {
		case <-ctx.Done():
			logger.Info("context was done.")
			return nil
		case talk := <-conf.NotificationEventReceiveChan:
			if err := c.Receive(talk); err != nil {
				logger.Error(xerrors.Errorf("message: %w", err), "notification failed")
			}
		}
	}
}
