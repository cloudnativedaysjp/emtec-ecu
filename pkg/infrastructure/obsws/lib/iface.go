package lib

import (
	"github.com/andreykaipov/goobs/api/requests/mediainputs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/andreykaipov/goobs/api/requests/scenes"
)

type ObsWsApi interface {
	GenerateClient(host, password string) error
	MediaInputs() ObsWsMediaInputsApi
	Scenes() ObsWsScenesApi
	SceneItems() ObsWsSceneItemsApi
}

type ObsWsMediaInputsApi interface {
	GetMediaInputStatus(params *mediainputs.GetMediaInputStatusParams) (*mediainputs.GetMediaInputStatusResponse, error)
}

type ObsWsScenesApi interface {
	GetSceneList(paramss ...*scenes.GetSceneListParams) (*scenes.GetSceneListResponse, error)
	SetCurrentProgramScene(params *scenes.SetCurrentProgramSceneParams) (*scenes.SetCurrentProgramSceneResponse, error)
}

type ObsWsSceneItemsApi interface {
	GetSceneItemList(params *sceneitems.GetSceneItemListParams) (*sceneitems.GetSceneItemListResponse, error)
}
