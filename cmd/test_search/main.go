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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("[Test] Executing Semantic Search Request...")
	res, err := c.SemanticSearch(ctx, &pb.SearchRequest{
		AgentId:    "test-boss-agent",
		Query:      "How can we prevent agent over-allocation in mesh networks?",
		MaxResults: 3,
	})
	if err != nil {
		log.Fatalf("search failed: %v", err)
	}

	log.Printf("[Test] Search Response Received!")
	log.Printf("Reasoning Context: %s", res.ReasoningContext)
	for _, r := range res.Results {
		log.Printf("Source: %s | Score: %f", r.Source, r.Score)
		// Truncate output for log readability
		content := r.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}
		log.Printf("Content: %s", content)
	}

	log.Printf("[Success] Phase 2 Semantic Search Verified.")
}
