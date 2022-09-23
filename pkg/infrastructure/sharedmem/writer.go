package sharedmem

import (
	"fmt"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

type WriterIface interface {
	SetDisableAutomation(trackId int32, disabled bool) error
	SetTalks(ttrackId int32, alks model.Talks) error
}

type Writer struct {
	UseStorageForDisableAutomation bool
	UseStorageForTalks             bool
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

func (w Writer) SetTalks(trackId int32, talks model.Talks) error {
	if !w.UseStorageForTalks {
		return fmt.Errorf("UseStorageForTalks was false")
	}
	storageForTalksMutex.Lock()
	defer storageForTalksMutex.Unlock()
	storageForTalks[talks.GetCurrentTalk().TrackId] = talks
	return nil
}
