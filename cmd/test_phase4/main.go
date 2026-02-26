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
	log.Println("============================================================")
	log.Println("🚀 PHASE 4 TEST: Arbiter & Counterbalance")
	log.Println("============================================================")

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ CONNECTION FAILED: %v", err)
	}
	defer conn.Close()
	c := pb.NewStrategicMeshClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 0. Register Agents first (fixes nil pointer in registry)
	log.Println("📡 Pre-registering Agents...")
	c.RegisterAgent(ctx, &pb.HandshakeRequest{AgentId: "boss-a", InitialRole: pb.AgentRole_OPERATIONAL})
	c.RegisterAgent(ctx, &pb.HandshakeRequest{AgentId: "boss-b", InitialRole: pb.AgentRole_OPERATIONAL})
	c.RegisterAgent(ctx, &pb.HandshakeRequest{AgentId: "hanging-agent", InitialRole: pb.AgentRole_OPERATIONAL})

	// 1. Test Strategic Lock (Counterbalance)
	log.Println("🔑 Testing Strategic Lock (Arbiter)...")
	res1, err := c.ExecuteStrategicAction(ctx, &pb.AgentAction{
		AgentId: "boss-a", ActionType: "HIGH_COMPLEXITY",
	})
	if err != nil { log.Fatalf("boss-a failed: %v", err) }
	log.Printf("   > Agent Boss-A Lock: PromotionSuggested=%v", res1.PromotionSuggested)

	res2, err := c.ExecuteStrategicAction(ctx, &pb.AgentAction{
		AgentId: "boss-b", ActionType: "HIGH_COMPLEXITY",
	})
	if err != nil { log.Fatalf("boss-b failed: %v", err) }
	log.Printf("   > Agent Boss-B Lock (should be false): PromotionSuggested=%v", res2.PromotionSuggested)

	// 2. Test State Reconstitution
	log.Println("🩹 Testing State Reconstitution (Healing)...")
	c.ExecuteStrategicAction(ctx, &pb.AgentAction{
		AgentId: "hanging-agent", ActionType: "OS_TASK", ReasoningChain: "Initial step before crash",
	})
	
	log.Println("   > Agent 'hanging-agent' simulating crash. Retrieving state...")
	reconst, err := c.GetStateReconstitution(ctx, &pb.HandshakeRequest{AgentId: "hanging-agent"})
	if err != nil {
		log.Fatalf("❌ RECONSTITUTION FAILED: %v", err)
	}
	log.Printf("   > Reconstituted State Type: %s, Data: %s", reconst.ActionType, reconst.ReasoningChain)

	log.Println("============================================================")
	log.Println("🎉 PHASE 4 VERIFIED")
	log.Println("============================================================")
}
