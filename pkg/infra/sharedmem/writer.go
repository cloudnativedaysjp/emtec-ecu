package sharedmem

import (
	"fmt"

	"github.com/cloudnativedaysjp/emtec-ecu/pkg/model"
)

var _ WriterIface = (*Writer)(nil)

type WriterIface interface {
	SetDisableAutomation(trackId int32, disabled bool) error
	SetTrack(track model.Track) error
}

type Writer struct {
	UseStorageForDisableAutomation bool
	UseStorageForTrack             bool
}

func (w Writer) SetDisableAutomation(trackId int32, disabled bool) error {
	if !w.UseStorageForDisableAutomation {
		return fmt.Errorf("UseStorageForDisableAutomation was false")
	}
	storageForDisableAutomationMutex.Lock()
	defer storageForDisableAutomationMutex.Unlock()
	storageForDisableAutomation[trackId] = disabled
	return nil
}

func (w Writer) SetTrack(track model.Track) error {
	if !w.UseStorageForTrack {
		return fmt.Errorf("UseStorageForTrack was false")
	}
	storageForTrackMutex.Lock()
	defer storageForTrackMutex.Unlock()
	storageForTrack[track.Id] = track
	return nil
}
