package notifier

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	slackgo "github.com/slack-go/slack"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/db"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/slack"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

const componentName = "notifier"

type Config struct {
	Development          bool
	Debug                bool
	Targets              []Target
	RedisHost            string
	NotificationRecvChan <-chan model.Notification
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

	slackClients := make(map[int32]slack.Client)
	channelIds := make(map[int32]string)

	redisClient, err := db.NewRedisClient(conf.RedisHost)
	if err != nil {
		return err
	}
	for _, target := range conf.Targets {
		slackClients[target.TrackId], err = slack.NewClient(target.SlackBotToken)
		if err != nil {
			msg := "slack.NewClient() was failed"
			logger.Error(err, msg)
			return xerrors.Errorf("%s: %w", msg, err)
		}
		channelIds[target.TrackId] = target.SlackChannelId
	}

	notifier := notifier{slackClients, channelIds, *redisClient, conf.NotificationRecvChan}
	if err := notifier.watch(ctx, logger); err != nil {
		return err
	}
	return nil
}

type notifier struct {
	slackClients         map[int32]slack.Client
	channelIds           map[int32]string
	db                   db.RedisClient
	notificationRecvChan <-chan model.Notification
}

func (n *notifier) watch(ctx context.Context, logger logr.Logger) error {
	for {
		select {
		case <-ctx.Done():
			logger.Info("context was done.")
			return nil
		case notification := <-n.notificationRecvChan:
			ctx := context.Background()
			messageWasPosted := false

			var trackId int32
			var msg slackgo.Msg
			switch m := notification.(type) {
			case *model.NotificationOnDkTimetable:
				trackId = m.TrackId()
				msg = ViewNextSessionWillBegin(m)
				defer func() {
					if messageWasPosted {
						if err := n.db.SetNextTalkNotification(ctx, *m); err != nil {
							logger.Error(xerrors.Errorf("message: %w", err), "set value to redis failed")
						}
					}
				}()
			case *model.NotificationSceneMovedToNext:
				trackId = m.TrackId()
				msg = ViewSceneMovedToNext(m)
			default:
				logger.Error(fmt.Errorf(
					"unknown Notification type: %v", reflect.TypeOf(m)), "unknown type")
				continue
			}
			sc, ok := n.slackClients[trackId]
			if !ok {
				logger.Info(fmt.Sprintf("notifier is disabled on trackId %d", trackId))
				return nil
			}
			if err := sc.PostMessage(ctx, n.channelIds[trackId], msg); err != nil {
				return xerrors.Errorf("message: %w", err)
			}
			messageWasPosted = true
		}
	}
}
