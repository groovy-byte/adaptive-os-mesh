#!/usr/bin/env python3
import sys
import os
import logging
import json
from flask import Flask, request, jsonify
from qdrant_client import QdrantClient
from sentence_transformers import SentenceTransformer

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(message)s')
logger = logging.getLogger("Mesh-Retriever-Service")

app = Flask(__name__)

# Global instances for Qdrant client and Sentence Transformer model
# These are loaded once when the service starts
qdrant_client = None
sentence_model = None

@app.before_request
def initialize_globals():
    global qdrant_client, sentence_model
    if qdrant_client is None:
        logger.info("Initializing QdrantClient...")
        qdrant_client = QdrantClient("localhost", port=6333)
    if sentence_model is None:
        logger.info("Initializing SentenceTransformer model (this may take a moment)...")
        sentence_model = SentenceTransformer('all-MiniLM-L6-v2')
        logger.info("SentenceTransformer model loaded.")

@app.route('/search', methods=['POST'])
def search_endpoint():
    try:
        data = request.get_json()
        query_text = data.get('query')
        limit = data.get('limit', 1)
        collections = data.get('collections', ["research_corpus", "llama_research"])

        if not query_text:
            return jsonify({"error": "Query text is required"}), 400

        vector = sentence_model.encode(query_text).tolist()
        
        results = []
        for coll in collections:
            try:
                hits = qdrant_client.query_points(
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
            except Exception as e:
                logger.warning(f"Error searching collection {coll}: {e}")
                # Don't fail the whole request if one collection has an issue

        results.sort(key=lambda x: x["score"], reverse=True)
        return jsonify(results[:limit])

    except Exception as e:
        logger.error(f"Error in search_endpoint: {e}")
        return jsonify({"error": str(e)}), 500

if __name__ == '__main__':
    # Flask defaults to port 5000
    # For production, use a more robust WSGI server like Gunicorn
    app.run(host='127.0.0.1', port=5000, debug=False)
