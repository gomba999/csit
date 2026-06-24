# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-24 21:04:42

**Modes:** write
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

- Modes: `write`
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

### Write Mode

#### Write Clients=1 Payload=16384B

Best offered aggregate rate: `128000` msg/sec
Estimated capacity offered-rate interval: `[128000, 130125]` msg/sec
Best sender write throughput: `6361.11` msg/sec with 95% CI [6335.03, 6387.19]
Best sender-completed throughput: `6361.11` msg/sec with 95% CI [6335.03, 6387.19]
Best node CPU: `44.21` % with 95% CI [44.01, 44.41]
Best total CPU: `185.89` % with 95% CI [185.51, 186.26]
Stop reason: refinement narrowed the estimated capacity to offered rates 128000 through 130125

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Sender Write Throughput | Sender Write Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 6361.11 | [6335.03, 6387.19] | 6361.11 | [6335.03, 6387.19] | 3991.44 | 0.00 | true | 44.21 | [44.01, 44.41] | 185.89 | [185.51, 186.26] | 0 |
| 2 | coarse | 144000 | 25 | 6338.05 | [6293.04, 6383.06] | 6338.05 | [6293.04, 6383.06] | 11891.79 | -0.36 | false | 44.19 | [43.84, 44.53] | 186.11 | [185.59, 186.64] | 0 |
| 3 | coarse | 162000 | 25 | 6367.56 | [6352.43, 6382.70] | 6367.56 | [6352.43, 6382.70] | 1344.42 | 0.10 | false | 44.26 | [44.12, 44.39] | 185.62 | [185.14, 186.10] | 0 |
| 4 | refine | 145000 | 25 | 6352.85 | [6330.72, 6374.99] | 6352.85 | [6330.72, 6374.99] | 2874.89 | -0.13 | false | 44.22 | [44.05, 44.39] | 185.62 | [185.21, 186.02] | 0 |
| 5 | refine | 136500 | 25 | 6362.66 | [6336.49, 6388.83] | 6362.66 | [6336.49, 6388.83] | 4019.53 | 0.02 | false | 44.25 | [44.00, 44.50] | 185.67 | [184.90, 186.44] | 0 |
| 6 | refine | 132250 | 25 | 6345.17 | [6306.49, 6383.86] | 6345.17 | [6306.49, 6383.86] | 8782.27 | -0.25 | false | 44.25 | [43.95, 44.55] | 185.48 | [184.84, 186.12] | 0 |
| 7 | refine | 130125 | 25 | 6325.22 | [6254.70, 6395.74] | 6325.22 | [6254.70, 6395.74] | 29187.92 | -0.56 | false | 44.05 | [43.58, 44.51] | 185.68 | [185.06, 186.29] | 0 |
