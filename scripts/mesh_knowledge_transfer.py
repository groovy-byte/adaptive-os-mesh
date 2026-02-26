#!/usr/bin/env python3
import sys
import os
import logging
import glob

# Clean logging
logging.basicConfig(level=logging.INFO, format='%(message)s')
logger = logging.getLogger("Mesh-Knowledge-Transfer")

def main():
    logger.info("============================================================")
    logger.info("🚀 STARTING KNOWLEDGE TRANSFER: Files -> Local Qdrant")
    logger.info("============================================================")

    try:
        from qdrant_client import QdrantClient
        from qdrant_client.http import models
        from sentence_transformers import SentenceTransformer
        from pypdf import PdfReader
    except ImportError as e:
        logger.error(f"❌ DEPENDENCY ERROR: {e}")
        logger.error("Please wait for the background pip install to finish.")
        return

    # 1. Initialize Clients
    logger.info("📡 Connecting to Qdrant...")
    q_client = QdrantClient("localhost", port=6333)
    
    logger.info("🧠 Loading Local Embedding Model (All-MiniLM-L6-v2)...")
    model = SentenceTransformer('all-MiniLM-L6-v2')

    # 2. Prepare Collection
    collection_name = "research_corpus"
    q_client.recreate_collection(
        collection_name=collection_name,
        vectors_config=models.VectorParams(size=384, distance=models.Distance.COSINE),
    )

    # 3. Read and Index Papers
    papers_dir = "/home/groovy-byte/agent-mesh-core/local_research"
    pdf_files = glob.glob(os.path.join(papers_dir, "*.pdf"))
    
    logger.info(f"📂 Found {len(pdf_files)} research papers for transfer.")

    point_id = 1
    for pdf_path in pdf_files:
        filename = os.path.basename(pdf_path)
        logger.info(f"📖 Processing: {filename}")
        
        try:
            reader = PdfReader(pdf_path)
            text = ""
            for page in reader.pages:
                text += page.extract_text() + "\n"
            
            # Simple Chunking (1000 chars)
            chunks = [text[i:i+1000] for i in range(0, len(text), 1000)]
            
            # Generate Embeddings & Upsert
            embeddings = model.encode(chunks)
            
            points = [
                models.PointStruct(
                    id=point_id + i,
                    vector=embeddings[i].tolist(),
                    payload={"filename": filename, "content": chunks[i], "source": "research_corpus"}
                )
                for i in range(len(chunks))
            ]
            
            q_client.upsert(collection_name=collection_name, points=points)
            point_id += len(chunks)
            logger.info(f"   ✅ Indexed {len(chunks)} chunks.")
            
        except Exception as e:
            logger.error(f"   ❌ FAILED to process {filename}: {e}")

    logger.info("============================================================")
    logger.info(f"🎉 TRANSFER COMPLETE: {point_id - 1} total points in Qdrant.")
    logger.info("============================================================")

if __name__ == "__main__":
    main()
