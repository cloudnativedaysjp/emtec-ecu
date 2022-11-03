package notifier

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-logr/logr"
	slackgo "github.com/slack-go/slack"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/db"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/slack"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

const componentName = "notifier"

type Config struct {
	Logger               logr.Logger
	Targets              []Target
	RedisClient          *db.RedisClient
	NotificationRecvChan <-chan model.Notification
}

type Target struct {
	TrackId        int32
	SlackBotToken  string
	SlackChannelId string
}

func Run(ctx context.Context, conf Config) error {
	logger := conf.Logger.WithName(componentName)

	slackClients := make(map[int32]slack.Client)
	channelIds := make(map[int32]string)

	// TODO(#57): move to cmd/server/main.go
	var err error
	for _, target := range conf.Targets {
		slackClients[target.TrackId], err = slack.NewClient(target.SlackBotToken)
		if err != nil {
			msg := "slack.NewClient() was failed"
			logger.Error(err, msg)
			return xerrors.Errorf("%s: %w", msg, err)
		}
		channelIds[target.TrackId] = target.SlackChannelId
	}

	notifier := notifier{slackClients, channelIds, conf.RedisClient, conf.NotificationRecvChan}
	if err := notifier.watch(ctx, logger); err != nil {
		return err
	}
	return nil
}

type notifier struct {
	slackClients         map[int32]slack.Client
	channelIds           map[int32]string
	db                   *db.RedisClient
	notificationRecvChan <-chan model.Notification
}

func (n *notifier) watch(ctx context.Context, logger logr.Logger) error {
	for {
		select {
		case <-ctx.Done():
			logger.Info("context was done.")
			return nil
		case notification := <-n.notificationRecvChan:
			if err := n.notify(logger, notification); err != nil {
				return err
			}
		}
	}
}

func (n *notifier) notify(logger logr.Logger, notification model.Notification) error {
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
	}
	sc, ok := n.slackClients[trackId]
	if !ok {
		logger.Info(fmt.Sprintf("notifier is disabled on trackId %d", trackId))
		return nil
	}
	if err := sc.PostMessage(ctx, n.channelIds[trackId], msg); err != nil {
		logger.Error(err, "PostMessage was failed")
		return nil
	}
	messageWasPosted = true
	return nil
}
