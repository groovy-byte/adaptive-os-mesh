#!/usr/bin/env python3
import sys
import os
import logging
from qdrant_client import QdrantClient
from sentence_transformers import SentenceTransformer

# Clean logging
logging.basicConfig(level=logging.INFO, format='%(message)s')
logger = logging.getLogger("Local-Mesh-Query")

def search(client, model, query_text, limit=2):
    # Vectorize query
    vector = model.encode(query_text).tolist()
    
    # Use the 'query' method which is available in this version
    results = client.query_points(
        collection_name="research_corpus",
        query=vector,
        limit=limit
    ).points
    
    return results

def main():
    client = QdrantClient("localhost", port=6333)
    model = SentenceTransformer('all-MiniLM-L6-v2')

    queries = [
        "Topology Evolution patterns AgentConductor paper",
        "Difficulty-Aware Agentic Orchestration query complexity",
        "Self-healing Differentiable Modal Logic paper",
        "Peak-Aware Orchestration Alignment in Time paper",
        "Exact Algorithms Resource Reallocation multi-agent"
    ]

    logger.info("============================================================")
    logger.info("🕵️  FINAL OPTIMIZATION RESEARCH: Local Qdrant Query")
    logger.info("============================================================")

    for q in queries:
        logger.info(f"🔍 Query: {q}")
        try:
            results = search(client, model, q)
            for i, res in enumerate(results):
                filename = res.payload.get("filename", "Unknown")
                content = res.payload.get("content", "")
                logger.info(f"   [{i+1}] Source: {filename} (Score: {res.score:.4f})")
                snippet = content.replace("\n", " ")[:250] + "..."
                logger.info(f"       {snippet}")
        except Exception as e:
            logger.error(f"   ❌ Search failed for '{q}': {e}")
        logger.info("-" * 20)

    logger.info("============================================================")

if __name__ == "__main__":
    main()
