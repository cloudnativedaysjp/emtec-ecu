package dreamkast

import (
	"context"
	"errors"

	"github.com/avast/retry-go"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/infrastructure/dreamkast/lib"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

type ClientIface interface {
	ListTalks(ctx context.Context) (model.Talks, error)
	SetSpecifiedTalkOnAir(ctx context.Context, talkId int32) error
	SetNextTalkOnAir(ctx context.Context, trackId int32) error
}

type Client struct {
	client lib.DreamkastApi

	eventAbbr         string
	auth0Domain       string
	auth0ClientId     string
	auth0ClientSecret string
	auth0Audience     string
}

func NewClient(eventAbbr, dkEndpointUrl string,
	auth0Domain, auth0ClientId, auth0ClientSecret, auth0Audience string,
) (ClientIface, error) {
	c, err := lib.NewClient(dkEndpointUrl)
	if err != nil {
		return nil, err
	}
	return &Client{
		c, eventAbbr, auth0Domain, auth0ClientId, auth0ClientSecret, auth0Audience}, nil
}

func (c *Client) ListTalks(ctx context.Context) (model.Talks, error) {
	tracks, err := c.client.ListTracks(ctx, c.eventAbbr)
	if err != nil {
		return nil, xerrors.Errorf("message: %w", err)
	}
	var result model.Talks
	for _, track := range tracks {
		talks, err := c.client.ListTalks(ctx, c.eventAbbr, track.ID)
		if err != nil {
			return nil, xerrors.Errorf("message: %w", err)
		}
		for _, talk := range talks {
			t := model.Talk{
				Id:        talk.ID,
				TalkName:  talk.Title,
				TrackId:   track.ID,
				TrackName: track.Name,
				EventAbbr: c.eventAbbr,
				// TODO (https://github.com/cloudnativedaysjp/dreamkast/issues/1490)
				//Type         TalkType
			}
			for _, speaker := range talk.Speakers {
				t.SpeakerNames = append(t.SpeakerNames, speaker.Name)
			}

			// TODO (#11)
			// response is as below, so calcurate YYYY-MM-DDThh:mm:ss from these fields
			// - "conferenceDayDate": "2022-08-05"
			// - "actualStartTime": "2000-01-01T13:05:00.000+09:00"
			t.StartAt = talk.StartTime
			t.EndAt = talk.EndTime

			result = append(result, t)
		}
	}
	result.FillCommercial()
	return result, nil
}

func (c *Client) SetSpecifiedTalkOnAir(ctx context.Context, talkId int32) error {
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

func (c *Client) setSpecifiedTalkOnAir(ctx context.Context, talkId int32) error {
	if err := c.client.GenerateAuth0Token(
		c.auth0Domain, c.auth0ClientId, c.auth0ClientSecret, c.auth0Audience,
	); err != nil {
		return xerrors.Errorf("message: %w", err)
	}

	if err := c.client.UpdateTalks(ctx, talkId, true); err != nil {
		return xerrors.Errorf("message: %w", err)
	}
	return nil
}

func (c *Client) SetNextTalkOnAir(ctx context.Context, trackId int32) error {
	// If Auth0Token has been expired, retry only once.
	err := retry.Do(
		func() (err error) {
			err = c.setNextTalkOnAir(ctx, trackId)
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

func (c *Client) setNextTalkOnAir(ctx context.Context, trackId int32) error {
	if err := c.client.GenerateAuth0Token(
		c.auth0Domain, c.auth0ClientId, c.auth0ClientSecret, c.auth0Audience,
	); err != nil {
		return xerrors.Errorf("message: %w", err)
	}

	talks, err := c.client.ListTalks(ctx, c.eventAbbr, trackId)
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

	if err := c.client.UpdateTalks(ctx, nextTalkId, true); err != nil {
		return xerrors.Errorf("message: %w", err)
	}
	return nil
}