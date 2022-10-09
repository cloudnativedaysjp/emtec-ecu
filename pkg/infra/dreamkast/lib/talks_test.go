//go:build dreamkast

// This test is for checking behavior.
// If you want to execute this test, you need to prepare real Dreamkast API.

package lib

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/testutils"
	"github.com/k0kubun/pp"
)

func TestPrimitiveClient_ListTalks(t *testing.T) {
	var (
		dkEndpointUrl = testutils.Getenv(t, "DK_ENDPOINT_URL")
		dkEventAbbr   = testutils.Getenv(t, "DK_EVENT_ABBR")
		trackIdStr    = testutils.Getenv(t, "DK_TRACK_ID_FOR_LISTUP")

		ctx = context.Background()
	)

	t.Run("test", func(t *testing.T) {
		u, err := url.Parse(dkEndpointUrl)
		if err != nil {
			t.Fatal(err)
		}
		c := &DreamkastClientImpl{
			client:        http.DefaultClient,
			dkEndpointUrl: *u,
		}

		trackId, err := strconv.Atoi(trackIdStr)
		if err != nil {
			t.Fatal(err)
		}

		got, err := c.ListTalks(ctx, dkEventAbbr, int32(trackId))
		if err != nil {
			t.Errorf("DreamkastClientImpl.ListTalks() error = %v", err)
			return
		}
		pp.Print(got)
	})
}

func TestPrimitiveClient_UpdateTalks(t *testing.T) {
	var (
		dkEndpointUrl     = testutils.Getenv(t, "DK_ENDPOINT_URL")
		talkIdStr         = testutils.Getenv(t, "DK_TALK_ID_FOR_SETTING_ONAIR")
		auth0Domain       = testutils.Getenv(t, "AUTH0_DOMAIN")
		auth0ClientId     = testutils.Getenv(t, "AUTH0_CLIENT_ID")
		auth0ClientSecret = testutils.Getenv(t, "AUTH0_CLIENT_SECRET")
		auth0Audience     = testutils.Getenv(t, "AUTH0_AUDIENCE")

		ctx = context.Background()
	)

	t.Run("test", func(t *testing.T) {
		u, err := url.Parse(dkEndpointUrl)
		if err != nil {
			t.Fatal(err)
		}
		c := &DreamkastClientImpl{
			client:        http.DefaultClient,
			dkEndpointUrl: *u,
		}
		if err := c.GenerateAuth0Token(auth0Domain,
			auth0ClientId, auth0ClientSecret, auth0Audience); err != nil {
			t.Fatal(err)
		}
		talkId, err := strconv.Atoi(talkIdStr)
		if err != nil {
			t.Fatal(err)
		}

		if err := c.UpdateTalks(ctx, int32(talkId), true); err != nil {
			t.Errorf("DreamkastClientImpl.UpdateTalks() error = %v", err)
		}
	})
}
