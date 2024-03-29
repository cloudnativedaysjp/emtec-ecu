package obsws

import (
	"context"
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/mediainputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/requests/scenes"
	"golang.org/x/xerrors"

	"github.com/cloudnativedaysjp/emtec-ecu/pkg/infra/obsws/lib"
	"github.com/cloudnativedaysjp/emtec-ecu/pkg/utils"
)

type Client interface {
	GetHost() string
	ListScenes(ctx context.Context) ([]Scene, error)
	MoveSceneToNext(ctx context.Context) error
	GetRemainingTimeOnCurrentScene(ctx context.Context) (*DurationAndCursor, error)
}

type ClientImpl struct {
	client lib.ObsWsClient

	host     string
	password string
}

func NewObsWebSocketClient(host, password string) (Client, error) {
	c := lib.NewClient()
	if err := c.GenerateClient(host, password); err != nil {
		return nil, err
	}
	return &ClientImpl{c, host, password}, nil
}

func (c *ClientImpl) GetHost() string {
	return c.host
}

type Scene struct {
	Name             string
	SceneIndex       int
	IsCurrentProgram bool
}

// ListScenes is output list of scenes. It is sorted order by as shown in OBS.
func (c *ClientImpl) ListScenes(ctx context.Context) ([]Scene, error) {
	if err := c.client.GenerateClient(c.host, c.password); err != nil {
		return nil, err
	}

	resp, err := c.client.Scenes().GetSceneList()
	if err != nil {
		c.client = nil
		return nil, err
	}

	var scenes []Scene
	for _, s := range resp.Scenes {
		scene := Scene{Name: s.SceneName, SceneIndex: s.SceneIndex}
		if s.SceneName == resp.CurrentProgramSceneName {
			scene.IsCurrentProgram = true
		}
		scenes = append(scenes, scene)
	}

	// reverse
	for i := 0; i < len(scenes)/2; i++ {
		scenes[i], scenes[len(scenes)-i-1] = scenes[len(scenes)-i-1], scenes[i]
	}

	return scenes, nil
}

func (c *ClientImpl) MoveSceneToNext(ctx context.Context) error {
	logger := utils.GetLogger(ctx)
	if err := c.client.GenerateClient(c.host, c.password); err != nil {
		return err
	}

	_scenes, err := c.ListScenes(ctx)
	if err != nil {
		return err
	}

	var nextSceneName string
	var currentProgramFlag, nextProgramFlag bool
	for _, scene := range _scenes {
		if scene.IsCurrentProgram {
			currentProgramFlag = true
			continue
		}
		if currentProgramFlag {
			currentProgramFlag = false
			nextProgramFlag = true
			nextSceneName = scene.Name
			break
		}
	}
	if currentProgramFlag {
		logger.Info("current scene is the last scene.")
		return nil
	} else if !nextProgramFlag {
		return fmt.Errorf("CurrentProgram is nothing")
	}

	if _, err := c.client.Scenes().SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneName: nextSceneName,
	}); err != nil {
		c.client = nil
		return err
	}
	return nil
}

type DurationAndCursor struct {
	DurationMilliSecond float64
	CursorMilliSecond   float64
}

func (c *ClientImpl) GetRemainingTimeOnCurrentScene(ctx context.Context) (*DurationAndCursor, error) {
	_ = utils.GetLogger(ctx)
	if err := c.client.GenerateClient(c.host, c.password); err != nil {
		return nil, err
	}

	listScenesResp, err := c.client.Scenes().GetSceneList()
	if err != nil {
		c.client = nil
		return nil, err
	}

	var currentSceneName string
	for _, s := range listScenesResp.Scenes {
		if s.SceneName == listScenesResp.CurrentProgramSceneName {
			currentSceneName = s.SceneName
		}
	}
	if currentSceneName == "" {
		return nil, xerrors.Errorf("message: %w", fmt.Errorf("CurrentProgramSceneName is none in Scenes"))
	}

	listSceneItemsResp, err := c.client.SceneItems().GetSceneItemList(
		&sceneitems.GetSceneItemListParams{SceneName: currentSceneName})
	if err != nil {
		c.client = nil
		return nil, err
	}

	if len(listSceneItemsResp.SceneItems) == 0 {
		return nil, xerrors.Errorf("message: %w", fmt.Errorf("Source is none in Scene '%s'", currentSceneName))
	}
	// TODO: 決め打ちであるのを直す
	mediaInputName := listSceneItemsResp.SceneItems[0].SourceName

	resp, err := c.client.MediaInputs().GetMediaInputStatus(&mediainputs.GetMediaInputStatusParams{InputName: mediaInputName})
	if err != nil {
		c.client = nil
		return nil, err
	}
	if resp.MediaState != "OBS_MEDIA_STATE_PLAYING" {
		return nil, xerrors.Errorf("message: %w", fmt.Errorf("state of MediaInput '%s' isn't OBS_MEDIA_STATE_PLAYING: now is %s",
			mediaInputName, resp.MediaState))
	} else if resp.MediaDuration == 0 {
		return nil, xerrors.Errorf("message: %w", fmt.Errorf("duration of MediaInput '%s' is zero", mediaInputName))
	}
	return &DurationAndCursor{resp.MediaDuration, resp.MediaCursor}, nil
}
