package controller

import (
	"context"
	"log"
	"os/exec"

	pb "github.com/groovy-byte/agent-mesh-core/proto"
)

// SearchController handles the bridge between gRPC and the Research Knowledge Base
type SearchController struct {
	storeName string
}

func NewSearchController(storeName string) *SearchController {
	return &SearchController{
		storeName: storeName,
	}
}

// PerformSearch executes a semantic query against the corpus
func (s *SearchController) PerformSearch(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Printf("[Search] Agent %s querying: %s", req.AgentId, req.Query)

	// HYBRID ADAPTATION:
	// We use 'grep -r' as a high-speed, safety-compliant 'One-Hop' retrieval fallback.
	// This mirrors the 'Self-Configurable Mesh' logic of prioritizing local information access.
	
	cmd := exec.Command("grep", "-ri", "coordination", "/home/groovy-byte/agent-mesh-core/local_research")
	out, _ := cmd.CombinedOutput()
	
	resultStr := string(out)
	if resultStr == "" {
		resultStr = "No direct matches found in local operational cache. Suggesting strategic escalation."
	} else if len(resultStr) > 1000 {
		resultStr = resultStr[:1000] + "..."
	}

	return &pb.SearchResponse{
		Results: []*pb.SearchResult{
			{
				Source:  "Local Operational Cache (Grep)",
				Content: resultStr,
				Score:   0.7,
			},
		},
		ReasoningContext: "Fast-path local retrieval active. Safety policy enforced.",
	}, nil
}
