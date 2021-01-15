package ha

import (
	"sync"
	"time"
)

type State struct {
	leader          string
	leaseStart      time.Time
	leaseUntil      time.Time
	minTimeSeen     time.Time
	maxTimeSeen     time.Time // max data time seen by any instance
	maxTimeInstance string    // the instance name that’s seen the maxtime
	_mu             sync.RWMutex
}

type StateView struct {
	leader          string
	leaseStart      time.Time
	leaseUntil      time.Time
	maxTimeSeen     time.Time
	maxTimeInstance string
}

func (h *State) updateStateFromDB(latestState *haLockState, maxT time.Time, replicaName string) {
	h._mu.Lock()
	defer h._mu.Unlock()
	if h.maxTimeSeen.IsZero() {
		h.maxTimeSeen = maxT
		h.maxTimeInstance = replicaName
	}
	h.leader = latestState.leader
	h.leaseStart = latestState.leaseStart
	h.leaseUntil = latestState.leaseUntil
}

func (h *State) updateState(currentReplica string, currentMaxT, currentMinT time.Time) {
	h._mu.Lock()
	defer h._mu.Unlock()
	if currentMaxT.After(h.maxTimeSeen) && currentMinT.After(h.minTimeSeen) {
		h.maxTimeInstance = currentReplica
		h.maxTimeSeen = currentMaxT
		h.minTimeSeen = currentMinT
	}
}

func (h *State) clone() *StateView {
	h._mu.RLock()
	defer h._mu.RUnlock()
	return &StateView{
		leader:          h.leader,
		leaseStart:      h.leaseStart,
		leaseUntil:      h.leaseUntil,
		maxTimeSeen:     h.maxTimeSeen,
		maxTimeInstance: h.maxTimeInstance,
	}
}