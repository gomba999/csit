# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-25 13:03:08

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
Best observed node throughput: `6296.43` msg/sec with 95% CI [6227.56, 6365.30]
Best sender-completed throughput: `6230.35` msg/sec with 95% CI [6163.61, 6297.09]
Best node CPU: `42.78` % with 95% CI [42.28, 43.28]
Best total CPU: `243.67` % with 95% CI [241.71, 245.62]
Stop reason: refinement narrowed the estimated capacity to offered rates 128000 through 130125

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 6230.35 | [6163.61, 6297.09] | 6296.43 | [6227.56, 6365.30] | 27835.33 | 0.00 | true | 42.78 | [42.28, 43.28] | 243.67 | [241.71, 245.62] | 0 |
| 2 | coarse | 144000 | 25 | 6323.43 | [6303.35, 6343.52] | 6388.42 | [6368.88, 6407.96] | 2240.96 | 1.46 | false | 43.68 | [43.55, 43.82] | 248.24 | [247.63, 248.86] | 0 |
| 3 | coarse | 162000 | 25 | 6355.81 | [6335.91, 6375.71] | 6421.09 | [6404.43, 6437.75] | 1629.05 | 1.98 | false | 44.04 | [43.90, 44.18] | 249.40 | [248.81, 250.00] | 0 |
| 4 | refine | 145000 | 25 | 6365.37 | [6342.67, 6388.06] | 6430.79 | [6409.29, 6452.29] | 2713.27 | 2.13 | false | 44.02 | [43.86, 44.17] | 249.00 | [248.42, 249.59] | 0 |
| 5 | refine | 136500 | 25 | 6346.95 | [6320.28, 6373.61] | 6424.23 | [6400.37, 6448.09] | 3341.75 | 2.03 | false | 43.83 | [43.65, 44.02] | 248.08 | [247.30, 248.86] | 0 |
| 6 | refine | 132250 | 25 | 6334.91 | [6310.85, 6358.97] | 6397.87 | [6374.69, 6421.04] | 3152.55 | 1.61 | false | 43.74 | [43.54, 43.94] | 248.09 | [247.08, 249.09] | 0 |
| 7 | refine | 130125 | 25 | 6284.99 | [6267.07, 6302.90] | 6351.79 | [6335.54, 6368.04] | 1549.58 | 0.88 | false | 43.38 | [43.25, 43.51] | 247.39 | [246.78, 248.01] | 0 |
