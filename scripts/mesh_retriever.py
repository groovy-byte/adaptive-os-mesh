#!/usr/bin/env python3
import sys
import json
from qdrant_client import QdrantClient
from sentence_transformers import SentenceTransformer

def main():
    if len(sys.argv) < 2:
        return
    
    query = sys.argv[1]
    limit = int(sys.argv[2]) if len(sys.argv) > 2 else 1
    
    client = QdrantClient("localhost", port=6333)
    model = SentenceTransformer('all-MiniLM-L6-v2')
    
    vector = model.encode(query).tolist()
    
    results = []
    # Search both collections for maximum grounding
    for coll in ["llama_research", "research_corpus"]:
        try:
            hits = client.query_points(
                collection_name=coll,
                query=vector,
                limit=limit
            ).points
            for hit in hits:
                results.append({
                    "source": hit.payload.get("filename", coll),
                    "content": hit.payload.get("content", ""),
                    "score": hit.score
                })
        except:
            continue
            
    # Sort and return top matches
    results.sort(key=lambda x: x["score"], reverse=True)
    print(json.dumps(results[:limit]))

if __name__ == "__main__":
    main()
