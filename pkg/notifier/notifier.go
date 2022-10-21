package notifier

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/infra"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/slack"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

const componentName = "notifier"

type Config struct {
	Development                  bool
	Debug                        bool
	Targets                      []Target
	RedisHost                    string
	NotificationEventReceiveChan <-chan model.CurrentAndNextTalk
}

type Target struct {
	TrackId        int32
	SlackBotToken  string
	SlackChannelId string
}

func Run(ctx context.Context, conf Config) error {
	// setup logger
	zapConf := zap.NewProductionConfig()
	if conf.Development {
		zapConf = zap.NewDevelopmentConfig()
	}
	zapConf.DisableStacktrace = true // due to output wrapped error in errorVerbose
	zapLogger, err := zapConf.Build()
	if err != nil {
		return err
	}
	logger := zapr.NewLogger(zapLogger).WithName(componentName)
	ctx = logr.NewContext(ctx, logger)

	slackClients := make(map[int32]slack.Client)
	channelIds := make(map[int32]string)

	redisClient, err := infra.NewRedisClient(conf.RedisHost)
	if err != nil {
		return err
	}
	for _, target := range conf.Targets {
		slackClients[target.TrackId], err = slack.NewClient(target.SlackBotToken)
		if err != nil {
			logger.Error(err, "slack.NewClient() was failed")
			return xerrors.Errorf("message: %w", err)
		}
		channelIds[target.TrackId] = target.SlackChannelId
	}
	c := NewController(logger, slackClients, channelIds)

	for {
		select {
		case <-ctx.Done():
			logger.Info("context was done.")
			return nil
		case talk := <-conf.NotificationEventReceiveChan:
			if err := c.Receive(talk); err != nil {
				logger.Error(xerrors.Errorf("message: %w", err), "notification failed")
			}
			if result := redisClient.Client.Set(ctx, infra.NextTalkNotificationKey, infra.NextTalkNotificationAlreadySent, infra.RedisExpiration); result.Err() != nil {
				return xerrors.Errorf("message: %w", result.Err())
			}
		}
	}
}
