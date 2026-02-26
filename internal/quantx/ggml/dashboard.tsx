import React from 'react';
import { Stack, Box, Text, Heading, Badge, Card } from "@voxel51/voodo";

interface ChainStep {
  agent: string;
  action: string;
  hardware: string;
}

interface AgentPerformance {
  id: string;
  avgLatency: number;
  throughput: number; // tokens/sec
}

interface VextraDashboardProps {
  reasoningChain: ChainStep[];
  agentStats: AgentPerformance[];
}

/*
  Mapping Hint: 
  const agentStats = benchmarkData.benchmarks.reduce((acc, b) => {
    if (!acc[b.agent_id]) acc[b.agent_id] = { id: b.agent_id, totalLat: 0, count: 0, tokens: 0 };
    acc[b.agent_id].totalLat += b.latency_ms;
    acc[b.agent_id].tokens += b.tokens_used;
    acc[b.agent_id].count++;
    return acc;
  }, {}).map(a => ({ id: a.id, avgLatency: a.totalLat/a.count, throughput: a.tokens / (a.totalLat/1000) }));
*/

export const VextraDashboard: React.FC<VextraDashboardProps> = ({ reasoningChain, agentStats }) => {
  return (
    <Box p="md">
      <Stack spacing="lg">
        <Heading size="md">󱐋 Vextra Mesh Dashboard</Heading>
        
        <Stack spacing="sm">
          <Text weight="bold">Live Reasoning Chain</Text>
          {reasoningChain.map((step, i) => (
            <Card key={i} p="sm">
              <Stack direction="row" align="center" spacing="md" justify="space-between">
                <Stack direction="row" align="center" spacing="md">
                  <Badge variant="primary">{step.agent}</Badge>
                  <Text>➜</Text>
                  <Stack>
                    <Text size="sm" weight="bold">{step.action}</Text>
                    <Text size="xs" color="neutral.darker">{step.hardware}</Text>
                  </Stack>
                </Stack>
                <Badge color={step.hardware.includes("GPU") ? "success" : "neutral"}>
                  {step.hardware}
                </Badge>
              </Stack>
            </Card>
          ))}
        </Stack>

        <Stack spacing="sm">
          <Text weight="bold">Individual Agent Performance</Text>
          <Box borderRadius="sm" overflow="hidden" border="1px solid" borderColor="neutral.light">
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '12px' }}>
              <thead style={{ backgroundColor: '#f4f4f4' }}>
                <tr>
                  <th style={{ textAlign: 'left', padding: '8px' }}>Agent ID</th>
                  <th style={{ textAlign: 'right', padding: '8px' }}>Avg Latency</th>
                  <th style={{ textAlign: 'right', padding: '8px' }}>Tokens/Sec</th>
                </tr>
              </thead>
              <tbody>
                {agentStats.map(stat => (
                  <tr key={stat.id} style={{ borderTop: '1px solid #eee' }}>
                    <td style={{ padding: '8px' }}>{stat.id}</td>
                    <td style={{ padding: '8px', textAlign: 'right' }}>{stat.avgLatency.toFixed(1)}ms</td>
                    <td style={{ padding: '8px', textAlign: 'right' }}>{stat.throughput.toFixed(1)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </Box>
        </Stack>

        <Box p="sm" borderRadius="sm" backgroundColor="neutral.light">
          <Stack direction="row" spacing="xl">
            <Stack>
              <Text size="xs">Total Mesh Throughput</Text>
              <Text weight="bold" color="success.main">2,450 tokens/min</Text>
            </Stack>
            <Stack>
              <Text size="xs">Hardware Utilization</Text>
              <Text weight="bold" color="success.main">82% (GPU)</Text>
            </Stack>
          </Stack>
        </Box>
      </Stack>
    </Box>
  );
};
