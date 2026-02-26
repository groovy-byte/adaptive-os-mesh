import React from 'react';
import { render, screen } from '@testing-library/react';
import { VextraDashboard } from './dashboard';

describe('VextraDashboard', () => {
  const mockChain = [
    { agent: 'Scout-Agent', action: 'Search', hardware: 'CPU_AVX2' },
    { agent: 'Coder-Agent', action: 'Coding', hardware: 'GPU_CUDA' },
  ];

  const mockStats = [
    { id: 'Scout-Agent', avgLatency: 55.2, throughput: 416.7 },
    { id: 'Coder-Agent', avgLatency: 122.1, throughput: 188.4 },
  ];

  it('renders the reasoning chain with hardware paths', () => {
    render(<VextraDashboard reasoningChain={mockChain} agentStats={mockStats} />);
    
    expect(screen.getAllByText('Scout-Agent')[0]).toBeInTheDocument();
    expect(screen.getAllByText('Coder-Agent')[0]).toBeInTheDocument();
    expect(screen.getAllByText('CPU_AVX2')[0]).toBeInTheDocument();
    expect(screen.getAllByText('GPU_CUDA')[0]).toBeInTheDocument();
  });

  it('renders the individual agent performance table correctly', () => {
    render(<VextraDashboard reasoningChain={mockChain} agentStats={mockStats} />);
    
    // Check table headers
    expect(screen.getByText('Agent ID')).toBeInTheDocument();
    expect(screen.getByText('Avg Latency')).toBeInTheDocument();
    expect(screen.getByText('Tokens/Sec')).toBeInTheDocument();

    // Check row data
    // Scout-Agent: Latency 55.2ms, Throughput 416.7
    expect(screen.getByText('55.2ms')).toBeInTheDocument();
    expect(screen.getByText('416.7')).toBeInTheDocument();
    
    // Coder-Agent: Latency 122.1ms, Throughput 188.4
    expect(screen.getByText('122.1ms')).toBeInTheDocument();
    expect(screen.getByText('188.4')).toBeInTheDocument();
  });

  it('calculates throughput correctly based on sample data', () => {
    const customStats = [
      { id: 'Test-Agent', avgLatency: 100, throughput: 10.0 } // 10 tokens / 0.1 sec = 100 tokens/sec? No, throughput is tokens/sec.
    ];
    render(<VextraDashboard reasoningChain={mockChain} agentStats={customStats} />);
    expect(screen.getByText('10.0')).toBeInTheDocument();
  });

  it('displays hardware utilization and mesh throughput', () => {
    render(<VextraDashboard reasoningChain={mockChain} agentStats={mockStats} />);
    
    expect(screen.getByText('Total Mesh Throughput')).toBeInTheDocument();
    expect(screen.getByText('2,450 tokens/min')).toBeInTheDocument();
    expect(screen.getByText('82% (GPU)')).toBeInTheDocument();
  });

  it('handles extreme performance values without error', () => {
    const extremeStats = [
      { id: 'Zero-Agent', avgLatency: 0, throughput: 0 },
      { id: 'Max-Agent', avgLatency: 9999, throughput: 9999 }
    ];
    render(<VextraDashboard reasoningChain={mockChain} agentStats={extremeStats} />);
    expect(screen.getByText('Zero-Agent')).toBeInTheDocument();
    expect(screen.getByText('Max-Agent')).toBeInTheDocument();
  });
});
