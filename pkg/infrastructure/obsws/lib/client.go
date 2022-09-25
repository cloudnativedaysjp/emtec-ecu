package lib

import (
	"github.com/andreykaipov/goobs"
	"golang.org/x/xerrors"
)

type PrimitiveClient struct {
	client *goobs.Client
}

func NewClient() ObsWsApi {
	return &PrimitiveClient{}
}

func (c *PrimitiveClient) GenerateClient(host, password string) error {
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

func (c *PrimitiveClient) MediaInputs() ObsWsMediaInputsApi {
	return c.client.MediaInputs
}

func (c *PrimitiveClient) Scenes() ObsWsScenesApi {
	return c.client.Scenes
}

func (c *PrimitiveClient) SceneItems() ObsWsSceneItemsApi {
	return c.client.SceneItems
}
