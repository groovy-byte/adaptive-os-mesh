#!/usr/bin/env python3
import argparse
import sys
import os
import logging
import grpc
import time
import random
import textwrap

# Add paths for proto
sys.path.append(os.path.join(os.path.dirname(__file__), '..'))
import mesh_pb2
import mesh_pb2_grpc

# --- UI Configuration (Professional CLI with Standard ANSI Colors) ---
class UI:
    # Standard 16 ANSI colors for best compatibility
    CYAN = '\033[96m'
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    MAGENTA = '\033[95m'
    BLUE = '\033[94m'
    RED = '\033[91m'
    BOLD = '\033[1m'
    RESET = '\033[0m'

    @staticmethod
    def header(text):
        print(f"\n{UI.BOLD}{UI.MAGENTA}{'='*80}{UI.RESET}")
        print(f"{UI.BOLD}{UI.MAGENTA}{text:^80}{UI.RESET}")
        print(f"{UI.BOLD}{UI.MAGENTA}{'='*80}{UI.RESET}\n")

    @staticmethod
    def speak(agent, role, text, color):
        prefix = f"{UI.BOLD}{color}[{agent} - {role}]{UI.RESET}"
        print(f"{prefix}")
        wrapped = textwrap.fill(text, width=76, initial_indent="  ", subsequent_indent="  ")
        print(f"{color}{wrapped}{UI.RESET}\n")

    @staticmethod
    def divider(text, color='\033[90m'): # Dark Gray
        print(f"{color}{UI.BOLD}{'--- ' + text + ' ' + '-'*(75-len(text))}{UI.RESET}")

# --- Agent Mesh Components ---

class MeshAgent:
    def __init__(self, agent_id, role, color, grpc_addr):
        self.agent_id = agent_id
        self.role = role
        self.color = color
        self.grpc_addr = grpc_addr
        self.research_queue = []

    def search(self, query, limit=3):
        try:
            with grpc.insecure_channel(self.grpc_addr) as channel:
                stub = mesh_pb2_grpc.StrategicMeshStub(channel)
                res = stub.SemanticSearch(mesh_pb2.SearchRequest(
                    agent_id=self.agent_id,
                    query=query,
                    max_results=limit
                ))
                if not res.results:
                    self.research_queue.append(query)
                    return None
                return res.results
        except Exception as e:
            print(f"Mesh Search Error: {e}")
            return None

    def cite(self, result):
        source = result.source
        content = result.content
        return f"\"{content.strip()[:250]}...\" (Source: {source})"

class LlamaDebate:
    def __init__(self, grpc_addr):
        self.grpc_addr = grpc_addr
        self.agent_a = MeshAgent("Student-A", "Hardware Specialist", UI.BLUE, self.grpc_addr)
        self.agent_b = MeshAgent("Student-B", "System Architect", UI.GREEN, self.grpc_addr)
        
    def run(self):
        UI.header("LLAMA.CPP DEBATE: QUANTIZATION VS. SCHEDULING")

        UI.divider("PHASE 1: OPENING ARGUMENTS")
        
        # Student A: Focusing on HAQA and QuantX (New Papers)
        res_a = self.agent_a.search("QuantX hardware-aware de-quantization speedup llama.cpp")
        arg_a = (
            "I advocate for the HAQA/QuantX framework. Our new research confirms that 3-bit resolution "
            "with hardware-aware constraints achieves de-quantization de-bottlenecking. "
            f"As stated in the source: {self.agent_a.cite(res_a[0])}. "
            "Hand-tuning is dead; we must jointly optimize kernels for 2.3x speedups on platforms like Snapdragon."
        )
        UI.speak(self.agent_a.agent_id, self.agent_a.role, arg_a, self.agent_a.color)
        time.sleep(1.5)

        # Student B: Focusing on ScheInfer and AI Compilers (New Papers)
        res_b = self.agent_b.search("ScheInfer model partitioning CPU GPU asynchronous processing")
        arg_b = (
            "My opponent overlooks the interconnect bottleneck. The ScheInfer paper introduces model partitioning "
            "for asynchronous processing. "
            f"Direct citation: {self.agent_b.cite(res_b[0])}. "
            "By saturating CPU, GPU, and the PCIe bridge simultaneously, we outperform hand-optimized llama.cpp "
            "by 1.25x to 2.04x regardless of quantization level."
        )
        UI.speak(self.agent_b.agent_id, self.agent_b.role, arg_b, self.agent_b.color)
        time.sleep(1.5)

        UI.divider("PHASE 2: CROSS-EXAMINATION")
        
        # Q1: A asks B about memory wall
        q1 = "Student B, how does your asynchronous partitioning solve the 'Memory Wall' described in the AI Compiler paper if the weights are still too large for the cache?"
        print(f"{UI.RED}{UI.BOLD}[STUDENT-A CHALLENGES]:{UI.RESET} {q1}\n")
        time.sleep(1)
        
        res_b_ans = self.agent_b.search("nncase e-graph based term rewriting memory locality optimization")
        ans_b = (
            "We solve this via e-graph based term rewriting and the NTT Tensor Template Library. "
            f"The literature confirms: {self.agent_b.cite(res_b_ans[0])}. "
            "We achieve register-level efficiency by static task partitioning, essentially 'hiding' the memory wall."
        )
        UI.speak(self.agent_b.agent_id, self.agent_b.role, ans_b, self.agent_b.color)
        time.sleep(1.5)

        # Q2: B asks A about vocabulary efficiency
        q2 = "Student A, HAQA focuses on weights, but what about the vocabulary floor overhead? Does your quantization handle the CDF-24 upgrade mentioned in Nacrith?"
        print(f"{UI.RED}{UI.BOLD}[STUDENT-B CHALLENGES]:{UI.RESET} {q2}\n")
        time.sleep(1)
        
        res_a_ans = self.agent_a.search("Nacrith CDF-24 precision upgrade arithmetic coding floor overhead")
        ans_a = (
            "Actually, integrating Nacrith's CDF-24 is our next step. "
            f"The research proves: {self.agent_a.cite(res_a_ans[0])}. "
            "Upgrading from 2^16 to 2^24 eliminates 75% of quantization overhead in large vocabularies. "
            "I concede we should merge HAQA recipes with Nacrith CDF-24 coding."
        )
        UI.speak(self.agent_a.agent_id, self.agent_a.role, ans_a, self.agent_a.color)

        UI.divider("PHASE 3: COPILOT SYNTHESIS")
        time.sleep(2)
        
        # Final Judgment using Synthesis RPC
        try:
            with grpc.insecure_channel(self.grpc_addr) as channel:
                stub = mesh_pb2_grpc.StrategicMeshStub(channel)
                synthesis = stub.SynthesizeOutputs(mesh_pb2.SynthesisRequest(
                    agent_ids=["student-a", "student-b"],
                    target_goal="Optimal llama.cpp architecture",
                    actions_to_merge=[
                        mesh_pb2.AgentAction(agent_id="student-a", reasoning_chain="Hardware-Aware Quantization"),
                        mesh_pb2.AgentAction(agent_id="student-b", reasoning_chain="Asynchronous Model Partitioning")
                    ]
                ))
                judgment = (
                    f"Decision: {synthesis.synthesized_state}. "
                    "Implementation Insight: The new research demonstrates that hand-optimized llama.cpp is "
                    "no longer the ceiling. We will implement QuantX hardware-aware kernels "
                    "distributed via the ScheInfer async scheduler, specifically using the CDF-24 "
                    "precision upgrade for vocabulary handling."
                )
                UI.speak("COPILOT", "Mesh Arbiter", judgment, UI.YELLOW)
        except Exception as e:
            print(f"Synthesis Error: {e}")

        UI.header("✅ DEBATE CONCLUDED")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Llama.cpp Deep Research Debate")
    parser.add_argument("--grpc-addr", default="localhost:50051", help="gRPC server address")
    args = parser.parse_args()
    
    debate = LlamaDebate(args.grpc_addr)
    debate.run()
