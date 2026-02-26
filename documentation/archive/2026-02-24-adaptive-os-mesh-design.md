# Design Document: Adaptive OS-Knowledgeable Mesh (Maestro)
**Date**: 2026-02-24
**Status**: Approved

## 1. Problem Statement
The objective is to implement an "OS-Knowledgeable" multi-agent mesh capable of semantic search across Student A/B shared context and a 25-paper research corpus. The system must overcome the "Strategic-Operational Gap" by supporting both high-level reasoning and fast task execution, while implementing "counterbalancing" mechanisms to prevent resource exhaustion (memory spikes).

## 2. Requirements
### Functional
- **Unified Semantic Search**: Search across live SQLite state (`agent_mesh.db`) and research papers.
- **Dual-Protocol Control Plane**: gRPC for structured "Strategic" work; NATS for fast "Operational" work.
- **Dynamic Role Evolution**: Agents (Student A/B) switch between Boss/Task roles based on workload.
- **Adaptive Protocol Bridge**: Protobuf-first contracts with a JSON-based "Flex-Layer" for adaptability.

### Non-Functional
- **Resilience**: A pure Go core to maintain heartbeats during Python memory spikes.
- **Counterbalancing**: Soft-throttling based on the "Value of Coordination" (semantic novelty).
- **Adaptability**: Evolutionary topology switching between Parallel and Hierarchical modes.

## 3. Architecture
### Components
1. **Vextra Mesh Controller (Go)**: Protocol bridge and mesh backbone.
2. **Shared Memory OS (SQLite)**: Local operational memory mirrored to Gemini Cloud.
3. **Copilot Intermediary**: Mesh supervisor and topology arbiter.
4. **Fluid Nodes (Student A/B)**: Adaptive agents with gRPC/NATS dual-stack.

### Data Flow
- **Fast Path**: Task Agent -> NATS -> Go Core -> SQLite.
- **Reasoning Path**: Boss Agent -> gRPC -> Go Core -> Gemini Research Store.
- **Sync Path**: Asynchronous SQLite-to-Cloud mirroring via Go worker.

## 4. Agent Team
- **Student A & B**: Fluid Mesh Nodes. Dynamic promotion to Strategic roles (gRPC) or demotion to Operational roles (NATS).
- **Copilot Agent**: The Supervisor. Executes the "Soft-Throttle" and manages state reconstitution.
- **Maestro (TechLead)**: Orchestrates the build and validation of the Go/Python ecosystem.

## 5. Risk Assessment & Mitigation
- **Risk**: Agent hangs during complex OS tasks.
- **Mitigation**: **State Reconstitution**. The Copilot agent uses NATS JetStream logs to rebuild the state of failed nodes.
- **Risk**: 15.4x memory spikes during tool calls.
- **Mitigation**: eBPF-inspired intent-driven resource adaptation (via AgentCgroup pattern).

## 6. Success Criteria
- [ ] Successful semantic query cross-referencing research papers and live SQLite logs.
- [ ] Go core maintains 100% uptime during simulated Python heavy-load bursts.
- [ ] Verified role-switching (NATS -> gRPC) triggered by task complexity.
