package lib

import (
	"github.com/andreykaipov/goobs"
	"golang.org/x/xerrors"
)

type ObsWsClientWrapper struct {
	client *goobs.Client
}

func NewClient() ObsWsClient {
	return &ObsWsClientWrapper{}
}

func (c *ObsWsClientWrapper) GenerateClient(host, password string) error {
	if c.client != nil {
		return nil
	}
	client, err := goobs.New(host, goobs.WithPassword(password))
	if err != nil {
		return xerrors.Errorf("message: %w", err)
	}
	c.client = client
	return nil
}

func (c *ObsWsClientWrapper) MediaInputs() ObsWsMediaInputsApi {
	return c.client.MediaInputs
}

func (c *ObsWsClientWrapper) Scenes() ObsWsScenesApi {
	return c.client.Scenes
}

func (c *ObsWsClientWrapper) SceneItems() ObsWsSceneItemsApi {
	return c.client.SceneItems
}
