import grpc
import sys
import os
import time

# Add proto path
sys.path.append(os.path.join(os.path.dirname(__file__), '..'))

import mesh_pb2
import mesh_pb2_grpc

def generate_llm_response(prompt, kv_cache_mb=0, controller_addr='localhost:50051'):
    print(f"\n🤖 Requesting LLM Inference...")
    print(f"   Prompt: '{prompt[:50]}...'")
    if kv_cache_mb > 0:
        print(f"   Expected KV Cache: {kv_cache_mb} MB")
    
    channel = grpc.insecure_channel(controller_addr)
    stub = mesh_pb2_grpc.StrategicMeshStub(channel)
    
    request = mesh_pb2.InferenceRequest(
        agent_id="hardware-aware-agent",
        prompt=prompt,
        max_tokens=100,
        temperature=0.7,
        expected_kv_cache_bytes=kv_cache_mb * 1024 * 1024
    )
    
    start_time = time.time()
    try:
        response = stub.GenerateResponse(request)
        duration = (time.time() - start_time) * 1000
        print(f"✅ RESPONSE RECEIVED ({duration:.2f}ms)")
        print(f"🎯 HARDWARE PATH: {response.hardware_path}")
        print(f"📝 TEXT: {response.text}")
        return response
    except grpc.RpcError as e:
        print(f"❌ gRPC Error: {e.code()} - {e.details()}")
        return None

def main():
    print("============================================================")
    print("🚀 HARDWARE-AWARE LLM INTEGRATION TEST")
    print("============================================================")
    
    # Test 1: Small Reasoning Task (Cache-Resident)
    small_prompt = "Explain the JADE pattern in distributed systems."
    generate_llm_response(small_prompt, kv_cache_mb=2)
    
    # Test 2: Large Context Task (DRAM-bound / GPU preferred)
    large_prompt = "Analyze the following codebase for security vulnerabilities..."
    generate_llm_response(large_prompt, kv_cache_mb=128)
    
    print("\n============================================================")
    print("🏁 PHASE 8 TEST CONCLUDED")
    print("============================================================")

if __name__ == "__main__":
    main()
