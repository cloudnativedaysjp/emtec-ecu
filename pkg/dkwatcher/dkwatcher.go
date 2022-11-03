package dkwatcher

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/emtec-ecu/pkg/infra/dreamkast"
	"github.com/cloudnativedaysjp/emtec-ecu/pkg/infra/sharedmem"
	"github.com/cloudnativedaysjp/emtec-ecu/pkg/metrics"
	"github.com/cloudnativedaysjp/emtec-ecu/pkg/model"
	"github.com/cloudnativedaysjp/emtec-ecu/pkg/utils"
)

const componentName = "dkwatcher"

type Config struct {
	Logger                           logr.Logger
	DkClient                         dreamkast.Client
	NotificationSendChan             chan<- model.Notification
	SyncPeriodSeconds                int
	HowManyMinutesBeforeNotification int
}

func Run(ctx context.Context, conf Config) error {
	logger := conf.Logger.WithName(componentName)
	mw := sharedmem.Writer{UseStorageForTrack: true}
	mr := sharedmem.Reader{UseStorageForDisableAutomation: true}

	dkwatcher := dkwatcher{conf.DkClient, mw, mr,
		conf.NotificationSendChan, time.Minute * time.Duration(conf.HowManyMinutesBeforeNotification)}
	if err := dkwatcher.watch(ctx, logger, time.NewTicker(
		time.Second*time.Duration(conf.SyncPeriodSeconds)),
	); err != nil {
		return err
	}
	return nil
}

type dkwatcher struct {
	dkClient             dreamkast.Client
	mw                   sharedmem.WriterIface
	mr                   sharedmem.ReaderIface
	notificationSendChan chan<- model.Notification
	// const variables
	HowManyDurationBeforeNotification time.Duration
}

func (w *dkwatcher) watch(ctx context.Context, logger logr.Logger, ticker *time.Ticker) error {
	if err := w.procedure(ctx); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			logger.Info("context was done.")
			return nil
		case <-ticker.C:
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

		if !track.Talks.IsStartNextTalkSoon(w.HowManyDurationBeforeNotification) {
			logger.Info("nextTalk is not start soon")
			continue
		}
		w.notificationSendChan <- notification
	}
	return nil
}
