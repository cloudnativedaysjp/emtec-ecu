package lib

import (
	"context"
	"time"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/metrics"
	"golang.org/x/xerrors"
)

type DreamkastClientWrapper struct {
	c             DreamkastClient
	nowFunc       func() time.Time
	dkEndpointUrl string
}

func NewDreamkastClientWrapper(dkEndpointUrl string) (DreamkastClient, error) {
	c, err := NewClient(dkEndpointUrl)
	if err != nil {
		return nil, xerrors.Errorf("message: %w", err)
	}
	return &DreamkastClientWrapper{
		c,
		func() time.Time { return time.Now() },
		dkEndpointUrl,
	}, nil
}

func (w DreamkastClientWrapper) GenerateAuth0Token(ctx context.Context, auth0Domain, auth0ClientId, auth0ClientSecret, auth0Audience string) error {
	return w.c.GenerateAuth0Token(ctx, auth0Domain, auth0ClientId, auth0ClientSecret, auth0Audience)
}

func (w DreamkastClientWrapper) ListTracks(ctx context.Context, eventAbbr string) (ListTracksResp, error) {
	metricsDao := metrics.DreamkastMetricsFromCtx(ctx)
	now := w.nowFunc()
	result, err := w.c.ListTracks(ctx, eventAbbr)
	metricsDao.ListTracks(w.nowFunc().Sub(now))
	return result, err
}

func (w DreamkastClientWrapper) ListTalks(ctx context.Context, eventAbbr string, trackId int32) (ListTalksResp, error) {
	metricsDao := metrics.DreamkastMetricsFromCtx(ctx)
	now := w.nowFunc()
	result, err := w.c.ListTalks(ctx, eventAbbr, trackId)
	metricsDao.ListTalks(w.nowFunc().Sub(now))
	return result, err
}

func (w DreamkastClientWrapper) UpdateTalk(ctx context.Context, talkId int32, onAir bool) error {
	metricsDao := metrics.DreamkastMetricsFromCtx(ctx)
	now := w.nowFunc()
	err := w.c.UpdateTalk(ctx, talkId, onAir)
	metricsDao.UpdateTalk(w.nowFunc().Sub(now))
	return err

}
