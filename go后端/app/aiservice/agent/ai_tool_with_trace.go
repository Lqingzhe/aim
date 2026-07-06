package agent

import (
	"sync"
	"time"
)

type TraceWithUser struct {
	TraceID string
	Epoch   int64
	Lock    sync.Mutex
}
type TraceWithUserManager struct {
	Info        map[int64]*TraceWithUser
	mu          sync.RWMutex
	GlobalEpoch int64
}

func NewTraceWithUserManager() *TraceWithUserManager {
	return &TraceWithUserManager{
		Info: make(map[int64]*TraceWithUser),
	}
}

func (t *TraceWithUserManager) deleteTraceWithUser(userID int64, Epoch int64) {
	select {
	case <-time.After(1 * time.Second):
	}
	t.mu.Lock()
	info, ok := t.Info[userID]
	if ok {
		info.Lock.Lock()
		if info.Epoch == Epoch {
			delete(t.Info, userID)
		}
		info.Lock.Unlock()
	}
	t.mu.Unlock()
}
func (t *TraceWithUserManager) SetTraceID(userID int64, traceID string) {
	t.mu.Lock()
	t.GlobalEpoch++
	globalEpoch := t.GlobalEpoch
	t.mu.Unlock()
	t.mu.RLock()
	info, ok := t.Info[userID]
	t.mu.RUnlock()
	if ok {
		info.Lock.Lock()
		info.Epoch = globalEpoch
		info.TraceID = traceID
	} else {
		t.mu.Lock()
		info, ok = t.Info[userID]
		if ok {
			t.mu.Unlock()
			info.Lock.Lock()
			info.Epoch = globalEpoch
			info.TraceID = traceID
		} else {
			t.Info[userID] = &TraceWithUser{TraceID: traceID, Epoch: globalEpoch}
			t.Info[userID].Lock.Lock()
			t.mu.Unlock()
		}
	}
}
func (t *TraceWithUserManager) GetTraceID(userID int64) string {
	t.mu.RLock()
	traceID := t.Info[userID].TraceID
	t.mu.RUnlock()
	return traceID
}
func (t *TraceWithUserManager) ReleaseTrace(userID int64) {
	t.mu.RLock()
	info, ok := t.Info[userID]
	if !ok {
		t.mu.RUnlock()
		return
	}
	go t.deleteTraceWithUser(userID, info.Epoch)
	info.Lock.Unlock()
	t.mu.RUnlock()
}
