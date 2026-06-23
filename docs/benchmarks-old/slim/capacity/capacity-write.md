# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-03 09:21:37

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
Best sender write throughput: `5960.55` msg/sec with 95% CI [5950.90, 5970.20]
Best sender-completed throughput: `5960.55` msg/sec with 95% CI [5950.90, 5970.20]
Best node CPU: `43.69` % with 95% CI [43.58, 43.80]
Best total CPU: `184.82` % with 95% CI [184.52, 185.11]
Stop reason: refinement narrowed the estimated capacity to offered rates 128000 through 130125

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Sender Write Throughput | Sender Write Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 5960.55 | [5950.90, 5970.20] | 5960.55 | [5950.90, 5970.20] | 546.06 | 0.00 | true | 43.69 | [43.58, 43.80] | 184.82 | [184.52, 185.11] | 0 |
| 2 | coarse | 144000 | 25 | 5745.27 | [5482.91, 6007.62] | 5745.27 | [5482.91, 6007.62] | 403950.69 | -3.61 | false | 41.99 | [40.00, 43.97] | 182.27 | [179.26, 185.27] | 0 |
| 3 | coarse | 162000 | 25 | 5971.06 | [5955.26, 5986.85] | 5971.06 | [5955.26, 5986.85] | 1464.95 | 0.18 | false | 43.67 | [43.52, 43.81] | 184.83 | [184.41, 185.25] | 0 |
| 4 | refine | 145000 | 25 | 5912.44 | [5794.14, 6030.73] | 5912.44 | [5794.14, 6030.73] | 82123.67 | -0.81 | false | 43.14 | [42.23, 44.04] | 184.11 | [182.64, 185.58] | 0 |
| 5 | refine | 136500 | 25 | 5873.71 | [5700.47, 6046.96] | 5873.71 | [5700.47, 6046.96] | 176148.28 | -1.46 | false | 42.85 | [41.54, 44.16] | 183.55 | [181.58, 185.51] | 0 |
| 6 | refine | 132250 | 25 | 5979.38 | [5965.37, 5993.39] | 5979.38 | [5965.37, 5993.39] | 1151.73 | 0.32 | false | 43.61 | [43.50, 43.71] | 184.98 | [184.69, 185.27] | 0 |
| 7 | refine | 130125 | 25 | 5887.21 | [5697.47, 6076.95] | 5887.21 | [5697.47, 6076.95] | 211285.50 | -1.23 | false | 43.00 | [41.58, 44.43] | 183.80 | [181.60, 186.01] | 0 |
