package sharedmem

import "github.com/cloudnativedaysjp/cnd-operation-server/pkg/model"

type DebuggerIface interface {
	ListAutomation() map[int32]bool
	ListTalks() map[int32]model.Talks
}

type Debugger struct{}

func (d Debugger) ListAutomation() map[int32]bool {
	result := make(map[int32]bool)
	storageForDisableAutomationMutex.RLock()
	defer storageForDisableAutomationMutex.RUnlock()
	for k, v := range storageForDisableAutomation {
		result[k] = v
	}
	return result
}

func (d Debugger) ListTalks() map[int32]model.Talks {
	result := make(map[int32]model.Talks)
	storageForTalksMutex.RLock()
	defer storageForTalksMutex.RUnlock()
	for k, v := range storageForTalks {
		result[k] = v
	}
	return result
}
