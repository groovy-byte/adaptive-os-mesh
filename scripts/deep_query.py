#!/usr/bin/env python3
import sys
import os
import logging
from qdrant_client import QdrantClient
from sentence_transformers import SentenceTransformer
import textwrap

# ANSI Colors for terminal visibility
class Colors:
    HEADER = '\033[95m'
    BLUE = '\033[94m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    BOLD = '\033[1m'
    RESET = '\033[0m'

logging.basicConfig(level=logging.INFO, format='%(message)s')

def deep_search(client, model, query, collection="llama_research", limit=5):
    vector = model.encode(query).tolist()
    results = client.query_points(
        collection_name=collection,
        query=vector,
        limit=limit
    ).points
    return results

def main():
    client = QdrantClient("localhost", port=6333)
    model = SentenceTransformer('all-MiniLM-L6-v2')

    queries = [
        "Hardware-aware quantization for ARM/x86 architectures from QuantX",
        "Adaptive Synthesis Protocol mathematical constraints AdaptOrch",
        "eBPF implementation AgentCgroup resource throttling",
        "Task scheduling algorithms LLM inference concurrency",
        "Failure modes mesh networks bandwidth saturation"
    ]

    print(Colors.HEADER + Colors.BOLD + "="*80 + Colors.RESET)
    print(Colors.HEADER + Colors.BOLD + "{:^80}".format("🕵️  MAESTRO DEEP KNOWLEDGE EXTRACTION") + Colors.RESET)
    print(Colors.HEADER + Colors.BOLD + "="*80 + Colors.RESET)

    for q in queries:
        print("\n" + Colors.YELLOW + Colors.BOLD + "🔎 DEEP QUERY:" + Colors.RESET + " " + Colors.BOLD + q + Colors.RESET)
        print(Colors.YELLOW + "-"*80 + Colors.RESET)
        
        all_results = []
        for coll in ["llama_research", "research_corpus"]:
            try:
                all_results.extend(deep_search(client, model, q, collection=coll))
            except:
                continue
        
        all_results.sort(key=lambda x: x.score, reverse=True)
        
        for i, res in enumerate(all_results[:3]):
            filename = res.payload.get("filename", "Unknown Source")
            content = res.payload.get("content", "No content available.")
            score = res.score
            
            print("\n  " + Colors.BLUE + Colors.BOLD + "[RESULT " + str(i+1) + "] Source: " + filename + " | Relevance: " + "{:.4f}".format(score) + Colors.RESET)
            
            clean_content = content.replace("\n", " ").strip()
            wrapped = textwrap.fill(clean_content, width=76, initial_indent="    ", subsequent_indent="    ")
            print(wrapped)

    print("\n" + Colors.HEADER + Colors.BOLD + "="*80 + Colors.RESET)
    print(Colors.HEADER + Colors.BOLD + "{:^80}".format("✅ DEEP EXTRACTION COMPLETE") + Colors.RESET)
    print(Colors.HEADER + Colors.BOLD + "="*80 + Colors.RESET + "\n")

if __name__ == "__main__":
    main()
