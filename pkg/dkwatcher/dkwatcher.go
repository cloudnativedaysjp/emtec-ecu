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
	NotificationEventSendChan chan<- model.Talk
}

const (
	syncPeriod = 30 * time.Second
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

	mw := sharedmem.Writer{UseStorageForTalks: true}
	mr := sharedmem.Reader{UseStorageForDisableAutomation: true}

	tick := time.NewTicker(syncPeriod)
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
	notificationEventSendChan chan<- model.Talk,
) error {
	logger := utils.GetLogger(ctx)

	talksList, err := dkClient.ListTalks(ctx)
	if err != nil {
		logger.Error(xerrors.Errorf("message: %w", err), "dkClient.ListTalks was failed")
		return nil
	}
	for _, talks := range talksList {
		track, currentTalkListNum := talks.GetCurrentTalk()
		if track == nil {
			logger.Info("Fail to GetCurrentTalk, skipped")
			continue
		}
		trackId := track.TrackId
		logger = logger.WithValues("trackId", trackId)

		if disabled, err := mr.DisableAutomation(trackId); err != nil {
			logger.Error(xerrors.Errorf("message: %w", err), "mr.DisableAutomation() was failed")
			return nil
		} else if disabled {
			logger.Info("DisableAutomation was true, skipped")
			continue
		}

		if err := mw.SetTalks(trackId, talks); err != nil {
			logger.Error(xerrors.Errorf("message: %w", err), "mw.SetTalks was failed")
			continue
		}
		if talks.WillStartNextTalkSince() {
			notificationEventSendChan <- talks.GetNextTalk(currentTalkListNum)
		}
	}
	return nil
}
