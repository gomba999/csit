# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-24 20:34:07

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
Best observed node throughput: `6211.38` msg/sec with 95% CI [6168.40, 6254.36]
Best sender-completed throughput: `6145.23` msg/sec with 95% CI [6100.29, 6190.16]
Best node CPU: `43.97` % with 95% CI [43.52, 44.42]
Best total CPU: `253.83` % with 95% CI [251.74, 255.92]
Stop reason: refinement narrowed the estimated capacity to offered rates 128000 through 130125

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 6145.23 | [6100.29, 6190.16] | 6211.38 | [6168.40, 6254.36] | 10841.76 | 0.00 | true | 43.97 | [43.52, 44.42] | 253.83 | [251.74, 255.92] | 0 |
| 2 | coarse | 144000 | 25 | 6173.91 | [6155.24, 6192.59] | 6248.64 | [6234.44, 6262.84] | 1183.45 | 0.60 | false | 44.23 | [44.10, 44.35] | 255.03 | [254.52, 255.54] | 0 |
| 3 | coarse | 162000 | 25 | 6163.68 | [6146.79, 6180.56] | 6233.83 | [6216.95, 6250.70] | 1671.63 | 0.36 | false | 44.35 | [44.20, 44.50] | 254.99 | [254.43, 255.54] | 0 |
| 4 | refine | 145000 | 25 | 6089.52 | [5925.57, 6253.46] | 6149.47 | [5983.63, 6315.31] | 161414.70 | -1.00 | false | 43.83 | [42.64, 45.01] | 253.81 | [250.24, 257.38] | 0 |
| 5 | refine | 136500 | 25 | 6179.47 | [6163.13, 6195.81] | 6248.06 | [6236.35, 6259.77] | 804.81 | 0.59 | false | 44.45 | [44.27, 44.62] | 255.52 | [254.85, 256.20] | 0 |
| 6 | refine | 132250 | 25 | 6092.65 | [5914.18, 6271.12] | 6154.48 | [5973.96, 6334.99] | 191243.41 | -0.92 | false | 43.66 | [42.38, 44.95] | 253.44 | [249.55, 257.34] | 0 |
| 7 | refine | 130125 | 25 | 6172.78 | [6151.28, 6194.27] | 6237.23 | [6219.59, 6254.86] | 1825.18 | 0.42 | false | 44.13 | [44.00, 44.25] | 255.01 | [254.50, 255.53] | 0 |
