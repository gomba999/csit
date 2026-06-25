# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-25 13:36:39

This CI report combines the sink-backed capacity sweeps and the write capacity sweep into one markdown artifact.

## Sink-Backed Modes

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

#### Request-Reply Clients=1 Payload=16384B

Best offered aggregate rate: `1000` msg/sec
Estimated capacity offered-rate interval: `[1000, 1250]` msg/sec
Best observed node throughput: `12.37` msg/sec with 95% CI [12.36, 12.37]
Best sender-completed throughput: `12.15` msg/sec with 95% CI [12.13, 12.16]
Best node CPU: `0.70` % with 95% CI [0.66, 0.74]
Best total CPU: `4.51` % with 95% CI [4.44, 4.59]
Stop reason: refinement narrowed the estimated capacity to offered rates 1000 through 1250

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 1000 | 25 | 12.15 | [12.13, 12.16] | 12.37 | [12.36, 12.37] | 0.00 | 0.00 | true | 0.70 | [0.66, 0.74] | 4.51 | [4.44, 4.59] | 0 |
| 2 | coarse | 2000 | 25 | 12.13 | [12.12, 12.15] | 12.37 | [12.36, 12.37] | 0.00 | -0.03 | false | 0.72 | [0.67, 0.76] | 4.56 | [4.49, 4.63] | 0 |
| 3 | refine | 1500 | 25 | 12.15 | [12.13, 12.16] | 12.37 | [12.36, 12.38] | 0.00 | 0.01 | false | 0.69 | [0.64, 0.74] | 4.46 | [4.36, 4.55] | 0 |
| 4 | refine | 1250 | 25 | 12.15 | [12.13, 12.16] | 12.37 | [12.37, 12.38] | 0.00 | 0.03 | false | 0.66 | [0.61, 0.70] | 4.36 | [4.28, 4.44] | 0 |

## Write Mode

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

