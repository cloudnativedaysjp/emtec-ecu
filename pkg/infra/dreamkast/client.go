package dreamkast

import (
	"context"
	"errors"

	"github.com/avast/retry-go"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/emtec-ecu/pkg/infra/dreamkast/lib"
	"github.com/cloudnativedaysjp/emtec-ecu/pkg/model"
	"github.com/cloudnativedaysjp/emtec-ecu/pkg/utils"
)

type Client interface {
	WithCredential(auth0Domain, auth0ClientId, auth0ClientSecret, auth0Audience string) Client
	EndpointUrl() string
	ListTracks(ctx context.Context, eventAbbr string) ([]model.Track, error)
	SetSpecifiedTalkOnAir(ctx context.Context, talkId int32) error
	SetNextTalkOnAir(ctx context.Context, eventAbbr string, trackId int32) error
}

type ClientImpl struct {
	client            lib.DreamkastClient
	endpointUrl       string
	auth0Domain       string
	auth0ClientId     string
	auth0ClientSecret string
	auth0Audience     string
}

func NewClient(dkEndpointUrl string) (Client, error) {
	c, err := lib.NewClient(dkEndpointUrl)
	if err != nil {
		return nil, err
	}
	return &ClientImpl{client: c, endpointUrl: dkEndpointUrl}, nil
}

func (c ClientImpl) WithCredential(auth0Domain, auth0ClientId, auth0ClientSecret, auth0Audience string) Client {
	c.auth0Domain = auth0Domain
	c.auth0ClientId = auth0ClientId
	c.auth0ClientSecret = auth0ClientSecret
	c.auth0Audience = auth0Audience
	return &c
}

func (c *ClientImpl) EndpointUrl() string {
	return c.endpointUrl
}

func (c *ClientImpl) ListTracks(ctx context.Context, eventAbbr string) ([]model.Track, error) {
	logger := utils.GetLogger(ctx)

	tracks, err := c.client.ListTracks(ctx, eventAbbr)
	if err != nil {
		return nil, xerrors.Errorf("message: %w", err)
	}
	var result []model.Track
	for _, track := range tracks {
		var talksModel model.Talks
		talks, err := c.client.ListTalks(ctx, eventAbbr, track.ID)
		if err != nil {
			return nil, xerrors.Errorf("message: %w", err)
		}
		for _, talk := range talks {
			t := model.Talk{
				Id:        talk.ID,
				TalkName:  talk.Title,
				TrackId:   track.ID,
				TrackName: track.Name,
				EventAbbr: eventAbbr,
			}
			talkType, err := t.GetTalkType(talk.Title, talk.PresentationMethod)
			if err != nil {
				err = xerrors.Errorf("message: %w", err)
				logger.Error(err, "GetTalkType() was failed")
				continue
			}
			t.Type = talkType
			for _, speaker := range talk.Speakers {
				t.SpeakerNames = append(t.SpeakerNames, speaker.Name)
			}

			t.StartAt, t.EndAt, err = t.GetActualStartAtAndEndAt(talk.ConferenceDayDate, talk.ActualStartTime, talk.ActualEndTime)
			if err != nil {
				return nil, xerrors.Errorf("message: %w", err)
			}
			talksModel = talksModel.AppendAndSort(t)
		}
		result = append(result, model.Track{
			Id:    track.ID,
			Name:  track.Name,
			Talks: talksModel,
		})
	}
	return result, nil
}

func (c *ClientImpl) SetSpecifiedTalkOnAir(ctx context.Context, talkId int32) error {
	// If Auth0Token has been expired, retry only once.
	err := retry.Do(
		func() (err error) {
			err = c.setSpecifiedTalkOnAir(ctx, talkId)
			return
		},
		retry.RetryIf(func(err error) bool {
			return errors.As(err, &lib.ErrorUnauthorized{})
		}),
		retry.Attempts(1),
		retry.Context(ctx),
	)
	return err
}

func (c *ClientImpl) setSpecifiedTalkOnAir(ctx context.Context, talkId int32) error {
	if err := c.client.GenerateAuth0Token(ctx,
		c.auth0Domain, c.auth0ClientId, c.auth0ClientSecret, c.auth0Audience,
	); err != nil {
		return xerrors.Errorf("message: %w", err)
	}

	if err := c.client.UpdateTalk(ctx, talkId, true); err != nil {
		return xerrors.Errorf("message: %w", err)
	}
	return nil
}

func (c *ClientImpl) SetNextTalkOnAir(ctx context.Context, eventAbbr string, trackId int32) error {
	// If Auth0Token has been expired, retry only once.
	err := retry.Do(
		func() (err error) {
			err = c.setNextTalkOnAir(ctx, eventAbbr, trackId)
			return
		},
		retry.RetryIf(func(err error) bool {
			return errors.As(err, &lib.ErrorUnauthorized{})
		}),
		retry.Attempts(1),
		retry.Context(ctx),
	)
	return err
}

func (c *ClientImpl) setNextTalkOnAir(ctx context.Context, eventAbbr string, trackId int32) error {
	if err := c.client.GenerateAuth0Token(ctx,
		c.auth0Domain, c.auth0ClientId, c.auth0ClientSecret, c.auth0Audience,
	); err != nil {
		return xerrors.Errorf("message: %w", err)
	}

	talks, err := c.client.ListTalks(ctx, eventAbbr, trackId)
	if err != nil {
		return xerrors.Errorf("message: %w", err)
	}

	var nextTalkId int32
	onAirFlag := false
	for idx, talk := range talks {
		if onAirFlag {
			nextTalkId = talk.ID
			break
		}
		if idx == len(talks)-1 {
			return xerrors.Errorf("message: Talks on specified track is the end. Next one is none.")
		}
		if talk.OnAir {
			onAirFlag = true
		}
	}

	if err := c.client.UpdateTalk(ctx, nextTalkId, true); err != nil {
		return xerrors.Errorf("message: %w", err)
	}
	return nil
}
