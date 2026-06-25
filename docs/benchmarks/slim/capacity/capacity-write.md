# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-25 13:36:39

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
Best sender write throughput: `6318.51` msg/sec with 95% CI [6295.12, 6341.89]
Best sender-completed throughput: `6318.51` msg/sec with 95% CI [6295.12, 6341.89]
Best node CPU: `43.85` % with 95% CI [43.66, 44.04]
Best total CPU: `184.44` % with 95% CI [183.95, 184.93]
Stop reason: refinement narrowed the estimated capacity to offered rates 128000 through 130125

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Sender Write Throughput | Sender Write Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 6318.51 | [6295.12, 6341.89] | 6318.51 | [6295.12, 6341.89] | 3209.68 | 0.00 | true | 43.85 | [43.66, 44.04] | 184.44 | [183.95, 184.93] | 0 |
| 2 | coarse | 144000 | 25 | 6278.62 | [6251.87, 6305.37] | 6278.62 | [6251.87, 6305.37] | 4199.57 | -0.63 | false | 43.64 | [43.50, 43.78] | 183.80 | [183.31, 184.29] | 0 |
| 3 | coarse | 162000 | 25 | 6297.25 | [6273.00, 6321.50] | 6297.25 | [6273.00, 6321.50] | 3451.29 | -0.34 | false | 43.70 | [43.52, 43.88] | 183.97 | [183.42, 184.51] | 0 |
| 4 | refine | 145000 | 25 | 6313.72 | [6298.05, 6329.40] | 6313.72 | [6298.05, 6329.40] | 1442.31 | -0.08 | false | 43.76 | [43.57, 43.94] | 184.02 | [183.48, 184.55] | 0 |
| 5 | refine | 136500 | 25 | 6328.24 | [6316.60, 6339.89] | 6328.24 | [6316.60, 6339.89] | 795.85 | 0.15 | false | 43.80 | [43.66, 43.94] | 183.93 | [183.62, 184.24] | 0 |
| 6 | refine | 132250 | 25 | 6337.25 | [6320.51, 6354.00] | 6337.25 | [6320.51, 6354.00] | 1645.28 | 0.30 | false | 43.89 | [43.75, 44.03] | 184.19 | [183.88, 184.51] | 0 |
| 7 | refine | 130125 | 25 | 6351.55 | [6317.55, 6385.55] | 6351.55 | [6317.55, 6385.55] | 6783.29 | 0.52 | false | 44.00 | [43.72, 44.28] | 185.42 | [184.87, 185.96] | 0 |
