package sharedmem

import (
	"fmt"

	"github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"
)

type WriterIface interface {
	SetDisableAutomation(trackId int32, disabled bool) error
	SetTalks(trackId int32, talk model.Talk) error
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

func (w Writer) SetTalks(talks []model.Talk) error {
	if !w.UseStorageForTalks {
		return fmt.Errorf("UseStorageForTalks was false")
	}
	storageForTalksMutex.Lock()
	defer storageForTalksMutex.Unlock()
	for _, talk := range talks {
		storageForTalks[talk.TrackId] = talk
	}
	return nil
}
