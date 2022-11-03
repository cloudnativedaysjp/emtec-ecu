package sharedmem

import (
	"fmt"

	"github.com/cloudnativedaysjp/emtec-ecu/pkg/model"
)

var _ ReaderIface = (*Reader)(nil)

type ReaderIface interface {
	DisableAutomation(trackId int32) (bool, error)
	Track(trackId int32) (*model.Track, error)
}

type Reader struct {
	UseStorageForDisableAutomation bool
	UseStorageForTrack             bool
}

func (r Reader) DisableAutomation(trackId int32) (bool, error) {
	if !r.UseStorageForDisableAutomation {
		return false, fmt.Errorf("UseStorageForDisableAutomation was false")
	}
	storageForDisableAutomationMutex.RLock()
	defer storageForDisableAutomationMutex.RUnlock()
	disabled, ok := storageForDisableAutomation[trackId]
	if !ok {
		return true, nil
	}
	return disabled, nil
}

func (r Reader) Track(trackId int32) (*model.Track, error) {
	if !r.UseStorageForTrack {
		return nil, fmt.Errorf("UseStorageForTrack was false")
	}
	storageForTrackMutex.RLock()
	defer storageForTrackMutex.RUnlock()
	track, ok := storageForTrack[trackId]
	if !ok {
		return nil, fmt.Errorf("trackId %d is not found", trackId)
	}
	return &track, nil
}
