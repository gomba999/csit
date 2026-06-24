# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-24 21:07:35

This CI report combines the sink-backed capacity sweeps and the write capacity sweep into one markdown artifact.

## Sink-Backed Modes

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

#### Request-Reply Clients=1 Payload=16384B

Best offered aggregate rate: `1000` msg/sec
Estimated capacity offered-rate interval: `[1000, 1250]` msg/sec
Best observed node throughput: `12.39` msg/sec with 95% CI [12.38, 12.39]
Best sender-completed throughput: `12.16` msg/sec with 95% CI [12.15, 12.18]
Best node CPU: `0.52` % with 95% CI [0.47, 0.57]
Best total CPU: `4.16` % with 95% CI [4.09, 4.22]
Stop reason: refinement narrowed the estimated capacity to offered rates 1000 through 1250

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 1000 | 25 | 12.16 | [12.15, 12.18] | 12.39 | [12.38, 12.39] | 0.00 | 0.00 | true | 0.52 | [0.47, 0.57] | 4.16 | [4.09, 4.22] | 0 |
| 2 | coarse | 2000 | 25 | 12.16 | [12.15, 12.18] | 12.39 | [12.38, 12.39] | 0.00 | -0.01 | false | 0.53 | [0.47, 0.60] | 4.18 | [4.09, 4.26] | 0 |
| 3 | refine | 1500 | 25 | 12.16 | [12.15, 12.18] | 12.39 | [12.39, 12.39] | 0.00 | 0.01 | false | 0.53 | [0.48, 0.58] | 4.16 | [4.10, 4.23] | 0 |
| 4 | refine | 1250 | 25 | 12.16 | [12.15, 12.17] | 12.39 | [12.39, 12.39] | 0.00 | 0.01 | false | 0.53 | [0.48, 0.58] | 4.16 | [4.10, 4.22] | 0 |

## Write Mode

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

