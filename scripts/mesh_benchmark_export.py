import grpc
import sys
import os
import time
import json
from datetime import datetime, timezone

# Add proto path
sys.path.append(os.path.join(os.path.dirname(__file__), '..'))

import mesh_pb2
import mesh_pb2_grpc

def run_benchmark(agent_id, kv_cache_mb, task_intent="REASONING", controller_addr=None):
    if controller_addr is None:
        controller_addr = os.getenv('MESH_CONTROLLER_ADDR', 'localhost:50051')
    channel = grpc.insecure_channel(controller_addr)
    stub = mesh_pb2_grpc.StrategicMeshStub(channel)
    
    prompt = f"Benchmark test for {kv_cache_mb}MB KV cache by {agent_id}."
    request = mesh_pb2.InferenceRequest(
        agent_id=agent_id,
        prompt=prompt,
        max_tokens=50,
        expected_kv_cache_bytes=kv_cache_mb * 1024 * 1024
    )
    
    start_time = time.perf_counter()
    try:
        response = stub.GenerateResponse(request)
        end_time = time.perf_counter()
        
        return {
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "agent_id": agent_id,
            "kv_cache_mb": kv_cache_mb,
            "task_intent": task_intent,
            "hardware_path": response.hardware_path,
            "latency_ms": (end_time - start_time) * 1000,
            "tokens_used": response.tokens_used,
            "throughput_gb_s": response.throughput_gbs,
            "avx512_active": response.avx512_usage,
            "status": "success"
        }
    except Exception as e:
        return {
            "timestamp": datetime.now(timezone.utc).isoformat(),
            "agent_id": agent_id,
            "status": "error",
            "error": str(e)
        }

def run_synthesis_benchmark(agent_ids, goal="Final Decision Synthesis", controller_addr=None):
    if controller_addr is None:
        controller_addr = os.getenv('MESH_CONTROLLER_ADDR', 'localhost:50051')
    channel = grpc.insecure_channel(controller_addr)
    stub = mesh_pb2_grpc.StrategicMeshStub(channel)
    
    actions = [mesh_pb2.AgentAction(agent_id=aid, reasoning_chain="Analyzed local state") for aid in agent_ids]
    request = mesh_pb2.SynthesisRequest(
        agent_ids=agent_ids,
        target_goal=goal,
        actions_to_merge=actions
    )
    
    try:
        stub.SynthesizeOutputs(request)
        return True
    except Exception as e:
        print(f"Synthesis failed: {e}")
        return False

def main():
    print("🚀 Running Adaptive OS Mesh Multi-Agent Benchmarks...")
    
    agents = ["Scout-Agent", "Coder-Agent", "Reviewer-Agent", "Performance-Engineer"]
    # Sizes targeting different hardware paths
    tasks = [
        (2, "SEARCH"), 
        (16, "CODING"), 
        (64, "REASONING"), 
        (128, "SYNTHESIS")
    ]
    
    results = []
    for agent in agents:
        print(f"📊 Benchmarking {agent}...")
        for size, intent in tasks:
            print(f"   - Task: {intent} ({size}MB)...")
            res = run_benchmark(agent, size, intent)
            results.append(res)
            time.sleep(0.1)

    print("🧩 Triggering Value of Contribution (VoC) Synthesis...")
    # Simulate inter-agent influence
    run_synthesis_benchmark(["Scout-Agent", "Coder-Agent"])
    run_synthesis_benchmark(["Coder-Agent", "Reviewer-Agent", "Performance-Engineer"])
    
    # Mocked contribution matrix for TUI demonstration
    contribution_matrix = {
        "Scout-Agent": {"Coder-Agent": 0.85, "Reviewer-Agent": 0.12},
        "Coder-Agent": {"Reviewer-Agent": 0.94, "Performance-Engineer": 0.45},
        "Reviewer-Agent": {"Coder-Agent": 0.78, "Scout-Agent": 0.23},
        "Performance-Engineer": {"Reviewer-Agent": 0.62, "Scout-Agent": 0.15}
    }

    # Mocked tool and failure logs
    agent_logs = {
        "Scout-Agent": {"tool_calls": 42, "failed_tasks": ["SearchTimeout", "PDFParseError"]},
        "Coder-Agent": {"tool_calls": 128, "failed_tasks": ["CompilationError", "LinterWarning"]},
        "Reviewer-Agent": {"tool_calls": 15, "failed_tasks": ["SecurityAuditInterrupted"]},
        "Performance-Engineer": {"tool_calls": 8, "failed_tasks": ["CudaOOM"]}
    }
    
    # Export to JSON for Voodoo UI and Waybar
    output_path = os.path.join(os.path.dirname(__file__), '../benchmark_results.json')
    with open(output_path, 'w') as f:
        json.dump({
            "mesh_id": "adaptive-os-mesh-v1",
            "last_updated": datetime.now(timezone.utc).isoformat(),
            "agents_active": len(agents),
            "benchmarks": results,
            "contribution_matrix": contribution_matrix,
            "agent_logs": agent_logs
        }, f, indent=2)
    
    print(f"✅ Multi-agent benchmarks complete. Results exported to {output_path}")

if __name__ == "__main__":
    main()
