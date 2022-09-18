//go:build dreamkast

// This test is for checking behavior.
// If you want to execute this test, you need to prepare real Dreamkast API.

package lib

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/testutils"
	"github.com/k0kubun/pp"
)

func TestPrimitiveClient_ListTracks(t *testing.T) {
	var (
		dkEndpointUrl = testutils.Getenv(t, "DK_ENDPOINT_URL")
		dkEventAbbr   = testutils.Getenv(t, "DK_EVENT_ABBR")

		ctx = context.Background()
	)

	t.Run("test", func(t *testing.T) {
		u, err := url.Parse(dkEndpointUrl)
		if err != nil {
			t.Fatal(err)
		}
		c := &PrimitiveClient{
			client:        http.DefaultClient,
			dkEndpointUrl: *u,
		}

		got, err := c.ListTracks(ctx, dkEventAbbr)
		if err != nil {
			t.Errorf("PrimitiveClient.ListTracks() error = %v", err)
			return
		}
		pp.Print(got)
	})
}
