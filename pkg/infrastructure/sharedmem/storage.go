package sharedmem

import (
	"sync"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

// key is TrackId
type StorageForDisableAutomation map[int32]bool

// key is TrackId
type StorageForTalks map[int32]model.Talk

var (
	storageForDisableAutomation      = make(map[int32]bool)
	storageForDisableAutomationMutex sync.RWMutex
	storageForTalks                  = make(map[int32]model.Talk)
	storageForTalksMutex             sync.RWMutex
)
