package obsws

import (
	"context"

	"github.com/andreykaipov/goobs"
)

type ObsWebSocketApi interface {
	GetHost() string
	ListScenes(ctx context.Context) ([]Scene, error)
	MoveSceneToNext(ctx context.Context) error
}

type ObsWebSocketClient struct {
	client *goobs.Client

	host     string
	password string
}

func NewObsWebSocketClient(host, password string) (ObsWebSocketApi, error) {
	c := &ObsWebSocketClient{host: host, password: password}
	if err := c.setClient(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *ObsWebSocketClient) GetHost() string {
	return c.host
}

func (c *ObsWebSocketClient) setClient() error {
	if c.client != nil {
		return nil
	}
	client, err := goobs.New(c.host, goobs.WithPassword(c.password))
	if err != nil {
		return err
	}
	c.client = client
	return nil
}
