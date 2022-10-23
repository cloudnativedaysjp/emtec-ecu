package obswatcher

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/obsws"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/sharedmem"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/utils"
)

const (
	componentName          = "obswatcher"
	syncPeriod             = 10 * time.Second
	startPreparetionPeriod = 60 * time.Second
)

type Config struct {
	Development          bool
	Debug                bool
	Obs                  []ConfigObs
	NotificationSendChan chan<- model.Notification
}

type ConfigObs struct {
	Host      string
	Password  string
	DkTrackId int32
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

	mr := &sharedmem.Reader{
		UseStorageForDisableAutomation: true, UseStorageForTrack: true}

	eg, ctx := errgroup.WithContext(ctx)
	for _, obs := range conf.Obs {
		obswsClient, err := obsws.NewObsWebSocketClient(obs.Host, obs.Password)
		if err != nil {
			err := xerrors.Errorf("message: %w", err)
			logger.Error(err, "obsws.NewObsWebSocketClient() was failed")
			return err
		}
		obswatcher := obswatcher{obswsClient, mr, conf.NotificationSendChan}
		eg.Go(obswatcher.watch(ctx, logger, obs.DkTrackId))
	}
	if err := eg.Wait(); err != nil {
		err := xerrors.Errorf("message: %w", err)
		logger.Error(err, "eg.Wait() was failed")
		return err
	}
	return nil
}

type obswatcher struct {
	obswsClient          obsws.Client
	mr                   sharedmem.ReaderIface
	notificationSendChan chan<- model.Notification
}

func (w *obswatcher) watch(ctx context.Context, logger logr.Logger, trackId int32) func() error {
	return func() error {
		logger := logger.WithValues("trackId", trackId)

		tick := time.NewTicker(syncPeriod)
		if err := w.procedure(ctx, trackId); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
		for {
			select {
			case <-ctx.Done():
				logger.Info("context was done.")
				return nil
			case <-tick.C:
				ctx := context.Background()
				ctx = logr.NewContext(ctx, logger)
				if err := w.procedure(ctx, trackId); err != nil {
					return xerrors.Errorf("message: %w", err)
				}
			}
		}
	}
}

func (w *obswatcher) procedure(ctx context.Context, trackId int32) error {
	logger := utils.GetLogger(ctx).WithValues("trackId", trackId)

	if disabled, err := w.mr.DisableAutomation(trackId); err != nil {
		logger.Error(xerrors.Errorf("message: %w", err), "mr.DisableAutomation() was failed")
		return nil
	} else if disabled {
		logger.Info("DisableAutomation was true, skipped")
		return nil
	}

	track, err := w.mr.Track(trackId)
	if err != nil {
		logger.Info(err.Error(), "trackId", trackId)
		return nil
	}
	currentTalk, err := track.Talks.GetCurrentTalk()
	if err != nil {
		logger.Info(fmt.Sprintf("talks.GetCurrentTalk was failed: %v", err))
		// カンファレンス開始前の場合は処理を続けたいため return しない
		currentTalk = &model.Talk{}
	}
	nextTalk, err := track.Talks.GetNextTalk()
	if err != nil {
		logger.Info(fmt.Sprintf("talks.GetNextTalk was failed: %v", err))
		return nil
	}

	t, err := w.obswsClient.GetRemainingTimeOnCurrentScene(ctx)
	if err != nil {
		logger.Error(xerrors.Errorf("message: %w", err), "obswsClient.GetRemainingTimeOnCurrentScene() was failed")
		return nil
	}
	remainingMilliSecond := t.DurationMilliSecond - t.CursorMilliSecond

	if float64(startPreparetionPeriod/time.Millisecond) < remainingMilliSecond {
		logger.Info(fmt.Sprintf("remainingTime on current Scene's MediaInput is over %ds: continue",
			startPreparetionPeriod/time.Second),
			"duration", t.DurationMilliSecond/float64(time.Millisecond),
			"cursor", t.CursorMilliSecond/float64(time.Millisecond))
		return nil
	}
	logger.Info(fmt.Sprintf("remainingTime on current Scene's MediaInput is within %ds",
		startPreparetionPeriod/time.Second),
		"duration", t.DurationMilliSecond/float64(time.Millisecond),
		"cursor", t.CursorMilliSecond/float64(time.Millisecond))

	// sleep until MediaInput is finished
	time.Sleep(time.Duration(remainingMilliSecond) * time.Millisecond)

	if err := w.obswsClient.MoveSceneToNext(context.Background()); err != nil {
		logger.Error(xerrors.Errorf("message: %w", err), "obswsClient.MoveSceneToNext() on automated task was failed")
		return nil
	}
	logger.Info("automated task was completed. Scene was moved to next.")
	w.notificationSendChan <- model.NewNotificationSceneMovedToNext(
		*currentTalk, *nextTalk)
	logger.Info("notified to Slack regarding Scene was moved to next")
	return nil
}
