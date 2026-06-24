# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-24 21:04:42

This CI report combines the sink-backed capacity sweeps and the write capacity sweep into one markdown artifact.

## Sink-Backed Modes

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

#### Request-Reply Clients=1 Payload=16384B

Best offered aggregate rate: `1000` msg/sec
Estimated capacity offered-rate interval: `[1000, 1250]` msg/sec
Best observed node throughput: `12.38` msg/sec with 95% CI [12.38, 12.38]
Best sender-completed throughput: `12.16` msg/sec with 95% CI [12.14, 12.17]
Best node CPU: `0.64` % with 95% CI [0.59, 0.69]
Best total CPU: `4.32` % with 95% CI [4.25, 4.39]
Stop reason: refinement narrowed the estimated capacity to offered rates 1000 through 1250

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 1000 | 25 | 12.16 | [12.14, 12.17] | 12.38 | [12.38, 12.38] | 0.00 | 0.00 | true | 0.64 | [0.59, 0.69] | 4.32 | [4.25, 4.39] | 0 |
| 2 | coarse | 2000 | 25 | 12.15 | [12.14, 12.16] | 12.38 | [12.37, 12.38] | 0.00 | -0.02 | false | 0.65 | [0.60, 0.70] | 4.33 | [4.27, 4.39] | 0 |
| 3 | refine | 1500 | 25 | 12.16 | [12.15, 12.17] | 12.38 | [12.38, 12.38] | 0.00 | 0.00 | false | 0.62 | [0.57, 0.67] | 4.26 | [4.21, 4.31] | 0 |
| 4 | refine | 1250 | 25 | 12.15 | [12.14, 12.17] | 12.38 | [12.38, 12.38] | 0.00 | 0.00 | false | 0.62 | [0.58, 0.66] | 4.25 | [4.20, 4.31] | 0 |

## Write Mode

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

