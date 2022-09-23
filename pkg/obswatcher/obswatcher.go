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

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/obsws"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/sharedmem"
)

const (
	componentName          = "obswatcher"
	syncPeriod             = 10 * time.Second
	startPreparetionPeriod = 60 * time.Second
)

type Config struct {
	Debug bool
	Obs   []ConfigObs
}

type ConfigObs struct {
	Host      string
	Password  string
	DkTrackId int32
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

	mr := &sharedmem.Reader{UseStorageForTalks: true, UseStorageForDisableAutomation: true}

	eg, ctx := errgroup.WithContext(ctx)
	for _, obs := range conf.Obs {
		obswsClient, err := obsws.NewObsWebSocketClient(obs.Host, obs.Password)
		if err != nil {
			return err
		}
		eg.Go(watch(ctx, obs.DkTrackId, obswsClient, mr))
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func watch(ctx context.Context, trackId int32,
	obswsClient obsws.ClientIface, mr sharedmem.ReaderIface,
) func() error {
	return func() error {
		logger, err := logr.FromContext(ctx)
		if err != nil {
			return err
		}
		logger = logger.WithValues("trackId", trackId)
		tick := time.NewTicker(syncPeriod)

		for {
			select {
			case <-ctx.Done():
				logger.Info("context was done.")
				return nil
			case <-tick.C:
				if err := procedure(ctx, trackId, obswsClient, mr); err != nil {
					return err
				}
			}
		}
	}
}

func procedure(ctx context.Context, trackId int32,
	obswsClient obsws.ClientIface, mr sharedmem.ReaderIface,
) error {
	logger, err := logr.FromContext(ctx)
	if err != nil {
		return err
	}
	logger = logger.WithValues("trackId", trackId)

	if ok, err := mr.DisableAutomation(trackId); err != nil {
		logger.Error(xerrors.Errorf("message: %w", err), "mr.DisableAutomation() was failed")
		return nil
	} else if ok {
		logger.Info("DisableAutomation was true, skipped")
		return nil
	}

	t, err := obswsClient.GetRemainingTimeOnCurrentScene(ctx)
	if err != nil {
		logger.Error(xerrors.Errorf("message: %w", err), "obswsClient.GetRemainingTimeOnCurrentScene() was failed")
		return nil
	}
	remainingTime := t.Duration - t.Cursor
	if float64(startPreparetionPeriod) > remainingTime {
		return nil
	}
	logger.Info(fmt.Sprintf("remainingTime on current Scene's MediaInput is within %d",
		startPreparetionPeriod), "duration", t.Duration, "cursor", t.Cursor)

	// sleep until MediaInput is finished
	time.Sleep(time.Duration(remainingTime) * time.Second)
	if err := obswsClient.MoveSceneToNext(context.Background()); err != nil {
		logger.Error(xerrors.Errorf("message: %w", err), "obswsClient.MoveSceneToNext() on automated task was failed")
		return nil
	}
	logger.Info("automated task was completed. Scene should be to next.")
	return nil
}
