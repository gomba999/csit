# Consensus scenarios (v2)

YAML `ConsensusScenario` plans for the SLIM vs A2A v2 consensus streaming benchmark.
Each scenario defines N agents that run a deterministic **distributed hypothesis convergence**
workload: agents think in parallel, emit findings, and must reach identical local consensus.

## Schema

```yaml
apiVersion: bench.agntcy.io/v2
kind: ConsensusScenario
metadata:
  name: hypothesis-convergence-10agents-10ms
  domain: hypothesis-convergence
spec:
  agents: 10
  thinkTimeMs: 10
  findingEmitDelayMs: 1
  maxRounds: 200
  targetMode: majority
  seed: 42
  valueSpace: 3
agents:
  - id: agent-0
    slimName: agntcy/bench-v2/agent-0
    a2aPort: 9700
    role: coordinator   # A2A only; SLIM treats all agents as peers
  - id: agent-1
    slimName: agntcy/bench-v2/agent-1
    a2aPort: 9711
    role: worker
```

## Generate sweep scenarios

```bash
go run ./tools/gen_scenario \
  -family hypothesis-convergence \
  -agents 10 \
  -think-ms 20 \
  -output plans/sweeps/hypothesis-convergence-10ag-20ms.yaml
```

## Run comparison

```bash
task build
task compare:plan PLAN=hypothesis-convergence-5ag-20ms
task compare:report
```

Primary metric: `consensus_wall_ms` in `reports/results.tsv`.
