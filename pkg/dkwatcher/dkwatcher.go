package dkwatcher

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/db"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/dreamkast"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/sharedmem"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/metrics"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/utils"
)

const componentName = "dkwatcher"

type Config struct {
	Development          bool
	Debug                bool
	EventAbbr            string
	DkEndpointUrl        string
	Auth0Domain          string
	Auth0ClientId        string
	Auth0ClientSecret    string
	Auth0ClientAudience  string
	RedisHost            string
	NotificationSendChan chan<- model.Notification
}

const (
	syncPeriod                = 30 * time.Second
	howManyMinutesUntilNotify = 5 * time.Minute
)

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

	dkClient, err := dreamkast.NewClient(conf.EventAbbr, conf.DkEndpointUrl,
		conf.Auth0Domain, conf.Auth0ClientId, conf.Auth0ClientSecret, conf.Auth0ClientAudience)
	if err != nil {
		return err
	}

	redisClient, err := db.NewRedisClient(conf.RedisHost)
	if err != nil {
		return err
	}

	mw := sharedmem.Writer{UseStorageForTrack: true}
	mr := sharedmem.Reader{UseStorageForDisableAutomation: true}

	dkwatcher := dkwatcher{dkClient, mw, mr, *redisClient, conf.NotificationSendChan}
	if err := dkwatcher.watch(ctx, logger); err != nil {
		return err
	}
	return nil
}

type dkwatcher struct {
	dkClient             dreamkast.Client
	mw                   sharedmem.WriterIface
	mr                   sharedmem.ReaderIface
	db                   db.RedisClient
	notificationSendChan chan<- model.Notification
}

func (w *dkwatcher) watch(ctx context.Context, logger logr.Logger) error {
	tick := time.NewTicker(syncPeriod)
	if err := w.procedure(ctx); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			logger.Info("context was done.")
			return nil
		case <-tick.C:
			ctx := context.Background()
			ctx = logr.NewContext(ctx, logger)
			ctx = metrics.SetDreamkastMetricsToCtx(ctx,
				metrics.NewDreamkastMetricsDao(w.dkClient.EndpointUrl()))
			if err := w.procedure(ctx); err != nil {
				return err
			}
		}
	}
}

func (w *dkwatcher) procedure(ctx context.Context) error {
	rootLogger := utils.GetLogger(ctx)

	tracks, err := w.dkClient.ListTracks(ctx)
	if err != nil {
		rootLogger.Error(xerrors.Errorf("message: %w", err), "dkClient.ListTalks was failed")
		return nil
	}
	for _, track := range tracks {
		logger := rootLogger.WithValues("trackId", track.Id)

		if err := w.mw.SetTrack(track); err != nil {
			logger.Error(xerrors.Errorf("message: %w", err), "mw.SetTrack was failed")
			continue
		}

		currentTalk, err := track.Talks.GetCurrentTalk()
		if err != nil {
			// カンファレンス開始前の場合は処理を続けたいため return しない
			logger.Info("currentTalk is none")
			currentTalk = &model.Talk{}
		}
		nextTalk, err := track.Talks.GetNextTalk()
		if err != nil {
			logger.Info("nextTalk is none")
			continue
		}
		notification := model.NewNotificationOnDkTimetable(
			*currentTalk, *nextTalk)

		if !track.Talks.IsStartNextTalkSoon(howManyMinutesUntilNotify) {
			logger.Info("nextTalk is not start soon. trackNo:%s", track.Id)
			continue
		}
		if notified, err := w.db.HasNextTalkNotificationAlreadyBeenSent(ctx, *notification); err != nil {
			logger.Error(xerrors.Errorf("message: %w", err), "db.GetNextTalkNotification() was failed")
			return err
		} else if notified {
			logger.Info("nextTalkNotification already sent . trackNo:%s", track.Id)
			continue
		}
		w.notificationSendChan <- notification
		logger.Info("notified to Slack regarding next talk will begin")
	}
	return nil
}
