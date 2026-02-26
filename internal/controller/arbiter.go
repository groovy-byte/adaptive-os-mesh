package controller

import (
	"log"
	"sync"
	"time"

	pb "github.com/groovy-byte/agent-mesh-core/proto"
)

const LockTTL = 30 * time.Second

// Arbiter manages global strategic locks and state recovery
type Arbiter struct {
	mu            sync.Mutex
	strategicLock string // Agent ID that currently holds the Strategic planning lock
	lockTime      time.Time
	lastStates    map[string]*pb.AgentAction
}

func NewArbiter() *Arbiter {
	return &Arbiter{
		lastStates: make(map[string]*pb.AgentAction),
	}
}

// RequestStrategicLock implements the Counterbalance mechanism to prevent 'Too many bosses'
func (a *Arbiter) RequestStrategicLock(agentID string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check for stale lock
	if a.strategicLock != "" && a.strategicLock != agentID {
		if time.Since(a.lockTime) > LockTTL {
			log.Printf("[Arbiter] ⚠️ Reclaiming stale lock from %s (Expired after %v)", a.strategicLock, LockTTL)
			a.strategicLock = ""
		}
	}

	if a.strategicLock == "" || a.strategicLock == agentID {
		a.strategicLock = agentID
		a.lockTime = time.Now()
		log.Printf("[Arbiter] 🔑 Strategic Lock granted to %s", agentID)
		return true
	}

	log.Printf("[Arbiter] ⚠️ Strategic Lock denied to %s (Held by %s)", agentID, a.strategicLock)
	return false
}

// ReleaseLock frees the strategic reasoning path
func (a *Arbiter) ReleaseLock(agentID string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.strategicLock == agentID {
		a.strategicLock = ""
		log.Printf("[Arbiter] 🔓 Strategic Lock released by %s", agentID)
	}
}

// SaveState records the last known operational state for reconstitution (Task 4.2)
func (a *Arbiter) SaveState(state *pb.AgentAction) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lastStates[state.AgentId] = state
}

func (a *Arbiter) GetState(agentID string) (*pb.AgentAction, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	state, ok := a.lastStates[agentID]
	return state, ok
}
