# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-24 17:11:29

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
Best sender write throughput: `7184.99` msg/sec with 95% CI [7138.94, 7231.04]
Best sender-completed throughput: `7184.99` msg/sec with 95% CI [7138.94, 7231.04]
Best node CPU: `41.80` % with 95% CI [41.56, 42.04]
Best total CPU: `182.33` % with 95% CI [181.72, 182.94]
Stop reason: refinement narrowed the estimated capacity to offered rates 128000 through 130125

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Sender Write Throughput | Sender Write Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 7184.99 | [7138.94, 7231.04] | 7184.99 | [7138.94, 7231.04] | 12444.53 | 0.00 | true | 41.80 | [41.56, 42.04] | 182.33 | [181.72, 182.94] | 0 |
| 2 | coarse | 144000 | 25 | 7233.79 | [7188.88, 7278.71] | 7233.79 | [7188.88, 7278.71] | 11839.15 | 0.68 | false | 41.94 | [41.81, 42.07] | 183.12 | [182.74, 183.49] | 0 |
| 3 | coarse | 162000 | 25 | 7078.69 | [7038.59, 7118.79] | 7078.69 | [7038.59, 7118.79] | 9436.55 | -1.48 | false | 42.30 | [42.09, 42.51] | 183.51 | [182.95, 184.07] | 0 |
| 4 | refine | 145000 | 25 | 7265.21 | [7241.48, 7288.94] | 7265.21 | [7241.48, 7288.94] | 3305.62 | 1.12 | false | 41.81 | [41.71, 41.91] | 183.01 | [182.65, 183.37] | 0 |
| 5 | refine | 136500 | 25 | 7086.90 | [7044.20, 7129.59] | 7086.90 | [7044.20, 7129.59] | 10697.69 | -1.37 | false | 42.09 | [41.92, 42.27] | 183.24 | [182.74, 183.73] | 0 |
| 6 | refine | 132250 | 25 | 7082.58 | [7053.87, 7111.30] | 7082.58 | [7053.87, 7111.30] | 4838.22 | -1.43 | false | 41.95 | [41.78, 42.12] | 182.84 | [182.34, 183.35] | 0 |
| 7 | refine | 130125 | 25 | 7277.24 | [7257.34, 7297.13] | 7277.24 | [7257.34, 7297.13] | 2323.16 | 1.28 | false | 41.81 | [41.70, 41.93] | 182.58 | [182.25, 182.92] | 0 |
