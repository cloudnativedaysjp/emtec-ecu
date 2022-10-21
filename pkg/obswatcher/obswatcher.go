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
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/utils"
)

const (
	componentName          = "obswatcher"
	syncPeriod             = 10 * time.Second
	startPreparetionPeriod = 60 * time.Second
)

type Config struct {
	Development bool
	Debug       bool
	Obs         []ConfigObs
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
	ctx = logr.NewContext(ctx, logger)

	mr := &sharedmem.Reader{UseStorageForDisableAutomation: true}

	eg, ctx := errgroup.WithContext(ctx)
	for _, obs := range conf.Obs {
		obswsClient, err := obsws.NewObsWebSocketClient(obs.Host, obs.Password)
		if err != nil {
			err := xerrors.Errorf("message: %w", err)
			logger.Error(err, "obsws.NewObsWebSocketClient() was failed")
			return err
		}
		obswatcher := obswatcher{obswsClient, mr}
		eg.Go(obswatcher.watch(ctx, obs.DkTrackId))
	}
	if err := eg.Wait(); err != nil {
		err := xerrors.Errorf("message: %w", err)
		logger.Error(err, "eg.Wait() was failed")
		return err
	}
	return nil
}

type obswatcher struct {
	obswsClient obsws.Client
	mr          sharedmem.ReaderIface
}

func (w *obswatcher) watch(ctx context.Context, trackId int32) func() error {
	return func() error {
		logger := utils.GetLogger(ctx).WithValues("trackId", trackId)

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
	logger.Info("automated task was completed. Scene should be to next.")
	return nil
}
