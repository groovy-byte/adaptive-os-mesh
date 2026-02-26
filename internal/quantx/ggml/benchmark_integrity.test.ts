import * as fs from 'fs';
import * as path from 'path';

describe('Benchmark Data Integrity', () => {
  const resultsPath = path.resolve(__dirname, '../../../benchmark_results.json');

  it('verifies that benchmark_results.json exists and is valid JSON', () => {
    expect(fs.existsSync(resultsPath)).toBe(true);
    const data = JSON.parse(fs.readFileSync(resultsPath, 'utf8'));
    expect(data).toHaveProperty('benchmarks');
    expect(data).toHaveProperty('contribution_matrix');
    expect(data).toHaveProperty('agent_logs');
    expect(Array.isArray(data.benchmarks)).toBe(true);
    expect(typeof data.contribution_matrix).toBe('object');
    expect(typeof data.agent_logs).toBe('object');
  });

  describe('Agent Audit Log Integrity', () => {
    it('verifies tool call counts and failure log structure', () => {
      const data = JSON.parse(fs.readFileSync(resultsPath, 'utf8'));
      const logs = data.agent_logs;

      Object.keys(logs).forEach(agent => {
        const log = logs[agent];
        expect(typeof log.tool_calls).toBe('number');
        expect(Array.isArray(log.failed_tasks)).toBe(true);
        expect(log.failed_tasks.length).toBeLessThanOrEqual(5);
        log.failed_tasks.forEach((f: any) => expect(typeof f).toBe('string'));
      });
    });
  });

  describe('Value of Contribution (VoC) Integrity', () => {
    it('verifies the structure of the contribution matrix', () => {
      const data = JSON.parse(fs.readFileSync(resultsPath, 'utf8'));
      const matrix = data.contribution_matrix;
      
      Object.keys(matrix).forEach(source => {
        const targets = matrix[source];
        expect(typeof targets).toBe('object');
        
        Object.keys(targets).forEach(target => {
          const score = targets[target];
          expect(typeof score).toBe('number');
          expect(score).toBeGreaterThanOrEqual(0);
          expect(score).toBeLessThanOrEqual(1.0);
        });
      });
    });
  });

  const getBenchmarks = () => {
    if (!fs.existsSync(resultsPath)) return [];
    return JSON.parse(fs.readFileSync(resultsPath, 'utf8')).benchmarks;
  };

  describe('Hardware Routing Logic', () => {
    getBenchmarks().forEach((b: any, i: number) => {
      it(`verifies routing for task ${i} (${b.agent_id} - ${b.kv_cache_mb}MB)`, () => {
        if (b.kv_cache_mb <= 8) {
          expect(b.hardware_path).toBe('CPU_AVX2');
        } else if (b.kv_cache_mb >= 16) {
          expect(b.hardware_path).toMatch(/GPU|AVX512/);
        }
      });
    });
  });

  describe('Latency Consistency', () => {
    getBenchmarks().forEach((b: any, i: number) => {
      it(`verifies latency bounds for task ${i} (${b.hardware_path})`, () => {
        if (b.hardware_path === 'CPU_AVX2') {
          expect(b.latency_ms).toBeGreaterThanOrEqual(40);
          expect(b.latency_ms).toBeLessThanOrEqual(80);
        } else if (b.hardware_path === 'CPU_AVX512') {
          expect(b.latency_ms).toBeGreaterThanOrEqual(30);
          expect(b.latency_ms).toBeLessThanOrEqual(50);
        } else if (b.hardware_path.includes('GPU')) {
          expect(b.latency_ms).toBeGreaterThanOrEqual(100);
          expect(b.latency_ms).toBeLessThanOrEqual(150);
        }
      });
    });
  });

  describe('Field Integrity', () => {
    getBenchmarks().forEach((b: any, i: number) => {
      it(`checks mandatory fields for task ${i}`, () => {
        expect(b).toHaveProperty('agent_id');
        expect(b).toHaveProperty('latency_ms');
        expect(b).toHaveProperty('tokens_used');
        expect(b).toHaveProperty('hardware_path');
        expect(b).toHaveProperty('status', 'success');
      });
    });
  });
});
