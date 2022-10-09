package dkwatcher

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/dreamkast"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/sharedmem"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/utils"
)

const componentName = "dkwatcher"

type Config struct {
	Development               bool
	Debug                     bool
	EventAbbr                 string
	DkEndpointUrl             string
	Auth0Domain               string
	Auth0ClientId             string
	Auth0ClientSecret         string
	Auth0ClientAudience       string
	NotificationEventSendChan chan<- model.CurrentAndNextTalk
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
	ctx = logr.NewContext(ctx, logger)

	dkClient, err := dreamkast.NewClient(conf.EventAbbr, conf.DkEndpointUrl,
		conf.Auth0Domain, conf.Auth0ClientId, conf.Auth0ClientSecret, conf.Auth0ClientAudience)
	if err != nil {
		return err
	}

	mw := sharedmem.Writer{UseStorageForTrack: true}
	mr := sharedmem.Reader{UseStorageForDisableAutomation: true}

	tick := time.NewTicker(syncPeriod)
	if err := procedure(ctx, dkClient, mw, mr, conf.NotificationEventSendChan); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			logger.Info("context was done.")
			return nil
		case <-tick.C:
			if err := procedure(ctx, dkClient, mw, mr, conf.NotificationEventSendChan); err != nil {
				return err
			}
		}
	}
}

func procedure(ctx context.Context,
	dkClient dreamkast.ClientIface, mw sharedmem.WriterIface, mr sharedmem.ReaderIface,
	notificationEventSendChan chan<- model.CurrentAndNextTalk,
) error {
	logger := utils.GetLogger(ctx)

	tracks, err := dkClient.ListTracks(ctx)
	if err != nil {
		logger.Error(xerrors.Errorf("message: %w", err), "dkClient.ListTalks was failed")
		return nil
	}
	for _, track := range tracks {
		logger = logger.WithValues("trackId", track.Id)

		if disabled, err := mr.DisableAutomation(track.Id); err != nil {
			logger.Error(xerrors.Errorf("message: %w", err), "mr.DisableAutomation() was failed")
			return nil
		} else if disabled {
			logger.Info("DisableAutomation was true, skipped")
			continue
		}

		if err := mw.SetTrack(track); err != nil {
			logger.Error(xerrors.Errorf("message: %w", err), "mw.SetTrack was failed")
			continue
		}
		if track.Talks.WillStartNextTalkSince(howManyMinutesUntilNotify) {
			currentTalk, err := track.Talks.GetCurrentTalk()
			if err != nil {
				logger.Error(xerrors.Errorf("message: %w", err), "dkClient.GetCurrentTalk was failed")
				continue
			}
			nextTalk, err := track.Talks.GetNextTalk(currentTalk)
			if err != nil {
				logger.Error(xerrors.Errorf("message: %w", err), "talks.GetNextTalk was failed")
				continue
			}
			notificationEventSendChan <- model.CurrentAndNextTalk{
				Current: *currentTalk, Next: *nextTalk}
		}
	}
	return nil
}
