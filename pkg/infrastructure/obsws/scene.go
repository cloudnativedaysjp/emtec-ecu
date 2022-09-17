package obsws

import (
	"context"
	"fmt"

	"github.com/andreykaipov/goobs/api/requests/scenes"
	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/utils"
)

type Scene struct {
	Name             string
	SceneIndex       int
	IsCurrentProgram bool
}

// ListScenes is output list of scenes. It is sorted order by as shown in OBS.
func (c *ObsWebSocketClient) ListScenes(ctx context.Context) ([]Scene, error) {
	_ = utils.GetLogger(ctx)
	if err := c.setClient(); err != nil {
		return nil, err
	}

	resp, err := c.client.Scenes.GetSceneList()
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

func (c *ObsWebSocketClient) MoveSceneToNext(ctx context.Context) error {
	logger := utils.GetLogger(ctx)
	if err := c.setClient(); err != nil {
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
		logger.Info("current scene is tha last scene.")
		nextSceneName = _scenes[0].Name
	} else if !nextProgramFlag {
		return fmt.Errorf("CurrentProgram is nothing")
	}

	if _, err := c.client.Scenes.SetCurrentProgramScene(&scenes.SetCurrentProgramSceneParams{
		SceneName: nextSceneName,
	}); err != nil {
		c.client = nil
		return err
	}
	return nil
}
