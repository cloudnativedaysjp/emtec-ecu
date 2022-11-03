package obswatcher

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/obsws"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/sharedmem"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/utils"
)

const componentName = "obswatcher"

type Config struct {
	Logger                        logr.Logger
	Obs                           []ConfigObs
	NotificationSendChan          chan<- model.Notification
	SyncPeriodSeconds             int
	StartPreparationPeriodSeconds int
}

type ConfigObs struct {
	Host      string
	Password  string
	DkTrackId int32
}

func Run(ctx context.Context, conf Config) error {
	logger := conf.Logger.WithName(componentName)

	mr := &sharedmem.Reader{
		UseStorageForDisableAutomation: true, UseStorageForTrack: true}

	eg, ctx := errgroup.WithContext(ctx)
	// TODO(#57): move to cmd/server/main.go
	for _, obs := range conf.Obs {
		obswsClient, err := obsws.NewObsWebSocketClient(obs.Host, obs.Password)
		if err != nil {
			err := xerrors.Errorf("message: %w", err)
			logger.Error(err, "obsws.NewObsWebSocketClient() was failed")
			return err
		}
		obswatcher := obswatcher{obswsClient, mr, conf.NotificationSendChan,
			time.Minute * time.Duration(conf.StartPreparationPeriodSeconds)}
		eg.Go(obswatcher.watch(ctx, logger, time.NewTicker(time.Second*time.Duration(conf.SyncPeriodSeconds)), obs.DkTrackId))
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
	// const variables
	StartPreparationPeriod time.Duration
}

func (w *obswatcher) watch(ctx context.Context,
	logger logr.Logger, ticker *time.Ticker, trackId int32,
) func() error {
	return func() error {
		logger := logger.WithValues("trackId", trackId)

		if err := w.procedure(ctx, trackId); err != nil {
			return xerrors.Errorf("message: %w", err)
		}
		for {
			select {
			case <-ctx.Done():
				logger.Info("context was done.")
				return nil
			case <-ticker.C:
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

	// currentTalk と nextTalk のどちらかがオンデマンドなセッションの場合、obswatcher からは何もする必要が無いため return する
	if currentTalk.IsOnDemand() || nextTalk.IsOnDemand() {
		return nil
	}

	var remainingMilliSecond float64
	if currentTalk.IsRepeatedConent() {

		// currentTalk が繰り返し流れるコンテンツ (cf. CM) の場合、タイムテーブルをもとに待ち時間を計算する
		remainingMilliSecond = float64(nextTalk.RemainingDurationUntilStart().Milliseconds())

	} else {

		// currentTalk が recording session の場合、動画の残り時間を取得してその時間だけ待つ
		t, err := w.obswsClient.GetRemainingTimeOnCurrentScene(ctx)
		if err != nil {
			logger.Error(xerrors.Errorf("message: %w", err), "obswsClient.GetRemainingTimeOnCurrentScene() was failed")
			return nil
		}
		remainingMilliSecond = t.DurationMilliSecond - t.CursorMilliSecond
		logger = logger.WithValues(
			"duration", t.DurationMilliSecond/float64(time.Millisecond),
			"cursor", t.CursorMilliSecond/float64(time.Millisecond),
		)

	}

	if float64(w.StartPreparationPeriod/time.Millisecond) < remainingMilliSecond {
		logger.Info(fmt.Sprintf("remainingTime on current Scene's MediaInput is over %ds: continue",
			w.StartPreparationPeriod/time.Second))
		return nil
	}
	logger.Info(fmt.Sprintf("remainingTime on current Scene's MediaInput is within %ds",
		w.StartPreparationPeriod/time.Second))

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
