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
	// 1. Connect to Vextra Controller
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStrategicMeshClient(conn)

	// 2. Test Task 1.3: One-Hop Handshake
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	log.Printf("[Test] Registering Student-A (Strategic)...")
	r1, err := c.RegisterAgent(ctx, &pb.HandshakeRequest{
		AgentId:      "student-a",
		Capabilities: []string{"reasoning", "planning"},
		InitialRole:  pb.AgentRole_STRATEGIC,
	})
	if err != nil {
		log.Fatalf("could not register student-a: %v", err)
	}
	log.Printf("[Test] Student-A Response: Approved=%v, Session=%s", r1.Approved, r1.SessionId)

	log.Printf("[Test] Registering Student-B (Operational)...")
	r2, err := c.RegisterAgent(ctx, &pb.HandshakeRequest{
		AgentId:      "student-b",
		Capabilities: []string{"tool-execution", "monitoring"},
		InitialRole:  pb.AgentRole_OPERATIONAL,
	})
	if err != nil {
		log.Fatalf("could not register student-b: %v", err)
	}
	log.Printf("[Test] Student-B Response: Approved=%v, Session=%s", r2.Approved, r2.SessionId)

	log.Printf("[Success] Phase 1 Foundation Verified.")
}
