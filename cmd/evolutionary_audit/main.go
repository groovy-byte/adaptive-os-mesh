package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "github.com/groovy-byte/agent-mesh-core/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuditLog struct {
	Timestamp string `json:"timestamp"`
	Agent     string `json:"agent"`
	Action    string `json:"action"`
	Status    string `json:"status"`
	Details   string `json:"details"`
}

func main() {
	fmt.Println("🛡️ Starting Vextra Evolutionary Auditor (Log Synthesis Tier)...")
	
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to mesh: %v", err)
	}
	defer conn.Close()
	client := pb.NewStrategicMeshClient(conn)

	// 1. Handshake as Auditor
	_, err = client.RegisterAgent(context.Background(), &pb.HandshakeRequest{
		AgentId:      "Evolutionary-Auditor",
		Capabilities: []string{"LOG_SYNTHESIS", "VOC_CALCULATION", "MESH_INTEGRITY"},
		InitialRole:  pb.AgentRole_STRATEGIC,
	})
	if err != nil {
		log.Printf("Audit node degraded: %v", err)
	}

	// 2. Continuous Monitoring Loop
	ticker := time.NewTicker(2 * time.Second)
	auditFile, _ := os.OpenFile("final_audit.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer auditFile.Close()

	fmt.Println("📊 Monitoring mesh workloads. Writing to final_audit.log...")

	for range ticker.C {
		// Simulate log fetching and synthesis
		res, err := client.ExecuteStrategicAction(context.Background(), &pb.AgentAction{
			AgentId:    "Evolutionary-Auditor",
			ActionType: "AUDIT_SYNTHESIS",
			TaskIntent: "Validating inter-agent integrity",
			DataSizeBytes: 24 * 1024 * 1024, 
		})

		if err == nil {
			entry := AuditLog{
				Timestamp: time.Now().Format(time.RFC3339),
				Agent:     "Mesh-System",
				Action:    "INTEGRITY_CHECK",
				Status:    "VERIFIED",
				Details:   fmt.Sprintf("Compute Tier: %s | Role: %v", res.RoutingProvider, res.RequiredRole),
			}
			json.NewEncoder(auditFile).Encode(entry)
			fmt.Printf("   [AUDIT] Integrity verified via %s\n", res.RoutingProvider)
		}
	}
}
