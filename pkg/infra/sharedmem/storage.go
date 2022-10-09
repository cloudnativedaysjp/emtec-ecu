package sharedmem

import (
	"sync"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

// key is TrackId
type StorageForDisableAutomation map[int32]bool

// key is TrackId
type StorageForTrack map[int32]model.Track

var (
	storageForDisableAutomation      = make(StorageForDisableAutomation)
	storageForDisableAutomationMutex sync.RWMutex
	storageForTrack                  = make(StorageForTrack)
	storageForTrackMutex             sync.RWMutex
)
