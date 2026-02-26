package main

import (
	"encoding/json"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WaybarOutput struct {
	Text    string `json:"text"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
}

type BenchmarkResult struct {
	AgentID   string  `json:"agent_id"`
	LatencyMS float32 `json:"latency_ms"`
}

type MeshResults struct {
	AgentsActive int               `json:"agents_active"`
	Benchmarks   []BenchmarkResult `json:"benchmarks"`
}

func main() {
	// 1. Check Controller
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		printError("Vextra: Offline")
		return
	}
	defer conn.Close()

	// 2. Read latest benchmarks
	data, _ := os.ReadFile("benchmark_results.json")
	var mesh MeshResults
	json.Unmarshal(data, &mesh)

	var totalLatency float32
	validCount := 0
	for _, b := range mesh.Benchmarks {
		if b.LatencyMS > 0 {
			totalLatency += b.LatencyMS
			validCount++
		}
	}

	avg := float32(0)
	if validCount > 0 {
		avg = totalLatency / float32(validCount)
	}

	output := WaybarOutput{
		Text:    fmt.Sprintf("󱐋 Mesh: %d Agents", mesh.AgentsActive),
		Tooltip: fmt.Sprintf("Avg Latency: %.1fms\nTasks Validated: %d", avg, validCount),
		Class:   "online",
	}

	if avg > 100 {
		output.Class = "degraded"
	}

	json.NewEncoder(os.Stdout).Encode(output)
}

func printError(msg string) {
	json.NewEncoder(os.Stdout).Encode(WaybarOutput{
		Text:  msg,
		Class: "offline",
	})
}
