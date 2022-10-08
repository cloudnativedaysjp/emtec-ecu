package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/cloudnativedaysjp/cnd-operation-server/cmd/server/config"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/dkwatcher"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/notifier"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/obswatcher"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/ws-proxy/server"
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

	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	talkStream := make(chan model.CurrentAndNextTalk, 16)

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
				Development: conf.Debug.Development,
				Debug:       conf.Debug.Debug,
				Obs:         configObs,
			})
		})
	}
	// dkwatcher
	if !conf.Debug.DisableDkWatcher {
		eg.Go(func() error {
			return dkwatcher.Run(ctx, dkwatcher.Config{
				Development:               conf.Debug.Development,
				Debug:                     conf.Debug.Debug,
				EventAbbr:                 conf.Dreamkast.EventAbbr,
				DkEndpointUrl:             conf.Dreamkast.EndpointUrl,
				Auth0Domain:               conf.Dreamkast.Auth0Domain,
				Auth0ClientId:             conf.Dreamkast.Auth0ClientId,
				Auth0ClientSecret:         conf.Dreamkast.Auth0ClientSecret,
				Auth0ClientAudience:       conf.Dreamkast.Auth0ClientAudience,
				NotificationEventSendChan: talkStream,
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
				Development:                  conf.Debug.Development,
				Debug:                        conf.Debug.Debug,
				Targets:                      targets,
				NotificationEventReceiveChan: talkStream,
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
				Debug:       conf.Debug.Debug,
				BindAddr:    conf.WsProxyBindAddr,
				Obs:         configObs,
			})
		})
	}

	if err := eg.Wait(); err != nil {
		cancel()
		log.Fatal(err)
	}
}
