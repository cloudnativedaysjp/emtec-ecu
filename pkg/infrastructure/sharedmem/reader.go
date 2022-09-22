package sharedmem

import (
	"fmt"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

type ReaderIface interface {
	DisableAutomation(trackId int32) (bool, error)
	Talks(trackId int32) (model.Talk, error)
}

type Reader struct {
	UseStorageForDisableAutomation bool
	UseStorageForTalks             bool
}

func (w Reader) DisableAutomation(trackId int32) (bool, error) {
	if !w.UseStorageForDisableAutomation {
		return false, fmt.Errorf("UseStorageForDisableAutomation was false")
	}
	storageForDisableAutomationMutex.RLock()
	defer storageForDisableAutomationMutex.RUnlock()
	disabled, ok := storageForDisableAutomation[trackId]
	if !ok {
		return false, fmt.Errorf("trackId %d is not found", trackId)
	}
	return disabled, nil
}

func (w Reader) Talks(trackId int32) (model.Talk, error) {
	if !w.UseStorageForTalks {
		return model.Talk{}, fmt.Errorf("UseStorageForTalks was false")
	}
	storageForTalksMutex.RLock()
	defer storageForTalksMutex.RUnlock()
	talk, ok := storageForTalks[trackId]
	if !ok {
		return model.Talk{}, fmt.Errorf("trackId %d is not found", trackId)
	}
	return talk, nil
}
