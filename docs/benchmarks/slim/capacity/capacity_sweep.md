# SLIM Adaptive Capacity Sweep Report

**Generated:** 2026-06-24 17:11:29

This CI report combines the sink-backed capacity sweeps and the write capacity sweep into one markdown artifact.

## Sink-Backed Modes

#### Fire-And-Forget Clients=1 Payload=16384B

Best offered aggregate rate: `128000` msg/sec
Estimated capacity offered-rate interval: `[128000, 130125]` msg/sec
Best observed node throughput: `7386.03` msg/sec with 95% CI [7343.40, 7428.65]
Best sender-completed throughput: `7308.84` msg/sec with 95% CI [7262.36, 7355.33]
Best node CPU: `41.47` % with 95% CI [41.11, 41.83]
Best total CPU: `247.44` % with 95% CI [245.91, 248.98]
Stop reason: refinement narrowed the estimated capacity to offered rates 128000 through 130125

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 128000 | 25 | 7308.84 | [7262.36, 7355.33] | 7386.03 | [7343.40, 7428.65] | 10661.86 | 0.00 | true | 41.47 | [41.11, 41.83] | 247.44 | [245.91, 248.98] | 0 |
| 2 | coarse | 144000 | 25 | 7391.44 | [7368.02, 7414.87] | 7469.64 | [7450.80, 7488.49] | 2083.28 | 1.13 | false | 41.95 | [41.84, 42.05] | 249.43 | [249.02, 249.83] | 0 |
| 3 | coarse | 162000 | 25 | 7378.01 | [7362.86, 7393.16] | 7465.60 | [7449.99, 7481.21] | 1430.77 | 1.08 | false | 41.90 | [41.80, 42.00] | 248.96 | [248.46, 249.46] | 0 |
| 4 | refine | 145000 | 25 | 7362.34 | [7346.47, 7378.22] | 7442.39 | [7428.08, 7456.71] | 1202.51 | 0.76 | false | 41.74 | [41.61, 41.87] | 248.21 | [247.56, 248.85] | 0 |
| 5 | refine | 136500 | 25 | 7385.30 | [7367.31, 7403.30] | 7472.54 | [7458.80, 7486.27] | 1106.94 | 1.17 | false | 41.88 | [41.79, 41.98] | 248.86 | [248.40, 249.33] | 0 |
| 6 | refine | 132250 | 25 | 7379.39 | [7359.48, 7399.30] | 7476.20 | [7461.50, 7490.89] | 1267.91 | 1.22 | false | 41.83 | [41.70, 41.96] | 248.66 | [248.11, 249.22] | 0 |
| 7 | refine | 130125 | 25 | 7393.22 | [7376.24, 7410.21] | 7471.29 | [7455.33, 7487.25] | 1494.41 | 1.15 | false | 41.90 | [41.78, 42.01] | 248.97 | [248.46, 249.48] | 0 |

#### Request-Reply Clients=1 Payload=16384B

Best offered aggregate rate: `1000` msg/sec
Estimated capacity offered-rate interval: `[1000, 1250]` msg/sec
Best observed node throughput: `12.38` msg/sec with 95% CI [12.38, 12.38]
Best sender-completed throughput: `12.15` msg/sec with 95% CI [12.14, 12.17]
Best node CPU: `0.53` % with 95% CI [0.48, 0.59]
Best total CPU: `3.53` % with 95% CI [3.47, 3.59]
Stop reason: refinement narrowed the estimated capacity to offered rates 1000 through 1250

| Step | Phase | Offered Aggregate Rate | Repeats | Sender Mean msg/sec | Sender 95% CI | Observed Node Throughput | Observed Node Throughput 95% CI | Observed Variance | Observed Gain % | Improved | Node CPU % | Node CPU 95% CI | Total CPU % | Total CPU 95% CI | Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | coarse | 1000 | 25 | 12.15 | [12.14, 12.17] | 12.38 | [12.38, 12.38] | 0.00 | 0.00 | true | 0.53 | [0.48, 0.59] | 3.53 | [3.47, 3.59] | 0 |
| 2 | coarse | 2000 | 25 | 12.15 | [12.14, 12.16] | 12.38 | [12.38, 12.38] | 0.00 | 0.01 | false | 0.59 | [0.53, 0.64] | 3.65 | [3.58, 3.71] | 0 |
| 3 | refine | 1500 | 25 | 12.15 | [12.14, 12.17] | 12.38 | [12.38, 12.38] | 0.00 | 0.00 | false | 0.59 | [0.53, 0.64] | 3.65 | [3.58, 3.72] | 0 |
| 4 | refine | 1250 | 25 | 12.16 | [12.15, 12.17] | 12.38 | [12.37, 12.38] | 0.00 | -0.01 | false | 0.59 | [0.54, 0.63] | 3.66 | [3.60, 3.72] | 0 |

## Write Mode

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

