package dkwatcher

import (
	"context"
	"time"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/dreamkast"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/sharedmem"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

const componentName = "dkwatcher"

type Config struct {
	EventAbbr                 string
	DkEndpointUrl             string
	Auth0Domain               string
	Auth0ClientId             string
	Auth0ClientSecret         string
	Auth0ClientAudience       string
	NotificationEventSendChan chan<- model.Talk
}

const (
	syncPeriod                = 30 * time.Second
	howManyMinutesUntilNotify = 5 * time.Minute
)

func Run(ctx context.Context, conf Config) error {
	// setup logger
	zapConf := zap.NewProductionConfig()
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

	mw := sharedmem.Writer{UseStorageForTalks: true}

	tick := time.NewTicker(syncPeriod)
	for {
		select {
		case <-ctx.Done():
			logger.Info("context was done.")
			return nil
		case <-tick.C:
			talks, err := dkClient.ListTalks(ctx)
			if err != nil {
				logger.Error(xerrors.Errorf("message: %w", err), "dkClient.ListTalks was failed")
			}
			if err := mw.SetTalks(talks); err != nil {
				logger.Error(xerrors.Errorf("message: %w", err), "mw.SetTalks was failed")
			}
			if talks.WillStartNextTalkSince(howManyMinutesUntilNotify) {
				conf.NotificationEventSendChan <- talks.GetNextTalk()
			}
		}
	}
}
