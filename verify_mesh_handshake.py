import grpc
import mesh_pb2
import mesh_pb2_grpc

def run():
    # Registration with controller at localhost:50051
    print("Connecting to StrategicMesh controller at localhost:50051...")
    channel = grpc.insecure_channel('localhost:50051')
    stub = mesh_pb2_grpc.StrategicMeshStub(channel)
    
    request = mesh_pb2.HandshakeRequest(
        agent_id='test-handshake-agent',
        capabilities=['handshake-verification', 'os-resource-tracking'],
        initial_role=mesh_pb2.STRATEGIC
    )
    
    print(f"Sending HandshakeRequest for agent: {request.agent_id}")
    try:
        response = stub.RegisterAgent(request)
        print("\nHandshakeResponse received:")
        print(f"Session ID: {response.session_id}")
        print(f"Approved: {response.approved}")
        if response.error_message:
            print(f"Error Message: {response.error_message}")
        
        if response.resource_limits:
            print(f"Resource Limits:")
            print(f"  CPU Usage Percent Limit: {response.resource_limits.cpu_usage_percent}")
            print(f"  Memory Used Bytes Limit: {response.resource_limits.memory_used_bytes}")
            print(f"  Memory Total Bytes Limit: {response.resource_limits.memory_total_bytes}")
            print(f"  Disk IO Wait Limit: {response.resource_limits.disk_io_wait}")
    except grpc.RpcError as e:
        print(f"\ngRPC Error: {e.code()} - {e.details()}")

if __name__ == '__main__':
    run()
