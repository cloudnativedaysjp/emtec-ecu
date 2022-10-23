package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/cloudnativedaysjp/cnd-operation-server/cmd/server/config"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/dkwatcher"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/db"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/dreamkast"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infra/sharedmem"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/metrics"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/notifier"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/obswatcher"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/ws-proxy/server"
	"github.com/go-logr/zapr"
)

func main() {
	var confFile string
	flag.StringVar(&confFile, "config", "", "filename of config (for example, refer to `example.yaml` on this repository)")
	flag.Parse()
	if confFile == "" {
		fmt.Println("flag --config must be specified")
		os.Exit(1)
	}
	conf, err := config.LoadConf(confFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// logger
	zapConf := zap.NewProductionConfig()
	if conf.Debug.Development {
		zapConf = zap.NewDevelopmentConfig()
	}
	zapConf.DisableStacktrace = true // due to output wrapped error in errorVerbose
	zapLogger, err := zapConf.Build()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger := zapr.NewLogger(zapLogger)

	// context
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	// metrics
	go func() {
		_ = metrics.RunCndOperationServer(conf.Metrics.BindAddr)
	}()

	// channels
	notificationStream := make(chan model.Notification, 16)

	//
	// Register
	//

	// dreamkast client
	dkClient, err := dreamkast.NewClient(
		conf.Dreamkast.EventAbbr, conf.Dreamkast.EndpointUrl,
		conf.Dreamkast.Auth0Domain, conf.Dreamkast.Auth0ClientId,
		conf.Dreamkast.Auth0ClientSecret, conf.Dreamkast.Auth0ClientAudience)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// redis client
	redisClient, err := db.NewRedisClient(conf.Redis.Host)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//
	// Initialize
	//
	mw := sharedmem.Writer{UseStorageForDisableAutomation: true}
	for _, track := range conf.Tracks {
		if err := mw.SetDisableAutomation(track.DkTrackId, false); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	//
	// Run components
	//

	// obswatcher
	if !conf.Debug.DisableObsWatcher {
		eg.Go(func() error {
			var configObs []obswatcher.ConfigObs
			for _, track := range conf.Tracks {
				configObs = append(configObs, obswatcher.ConfigObs{
					DkTrackId: track.DkTrackId,
					Host:      track.Obs.Host,
					Password:  track.Obs.Password,
				})
			}
			return obswatcher.Run(ctx, obswatcher.Config{
				Logger:               logger,
				Obs:                  configObs,
				NotificationSendChan: notificationStream,
			})
		})
	}
	// dkwatcher
	if !conf.Debug.DisableDkWatcher {
		eg.Go(func() error {
			return dkwatcher.Run(ctx, dkwatcher.Config{
				Logger:               logger,
				DkClient:             dkClient,
				RedisClient:          redisClient,
				NotificationSendChan: notificationStream,
			})
		})
	}
	// notifier
	if !conf.Debug.DisableNotifier {
		var targets []notifier.Target
		for _, track := range conf.Tracks {
			targets = append(targets, notifier.Target{
				TrackId:        track.DkTrackId,
				SlackBotToken:  track.Slack.BotToken,
				SlackChannelId: track.Slack.ChannelId,
			})
		}
		eg.Go(func() error {
			return notifier.Run(ctx, notifier.Config{
				Logger:               logger,
				Targets:              targets,
				RedisClient:          redisClient,
				NotificationRecvChan: notificationStream,
			})
		})
	}
	// ws-proxy
	if !conf.Debug.DisableWsProxy {
		eg.Go(func() error {
			var configObs []server.ConfigObs
			for _, track := range conf.Tracks {
				configObs = append(configObs, server.ConfigObs{
					DkTrackId: track.DkTrackId,
					Host:      track.Obs.Host,
					Password:  track.Obs.Password,
				})
			}
			return server.Run(ctx, server.Config{
				Development: conf.Debug.Development,
				Logger:      logger,
				ZapLogger:   zapLogger,
				BindAddr:    conf.WsProxy.BindAddr,
				Obs:         configObs,
			})
		})
	}

	if err := eg.Wait(); err != nil {
		cancel()
		log.Fatal(err)
	}
}
