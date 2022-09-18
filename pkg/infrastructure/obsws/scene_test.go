//go:build obsws

// This test is for checking behavior.
// If you want to execute this test, you need to prepare OBS with obs-websocket &
// expose tcp/4455 port with no password.

package obsws

import (
	"context"
	"testing"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/testutils"
	"github.com/k0kubun/pp"
)

func TestObsWebSocketClient_ListScenes(t *testing.T) {
	var (
		host = testutils.Getenv(t, "OBS_HOST")

		ctx = context.Background()
	)

	t.Run("test", func(t *testing.T) {
		c := ObsWebSocketClient{
			host: host,
		}
		scenes, err := c.ListScenes(ctx)
		if err != nil {
			t.Errorf("ObsWebSocketClient.ListScenes() error = %v", err)
		}
		pp.Print(scenes)
	})
}

func TestObsWebSocketClient_MoveSceneToNext(t *testing.T) {
	var (
		host = testutils.Getenv(t, "OBS_HOST")

		ctx = context.Background()
	)

	t.Run("test", func(t *testing.T) {
		c := ObsWebSocketClient{
			host: host,
		}
		if err := c.MoveSceneToNext(ctx); err != nil {
			t.Errorf("ObsWebSocketClient.ListScenes() error = %v", err)
		}
	})
}
