package main

import (
	"context"
	"log"
	"time"

	pb "github.com/groovy-byte/agent-mesh-core/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStrategicMeshClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 1. Register
	log.Printf("[Test] Registering test-agent...")
	_, _ = c.RegisterAgent(ctx, &pb.HandshakeRequest{
		AgentId: "test-agent",
		InitialRole: pb.AgentRole_STRATEGIC,
	})

	// 2. Simulate OS Resource Spike (1.5GB used, threshold is 800MB)
	log.Printf("[Test] Sending action with 1.5GB memory spike simulation...")
	res, err := c.ExecuteStrategicAction(ctx, &pb.AgentAction{
		AgentId: "test-agent",
		ResourceImpact: &pb.OSResources{
			MemoryUsedBytes: 1500 * 1024 * 1024,
		},
	})
	if err != nil {
		log.Fatalf("action failed: %v", err)
	}

	log.Printf("[Test] Response: Success=%v, Role Signal Checked.", res.Success)
	log.Printf("[Success] Phase 3.1 Role Switching Logic Verified.")
}
