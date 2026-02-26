#!/usr/bin/env python3
import sys
import os
import json
from qdrant_client import QdrantClient
from sentence_transformers import SentenceTransformer

def deep_research(query, collection="llama_research", limit=3):
    client = QdrantClient("localhost", port=6333)
    model = SentenceTransformer('all-MiniLM-L6-v2')
    vector = model.encode(query).tolist()
    
    results = client.query_points(
        collection_name=collection,
        query=vector,
        limit=limit
    ).points
    
    formatted_results = []
    for hit in results:
        formatted_results.append({
            "source": hit.payload.get("filename", "Unknown"),
            "content": hit.payload.get("content", ""),
            "score": hit.score
        })
    return formatted_results

def main():
    topics = [
        "Granular AI Compiler tricks for llama.cpp from nncase and AI Compiler.pdf",
        "3-bit quantization de-quantization speed Pareto frontier QuantX",
        "Model partitioning and CPU-GPU task scheduling ScheInfer",
        "Nacrith CDF-24 precision upgrade vocabulary floor overhead",
        "LLM-based hardware-aware quantization agent HAQA Llama speedup"
    ]
    
    all_knowledge = {}
    for topic in topics:
        all_knowledge[topic] = deep_research(topic)
        
    print(json.dumps(all_knowledge, indent=2))

if __name__ == "__main__":
    main()
