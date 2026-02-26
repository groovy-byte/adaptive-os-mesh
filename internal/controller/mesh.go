package controller

import (
	"fmt"
	"log"
	"sync"

	pb "github.com/groovy-byte/agent-mesh-core/proto"
)

// AgentInfo represents an active agent in the mesh
type AgentInfo struct {
	ID            string
	Role          pb.AgentRole
	Capabilities  []string
	Neighbors     []string
	UtilityScore  float64 // DSBO Utility tracking
	TotalLatency  float32
	TotalTokens   uint32
	RequestCount  uint32
	ToolCalls     uint32
	FailedTasks   []string // Last 5 failed tasks
	MaxThroughput float32  // NEW: Peak GB/s tracked
}

// MeshRegistry manages the active agents and their communication neighborhoods
type MeshRegistry struct {
	mu                 sync.RWMutex
	agents             map[string]*AgentInfo
	contributionMatrix map[string]map[string]float64 // Source -> {Target: Score}
}

func NewMeshRegistry() *MeshRegistry {
	return &MeshRegistry{
		agents:             make(map[string]*AgentInfo),
		contributionMatrix: make(map[string]map[string]float64),
	}
}

// RecordContribution tracks the Value of Contribution (VoC) between agents
func (r *MeshRegistry) RecordContribution(sourceID, targetID string, score float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.contributionMatrix[sourceID]; !ok {
		r.contributionMatrix[sourceID] = make(map[string]float64)
	}
	// Aggregate influence with normalization clamp
	r.contributionMatrix[sourceID][targetID] += score
	if r.contributionMatrix[sourceID][targetID] > 1.0 {
		r.contributionMatrix[sourceID][targetID] = 1.0
	}
}

// GetContributionDetail returns how a specific agent influenced others
func (r *MeshRegistry) GetContributionDetail(id string) map[string]float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	detail := make(map[string]float64)
	if targets, ok := r.contributionMatrix[id]; ok {
		for target, score := range targets {
			detail[target] = score
		}
	}
	return detail
}

// GetAgent returns a thread-safe copy of an agent's info
func (r *MeshRegistry) GetAgent(id string) (AgentInfo, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.agents[id]
	if !ok {
		return AgentInfo{}, false
	}
	
	// Deep copy slices to prevent external modification
	failedTasks := make([]string, len(info.FailedTasks))
	copy(failedTasks, info.FailedTasks)
	
	caps := make([]string, len(info.Capabilities))
	copy(caps, info.Capabilities)

	neighbors := make([]string, len(info.Neighbors))
	copy(neighbors, info.Neighbors)

	return AgentInfo{
		ID:           info.ID,
		Role:         info.Role,
		Capabilities: caps,
		Neighbors:    neighbors,
		UtilityScore: info.UtilityScore,
		TotalLatency: info.TotalLatency,
		TotalTokens:  info.TotalTokens,
		RequestCount: info.RequestCount,
		ToolCalls:    info.ToolCalls,
		FailedTasks:  failedTasks,
	}, true
}

// RegisterAgent handles the "One-Hop" handshake logic
func (r *MeshRegistry) RegisterAgent(req *pb.HandshakeRequest) (*pb.HandshakeResponse, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	log.Printf("[Mesh] Handshake received from Agent: %s", req.AgentId)

	// DSBO: Initial neighborhood selection
	neighbors := []string{}
	for id := range r.agents {
		if id != req.AgentId {
			neighbors = append(neighbors, id)
			if len(neighbors) >= 2 { // Increased limit for evaluation
				break
			}
		}
	}

	agent := &AgentInfo{
		ID:           req.AgentId,
		Role:         req.InitialRole,
		Capabilities: req.Capabilities,
		Neighbors:    neighbors,
		UtilityScore: 1.0, // Initial base utility
		TotalLatency: 0,
		TotalTokens:  0,
		RequestCount: 0,
		ToolCalls:    0,
		FailedTasks:  []string{},
	}
	r.agents[req.AgentId] = agent

	return &pb.HandshakeResponse{
		SessionId: fmt.Sprintf("mesh_sess_%s", req.AgentId),
		Approved:  true,
		ResourceLimits: &pb.OSResources{
			CpuUsagePercent:  75.0,
			MemoryTotalBytes: 512 * 1024 * 1024,
		},
	}, nil
}

// RecordTaskResult logs the outcome of an agent operation
func (r *MeshRegistry) RecordTaskResult(id string, success bool, taskName string, toolCount uint32) {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, ok := r.agents[id]
	if !ok {
		return
	}

	agent.ToolCalls += toolCount
	if !success {
		// Circular buffer of last 5 failures
		agent.FailedTasks = append([]string{taskName}, agent.FailedTasks...)
		if len(agent.FailedTasks) > 5 {
			agent.FailedTasks = agent.FailedTasks[:5]
		}
	}
}

// RecordMetrics updates the performance tracking for an agent
func (r *MeshRegistry) RecordMetrics(id string, latency float32, tokens uint32, throughput float32) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if agent, ok := r.agents[id]; ok {
		agent.TotalLatency += latency
		agent.TotalTokens += tokens
		agent.RequestCount++
		if throughput > agent.MaxThroughput {
			agent.MaxThroughput = throughput
		}
	}
}

// AgentStats represents performance metrics for an agent
type AgentStats struct {
	ID         string
	AvgLatency float32
	Requests   uint32
	Tokens     uint32
}

func (r *MeshRegistry) GetStatsSummary() []AgentStats {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := make([]AgentStats, 0, len(r.agents))
	for id, agent := range r.agents {
		avg := float32(0)
		if agent.RequestCount > 0 {
			avg = agent.TotalLatency / float32(agent.RequestCount)
		}
		stats = append(stats, AgentStats{
			ID:         id,
			AvgLatency: avg,
			Requests:   agent.RequestCount,
			Tokens:     agent.TotalTokens,
		})
	}
	return stats
}

// ReevaluateNeighbors implements DSBO and Proactive Healing signals
func (r *MeshRegistry) ReevaluateNeighbors(agentID string, novelContext bool, arbiter *Arbiter) {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, ok := r.agents[agentID]
	if !ok {
		return
	}

	// Update Utility Score
	if novelContext {
		agent.UtilityScore += 0.1
	} else {
		agent.UtilityScore -= 0.05
	}

	// Proactive Healing Signal: If utility drops rapidly, trigger state snapshot
	if agent.UtilityScore < 0.7 {
		log.Printf("[Mesh] 🩹 Utility Drop: Notifying Arbiter to snapshot %s for proactive healing", agentID)
		// Signal logic handled in the main controller loop using the arbiter
	}

	// Pruning Logic
	if agent.UtilityScore < 0.5 && len(agent.Neighbors) > 1 {
		log.Printf("[DSBO] Pruning communication path for %s due to low novelty", agentID)
		agent.Neighbors = agent.Neighbors[:len(agent.Neighbors)-1]
		agent.UtilityScore = 0.8 
	}
}

func (r *MeshRegistry) UpdateRole(id string, role pb.AgentRole) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if agent, ok := r.agents[id]; ok {
		agent.Role = role
	}
}
