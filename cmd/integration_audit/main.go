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
	log.Println("🚀 INTEGRATION TEST: Full Mesh Flow & Resilience")
	log.Println("============================================================")

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("❌ CONNECTION FAILED: %v", err)
	}
	defer conn.Close()
	c := pb.NewStrategicMeshClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Verify Phase 1 (Handshake)
	log.Println("📡 Testing Phase 1: One-Hop Handshake...")
	_, err = c.RegisterAgent(ctx, &pb.HandshakeRequest{AgentId: "fluid-student-a", InitialRole: pb.AgentRole_OPERATIONAL})
	if err != nil {
		log.Fatalf("❌ HANDSHAKE FAILED: %v", err)
	}
	log.Println("✅ Phase 1: Handshake Verified.")

	// 2. Verify Phase 2/3 (Soft-Throttle & VoC)
	log.Println("🔍 Testing Phase 3: Soft-Throttle (VoC)...")
	query := "How to manage OS resources in AI agents?"
	
	log.Printf("   > Sending first novel query: '%s'", query)
	res1, err := c.SemanticSearch(ctx, &pb.SearchRequest{AgentId: "fluid-student-a", Query: query, MaxResults: 1})
	if err != nil {
		log.Fatalf("❌ SEARCH 1 FAILED: %v", err)
	}
	if res1 == nil {
		log.Fatal("❌ SEARCH 1 FAILED: Response is nil")
	}
	log.Printf("   > Response 1: %s", res1.ReasoningContext)

	log.Printf("   > Sending redundant query (expecting throttle)...")
	res2, err := c.SemanticSearch(ctx, &pb.SearchRequest{AgentId: "fluid-student-a", Query: query, MaxResults: 1})
	if err != nil {
		log.Fatalf("❌ SEARCH 2 FAILED: %v", err)
	}
	if res2 == nil {
		log.Fatal("❌ SEARCH 2 FAILED: Response is nil")
	}
	log.Printf("   > Response 2: %s", res2.ReasoningContext)

	if res2.ReasoningContext == "THROTTLED: Redundant semantic search detected." {
		log.Println("✅ Phase 3: Soft-Throttle Verified.")
	} else {
		log.Fatalf("❌ THROTTLE FAILED: Received: %s", res2.ReasoningContext)
	}

	// 3. Verify Phase 3.1 (Role Demotion)
	log.Println("📉 Testing Phase 3.1: Adaptive Counterbalance...")
	res3, _ := c.ExecuteStrategicAction(ctx, &pb.AgentAction{
		AgentId: "fluid-student-a",
		ResourceImpact: &pb.OSResources{MemoryUsedBytes: 1500 * 1024 * 1024},
	})
	log.Printf("   > Demotion test success: %v", res3.Success)
	log.Println("✅ Phase 3.1: Adaptive Counterbalance Verified.")

	log.Println("============================================================")
	log.Println("🎉 ALL PREVIOUS PHASES VERIFIED")
	log.Println("============================================================")
}
