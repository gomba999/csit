# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-24 16:38:03

**Modes:** fire-and-forget
**Clients:** 1
**Sizes:** 16384
**Start Rate:** 128000
**Max Rate:** 176000
**Growth Factor:** 1.12
**Plateau Threshold:** 3.00%
**Plateau Steps:** 2
**Max Steps:** 6
**Repeats Per Sweep Step:** 25

## Adaptive Capacity Sweep

This sweep first increases the configured send rate geometrically to find the saturation region, then performs midpoint refinement to narrow the offered-rate interval that saturates the node.
Results are reported separately for each fixed `(mode, clients, payload)` case. The reported rate is the aggregate offered load across all clients in that case. For request-reply and fire-and-forget, effective throughput is sink-observed total node throughput. For write mode, effective throughput is sender-completed write throughput because no responder is running.

- Modes: `fire-and-forget`
- Clients: `1`
- Sizes: `16384` bytes
- Start rate: `128000` msg/sec
- Max rate: `176000` msg/sec (0 means unbounded by rate cap)
- Growth factor: `1.12`
- Plateau threshold: `3.00%` effective throughput gain
- Plateau steps: `2`
- Max steps: `6`
- Repeats per sweep step: `25`
- Refinement steps after coarse sweep: `4`
- Minimum offered-rate interval after refinement: `1000` msg/sec

### Sink-Backed Modes

#### Fire-And-Forget Clients=1 Payload=16384B

Best offered aggregate rate: `128000` msg/sec
Estimated capacity offered-rate interval: `[128000, 130125]` msg/sec
Best observed node throughput: `7386.03` msg/sec with 95% CI [7343.40, 7428.65]
Best sender-completed throughput: `7308.84` msg/sec with 95% CI [7262.36, 7355.33]
Best node CPU: `41.47` % with 95% CI [41.11, 41.83]
Best total CPU: `247.44` % with 95% CI [245.91, 248.98]
Stop reason: refinement narrowed the estimated capacity to offered rates 128000 through 130125

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 7308.84 | [7262.36, 7355.33] | 7386.03 | [7343.40, 7428.65] | 10661.86 | 0.00 | true | 41.47 | [41.11, 41.83] | 247.44 | [245.91, 248.98] | 0 |
| 2 | coarse | 144000 | 25 | 7391.44 | [7368.02, 7414.87] | 7469.64 | [7450.80, 7488.49] | 2083.28 | 1.13 | false | 41.95 | [41.84, 42.05] | 249.43 | [249.02, 249.83] | 0 |
| 3 | coarse | 162000 | 25 | 7378.01 | [7362.86, 7393.16] | 7465.60 | [7449.99, 7481.21] | 1430.77 | 1.08 | false | 41.90 | [41.80, 42.00] | 248.96 | [248.46, 249.46] | 0 |
| 4 | refine | 145000 | 25 | 7362.34 | [7346.47, 7378.22] | 7442.39 | [7428.08, 7456.71] | 1202.51 | 0.76 | false | 41.74 | [41.61, 41.87] | 248.21 | [247.56, 248.85] | 0 |
| 5 | refine | 136500 | 25 | 7385.30 | [7367.31, 7403.30] | 7472.54 | [7458.80, 7486.27] | 1106.94 | 1.17 | false | 41.88 | [41.79, 41.98] | 248.86 | [248.40, 249.33] | 0 |
| 6 | refine | 132250 | 25 | 7379.39 | [7359.48, 7399.30] | 7476.20 | [7461.50, 7490.89] | 1267.91 | 1.22 | false | 41.83 | [41.70, 41.96] | 248.66 | [248.11, 249.22] | 0 |
| 7 | refine | 130125 | 25 | 7393.22 | [7376.24, 7410.21] | 7471.29 | [7455.33, 7487.25] | 1494.41 | 1.15 | false | 41.90 | [41.78, 42.01] | 248.97 | [248.46, 249.48] | 0 |
