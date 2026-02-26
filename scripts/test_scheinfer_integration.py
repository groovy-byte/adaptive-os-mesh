import grpc
import sys
import os
import json

# Add proto path
sys.path.append(os.path.join(os.path.dirname(__file__), '..'))

import mesh_pb2
import mesh_pb2_grpc

def test_routing(agent_id, data_size_mb, controller_addr='localhost:50051'):
    print(f"📡 Testing ScheInfer Integration for {agent_id} with {data_size_mb}MB Task...")
    channel = grpc.insecure_channel(controller_addr)
    stub = mesh_pb2_grpc.StrategicMeshStub(channel)
    
    # 1. Handshake
    handshake = mesh_pb2.HandshakeRequest(
        agent_id=agent_id,
        capabilities=['scheinfer-test'],
        initial_role=mesh_pb2.STRATEGIC
    )
    stub.RegisterAgent(handshake)
    
    # 2. Execute Action with Data Size hint
    action = mesh_pb2.AgentAction(
        agent_id=agent_id,
        action_type="QUANTX_DEQUANT",
        data_size_bytes=data_size_mb * 1024 * 1024,
        task_intent="PERFORMANCE_BENCHMARK"
    )
    
    try:
        response = stub.ExecuteStrategicAction(action)
        print(f"✅ ACTION SUCCESS: {response.success}")
        print(f"🎯 SCHEINFER ROUTING: {response.routing_provider}")
        return response.routing_provider
    except grpc.RpcError as e:
        print(f"❌ gRPC Error: {e.code()} - {e.details()}")
        return None

def main():
    print("============================================================")
    print("🚀 SCHEINFER INTEGRATION TEST")
    print("============================================================")
    
    # Test 1: Small task (Cache-resident)
    provider_small = test_routing("agent-small", 4)
    
    # Test 2: Large task (DRAM-bound)
    provider_large = test_routing("agent-large", 64)
    
    print("\n============================================================")
    print("🏁 SUMMARY")
    print(f"Small Task (4MB)  -> {provider_small}")
    print(f"Large Task (64MB) -> {provider_large}")
    print("============================================================")

if __name__ == "__main__":
    main()
