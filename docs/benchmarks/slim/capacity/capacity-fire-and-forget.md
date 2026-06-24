# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-24 16:34:26

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

Best offered aggregate rate: `136500` msg/sec
Estimated capacity offered-rate interval: `[136500, 138625]` msg/sec
Best observed node throughput: `5741.66` msg/sec with 95% CI [5725.74, 5757.58]
Best sender-completed throughput: `5675.57` msg/sec with 95% CI [5657.96, 5693.18]
Best node CPU: `42.75` % with 95% CI [42.62, 42.89]
Best total CPU: `248.28` % with 95% CI [247.70, 248.86]
Stop reason: refinement narrowed the estimated capacity to offered rates 136500 through 138625

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 5506.40 | [5353.76, 5659.04] | 5569.29 | [5414.59, 5724.00] | 140460.69 | 0.00 | true | 41.64 | [40.49, 42.78] | 244.14 | [240.68, 247.60] | 0 |
| 2 | coarse | 144000 | 25 | 5472.40 | [5400.15, 5544.66] | 5536.57 | [5463.93, 5609.20] | 30965.07 | -0.59 | false | 42.35 | [41.80, 42.91] | 246.56 | [244.99, 248.13] | 0 |
| 3 | coarse | 162000 | 25 | 5571.70 | [5542.04, 5601.36] | 5633.21 | [5604.62, 5661.80] | 4796.31 | 1.15 | false | 43.00 | [42.77, 43.22] | 249.29 | [248.44, 250.14] | 0 |
| 4 | refine | 145000 | 25 | 5534.58 | [5427.73, 5641.44] | 5602.60 | [5494.76, 5710.44] | 68252.41 | 0.60 | false | 42.10 | [41.35, 42.85] | 245.71 | [243.50, 247.92] | 0 |
| 5 | refine | 136500 | 25 | 5675.57 | [5657.96, 5693.18] | 5741.66 | [5725.74, 5757.58] | 1487.84 | 3.09 | true | 42.75 | [42.62, 42.89] | 248.28 | [247.70, 248.86] | 0 |
| 6 | refine | 140750 | 25 | 5580.25 | [5449.72, 5710.78] | 5636.37 | [5504.41, 5768.32] | 102188.19 | -1.83 | false | 42.39 | [41.37, 43.40] | 247.81 | [244.83, 250.80] | 0 |
| 7 | refine | 138625 | 25 | 5539.39 | [5434.41, 5644.38] | 5600.36 | [5494.31, 5706.40] | 65999.40 | -2.46 | false | 41.97 | [41.21, 42.73] | 245.40 | [243.24, 247.56] | 0 |
