#!/usr/bin/env python3
import time
import json
import logging
import grpc
import sys
import os

# Add proto path
sys.path.append(os.path.join(os.path.dirname(__file__), '..'))

import mesh_pb2
import mesh_pb2_grpc

logging.basicConfig(level=logging.INFO, format='%(message)s')
logger = logging.getLogger("Fluid-Agent")

class FluidAgent:
    def __init__(self, agent_id):
        self.agent_id = agent_id
        self.channel = grpc.insecure_channel('localhost:50051')
        self.stub = mesh_pb2_grpc.StrategicMeshStub(self.channel)
        self.role = mesh_pb2.OPERATIONAL
        logger.info(f"🚀 INITIALIZED: Fluid Agent {self.agent_id}")

    def register(self):
        logger.info(f"📡 Registering with Mesh...")
        req = mesh_pb2.HandshakeRequest(agent_id=self.agent_id, initial_role=self.role)
        res = self.stub.RegisterAgent(req)
        logger.info(f"✅ REGISTERED: Session {res.session_id}")

    def perform_task(self, complexity="LOW", memory_mb=100):
        logger.info(f"🛠️  Task Complexity: {complexity} | Memory: {memory_mb}MB")
        
        # Simulate OS resource usage
        resource_impact = mesh_pb2.OSResources(memory_used_bytes=memory_mb * 1024 * 1024)
        
        action = mesh_pb2.AgentAction(
            agent_id=self.agent_id,
            action_type="OS_TASK",
            resource_impact=resource_impact,
            reasoning_chain=f"Executing {complexity} complexity task"
        )
        
        response = self.stub.ExecuteStrategicAction(action)
        
        # ADAPTIVE ROLE SWITCHING (Single Source of Truth)
        if response.required_role != self.role:
            role_name = "STRATEGIC" if response.required_role == mesh_pb2.STRATEGIC else "OPERATIONAL"
            if response.required_role == mesh_pb2.STRATEGIC:
                logger.info(f"🔼 ROLE EVOLUTION: Promoting to {role_name} (Boss Mode)")
            else:
                logger.info(f"🔽 COUNTERBALANCE: Controller mandated demotion to {role_name}")
            self.role = response.required_role

    def shutdown(self):
        self.channel.close()

def main():
    logger.info("============================================================")
    logger.info(f"🏃 FLUID AGENT DEMO: {sys.argv[1] if len(sys.argv) > 1 else 'Student-X'}")
    logger.info("============================================================")
    
    agent = FluidAgent(sys.argv[1] if len(sys.argv) > 1 else 'student-x')
    agent.register()
    
    # Task 1: Normal task (Stays OPERATIONAL)
    agent.perform_task("LOW", 200)
    
    # Task 2: High Difficulty task (Tests STRATEGIC promotion)
    logger.info("📡 Requesting Strategic Promotion...")
    action = mesh_pb2.AgentAction(
        agent_id=agent.agent_id,
        action_type="STRATEGIC_REASONING",
        task_intent="STRATEGIC_REASONING_PLANNING",
        reasoning_chain="Analyzing complex mesh topology"
    )
    res = agent.stub.ExecuteStrategicAction(action)
    logger.info(f"DEBUG: Response Required Role: {res.required_role} | Current Role: {agent.role}")
    if res.required_role != agent.role:
        role_name = "STRATEGIC" if res.required_role == mesh_pb2.STRATEGIC else "OPERATIONAL"
        if res.required_role == mesh_pb2.STRATEGIC:
            logger.info(f"🔼 ROLE EVOLUTION: Promoted to {role_name} (Boss Mode)")
        else:
            logger.info(f"🔽 COUNTERBALANCE: Controller mandated demotion to {role_name}")
        agent.role = res.required_role

    # Task 3: Resource spike test (Tests Controller-mandated demotion)
    agent.perform_task("HIGH", 1200)
    
    logger.info("✅ FLUID TASK COMPLETED")
    logger.info("============================================================")

if __name__ == "__main__":
    main()
