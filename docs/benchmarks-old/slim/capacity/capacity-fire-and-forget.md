# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-03 08:48:09

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
Best observed node throughput: `5962.24` msg/sec with 95% CI [5930.07, 5994.41]
Best sender-completed throughput: `5890.64` msg/sec with 95% CI [5854.51, 5926.76]
Best node CPU: `43.00` % with 95% CI [42.61, 43.40]
Best total CPU: `251.58` % with 95% CI [249.78, 253.39]
Stop reason: refinement narrowed the estimated capacity to offered rates 128000 through 130125

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 5890.64 | [5854.51, 5926.76] | 5962.24 | [5930.07, 5994.41] | 6072.32 | 0.00 | true | 43.00 | [42.61, 43.40] | 251.58 | [249.78, 253.39] | 0 |
| 2 | coarse | 144000 | 25 | 5932.28 | [5920.14, 5944.42] | 5989.34 | [5978.53, 6000.15] | 686.25 | 0.45 | false | 43.51 | [43.41, 43.61] | 253.64 | [253.20, 254.07] | 0 |
| 3 | coarse | 162000 | 25 | 5959.88 | [5946.23, 5973.52] | 6020.80 | [6007.06, 6034.53] | 1107.31 | 0.98 | false | 43.53 | [43.44, 43.63] | 253.62 | [253.25, 254.00] | 0 |
| 4 | refine | 145000 | 25 | 5962.15 | [5944.25, 5980.04] | 6026.74 | [6010.60, 6042.87] | 1527.73 | 1.08 | false | 43.43 | [43.31, 43.56] | 253.29 | [252.65, 253.92] | 0 |
| 5 | refine | 136500 | 25 | 5978.14 | [5965.78, 5990.51] | 6037.17 | [6025.92, 6048.42] | 743.11 | 1.26 | false | 43.56 | [43.43, 43.69] | 253.89 | [253.42, 254.37] | 0 |
| 6 | refine | 132250 | 25 | 5954.56 | [5940.68, 5968.44] | 6026.86 | [6013.86, 6039.85] | 990.68 | 1.08 | false | 43.51 | [43.41, 43.61] | 253.32 | [252.86, 253.78] | 0 |
| 7 | refine | 130125 | 25 | 5967.13 | [5947.38, 5986.88] | 6027.79 | [6011.16, 6044.42] | 1623.44 | 1.10 | false | 43.55 | [43.38, 43.71] | 253.73 | [253.08, 254.39] | 0 |
