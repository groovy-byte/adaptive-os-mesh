import grpc
import time
import json
import random
import concurrent.futures
from datetime import datetime
import os
import sys

# Add proto path
sys.path.append(os.path.join(os.path.dirname(__file__), '..'))

import mesh_pb2
import mesh_pb2_grpc

def agent_workflow(agent_id, tasks):
    channel = grpc.insecure_channel('localhost:50051')
    stub = mesh_pb2_grpc.StrategicMeshStub(channel)
    
    print(f"🚀 [{agent_id}] Registering...")
    try:
        stub.RegisterAgent(mesh_pb2.HandshakeRequest(
            agent_id=agent_id,
            capabilities=["STRESS_TEST"],
            initial_role=mesh_pb2.AgentRole.OPERATIONAL
        ))
    except Exception as e:
        print(f"   [{agent_id}] ✖ Registration FAILED: {e}")
        return []

    print(f"🚀 [{agent_id}] Starting workload...")
    
    logs = []
    for i, (task_type, size_mb) in enumerate(tasks):
        # 1. Action
        try:
            action = stub.ExecuteStrategicAction(mesh_pb2.AgentAction(
                agent_id=agent_id,
                action_type=task_type,
                data_size_bytes=size_mb * 1024 * 1024,
                reasoning_chain=f"Task {i} for {agent_id}",
                task_intent=f"Project stress test component {i}"
            ))
            
            # 2. Inference
            inference = stub.GenerateResponse(mesh_pb2.InferenceRequest(
                agent_id=agent_id,
                prompt=f"Process results for {task_type}",
                expected_kv_cache_bytes=size_mb * 1024 * 1024
            ))
            
            log_entry = {
                "agent": agent_id,
                "step": i,
                "type": task_type,
                "hardware": action.routing_provider,
                "latency": inference.latency_ms,
                "status": "SUCCESS"
            }
            print(f"   [{agent_id}] Step {i} ({task_type}) verified on {action.routing_provider}")
            
        except Exception as e:
            log_entry = {"agent": agent_id, "step": i, "status": "FAILED", "error": str(e)}
            print(f"   [{agent_id}] ✖ Step {i} FAILED: {e}")
            
        logs.append(log_entry)
        time.sleep(random.uniform(0.1, 0.5))
    
    return logs

def main():
    print("🔥 Starting Distributed Mesh Stress Workload...")
    
    agents = {
        "Scout-Alpha": [("SEARCH", 2), ("REASONING", 32)],
        "Coder-Beta": [("CODING", 16), ("OS_COMPILATION", 64)],
        "Reviewer-Gamma": [("REASONING", 128), ("SYNTHESIS", 8)],
        "Perf-Delta": [("SEARCH", 1), ("REASONING", 256)],
        "Strategic-Omega": [("SYNTHESIS", 512), ("REASONING", 1024)]
    }

    with concurrent.futures.ThreadPoolExecutor(max_workers=5) as executor:
        futures = {executor.submit(agent_workflow, aid, tasks): aid for aid, tasks in agents.items()}
        
        all_logs = []
        for future in concurrent.futures.as_completed(futures):
            all_logs.extend(future.result())

    # Export structured logs for the Auditor
    with open("workload_execution.jsonl", "w") as f:
        for entry in all_logs:
            f.write(json.dumps(entry) + "\n")
            
    print(f"✅ Workload complete. {len(all_logs)} log entries generated in workload_execution.jsonl")

if __name__ == "__main__":
    main()
