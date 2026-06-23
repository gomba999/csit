# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-03 09:21:37

This CI report combines the sink-backed capacity sweeps and the write capacity sweep into one markdown artifact.

## Sink-Backed Modes

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

#### Request-Reply Clients=1 Payload=16384B

Best offered aggregate rate: `1000` msg/sec
Estimated capacity offered-rate interval: `[1000, 1250]` msg/sec
Best observed node throughput: `12.39` msg/sec with 95% CI [12.38, 12.39]
Best sender-completed throughput: `12.15` msg/sec with 95% CI [12.14, 12.17]
Best node CPU: `0.55` % with 95% CI [0.51, 0.59]
Best total CPU: `4.41` % with 95% CI [4.34, 4.48]
Stop reason: refinement narrowed the estimated capacity to offered rates 1000 through 1250

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 1000 | 25 | 12.15 | [12.14, 12.17] | 12.39 | [12.38, 12.39] | 0.00 | 0.00 | true | 0.55 | [0.51, 0.59] | 4.41 | [4.34, 4.48] | 0 |
| 2 | coarse | 2000 | 25 | 12.15 | [12.14, 12.17] | 12.38 | [12.38, 12.39] | 0.00 | -0.01 | false | 0.54 | [0.48, 0.60] | 4.38 | [4.30, 4.47] | 0 |
| 3 | refine | 1500 | 25 | 12.17 | [12.16, 12.18] | 12.38 | [12.38, 12.39] | 0.00 | -0.00 | false | 0.55 | [0.49, 0.61] | 4.39 | [4.31, 4.46] | 0 |
| 4 | refine | 1250 | 25 | 12.16 | [12.15, 12.17] | 12.38 | [12.38, 12.39] | 0.00 | -0.01 | false | 0.55 | [0.50, 0.59] | 4.42 | [4.35, 4.50] | 0 |

## Write Mode

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

