# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-24 20:31:11

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
Best observed node throughput: `6366.36` msg/sec with 95% CI [6310.61, 6422.11]
Best sender-completed throughput: `6295.70` msg/sec with 95% CI [6236.76, 6354.64]
Best node CPU: `43.64` % with 95% CI [43.08, 44.20]
Best total CPU: `247.31` % with 95% CI [244.95, 249.67]
Stop reason: refinement narrowed the estimated capacity to offered rates 128000 through 130125

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 6295.70 | [6236.76, 6354.64] | 6366.36 | [6310.61, 6422.11] | 18242.85 | 0.00 | true | 43.64 | [43.08, 44.20] | 247.31 | [244.95, 249.67] | 0 |
| 2 | coarse | 144000 | 25 | 6289.71 | [6260.61, 6318.80] | 6360.58 | [6334.70, 6386.45] | 3928.41 | -0.09 | false | 43.98 | [43.74, 44.21] | 248.27 | [247.41, 249.14] | 0 |
| 3 | coarse | 162000 | 25 | 6254.58 | [6239.70, 6269.47] | 6321.33 | [6309.43, 6333.24] | 831.28 | -0.71 | false | 43.65 | [43.52, 43.79] | 247.52 | [246.91, 248.12] | 0 |
| 4 | refine | 145000 | 25 | 6236.00 | [6215.21, 6256.80] | 6314.72 | [6296.67, 6332.76] | 1911.48 | -0.81 | false | 43.53 | [43.35, 43.71] | 246.90 | [246.00, 247.79] | 0 |
| 5 | refine | 136500 | 25 | 6252.63 | [6237.02, 6268.24] | 6317.24 | [6304.79, 6329.69] | 909.18 | -0.77 | false | 43.45 | [43.32, 43.59] | 246.50 | [245.92, 247.07] | 0 |
| 6 | refine | 132250 | 25 | 6264.51 | [6248.15, 6280.87] | 6343.42 | [6329.22, 6357.63] | 1184.38 | -0.36 | false | 43.77 | [43.64, 43.90] | 248.10 | [247.46, 248.75] | 0 |
| 7 | refine | 130125 | 25 | 6257.21 | [6237.19, 6277.23] | 6329.92 | [6310.73, 6349.10] | 2159.74 | -0.57 | false | 43.80 | [43.69, 43.91] | 248.59 | [248.11, 249.06] | 0 |
