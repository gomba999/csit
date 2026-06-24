# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-24 21:07:35

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
Best sender write throughput: `6109.03` msg/sec with 95% CI [5940.49, 6277.58]
Best sender-completed throughput: `6109.03` msg/sec with 95% CI [5940.49, 6277.58]
Best node CPU: `43.60` % with 95% CI [42.37, 44.84]
Best total CPU: `184.59` % with 95% CI [182.53, 186.65]
Stop reason: refinement narrowed the estimated capacity to offered rates 128000 through 130125

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Sender Write Throughput | Sender Write Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 6109.03 | [5940.49, 6277.58] | 6109.03 | [5940.49, 6277.58] | 166717.49 | 0.00 | true | 43.60 | [42.37, 44.84] | 184.59 | [182.53, 186.65] | 0 |
| 2 | coarse | 144000 | 25 | 6147.72 | [6071.41, 6224.03] | 6147.72 | [6071.41, 6224.03] | 34174.74 | 0.63 | false | 43.96 | [43.37, 44.54] | 185.31 | [184.39, 186.22] | 0 |
| 3 | coarse | 162000 | 25 | 6179.50 | [6168.76, 6190.24] | 6179.50 | [6168.76, 6190.24] | 676.98 | 1.15 | false | 44.19 | [44.08, 44.30] | 185.52 | [185.25, 185.79] | 0 |
| 4 | refine | 145000 | 25 | 6136.23 | [6049.95, 6222.51] | 6136.23 | [6049.95, 6222.51] | 43690.26 | 0.45 | false | 43.87 | [43.23, 44.52] | 185.17 | [184.14, 186.19] | 0 |
| 5 | refine | 136500 | 25 | 6130.00 | [6000.38, 6259.62] | 6130.00 | [6000.38, 6259.62] | 98607.02 | 0.34 | false | 43.78 | [42.84, 44.71] | 184.98 | [183.53, 186.43] | 0 |
| 6 | refine | 132250 | 25 | 6194.88 | [6184.59, 6205.18] | 6194.88 | [6184.59, 6205.18] | 621.91 | 1.41 | false | 44.25 | [44.16, 44.35] | 185.77 | [185.44, 186.09] | 0 |
| 7 | refine | 130125 | 25 | 6205.00 | [6190.76, 6219.24] | 6205.00 | [6190.76, 6219.24] | 1190.60 | 1.57 | false | 44.27 | [44.14, 44.40] | 185.87 | [185.51, 186.23] | 0 |
