package controller

import (
	"log"
	"strings"
	pb "github.com/groovy-byte/agent-mesh-core/proto"
)

// RoleSwitcher evaluates system health and suggests agent role transitions
type RoleSwitcher struct {
	memoryThreshold uint64
	cpuThreshold    float64
}

func NewRoleSwitcher() *RoleSwitcher {
	return &RoleSwitcher{
		memoryThreshold: 800 * 1024 * 1024,
		cpuThreshold:    70.0,
	}
}

// EvaluateTransition implements Difficulty-Aware Topology Selection
func (s *RoleSwitcher) EvaluateTransition(info *AgentInfo, currentLoad *pb.OSResources, intent string) pb.AgentRole {
	intent = strings.ToUpper(intent)
	
	// 1. Peak-Aware Predictive Control (Alignment in Time)
	if strings.Contains(intent, "LONG_HORIZON") || strings.Contains(intent, "COMPILATION") {
		if info.Role == pb.AgentRole_STRATEGIC {
			log.Printf("[Role] ⏳ Peak-Aware Throttling: Scaling down %s for long-horizon task stability", info.ID)
			return pb.AgentRole_OPERATIONAL
		}
	}

	// 2. Difficulty-Aware Promotion (Difficulty-Aware Orchestration)
	if strings.Contains(intent, "STRATEGIC_REASONING") && info.Role == pb.AgentRole_OPERATIONAL {
		log.Printf("[Role] 🧠 High Difficulty detected: Suggesting STRATEGIC promotion for %s", info.ID)
		return pb.AgentRole_STRATEGIC
	}

	// 3. Reactive Fallback
	if currentLoad != nil && (currentLoad.MemoryUsedBytes > s.memoryThreshold || currentLoad.CpuUsagePercent > s.cpuThreshold) {
		return pb.AgentRole_OPERATIONAL
	}

	return info.Role
}
