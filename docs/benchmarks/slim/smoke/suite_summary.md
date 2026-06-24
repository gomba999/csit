# SLIM Benchmark Statistical Summary

**Generated:** 2026-06-24 20:13:30

**Server:** http://127.0.0.1:45605
**Destination:** agntcy/demo/echo
**Modes:** request-reply fire-and-forget write
**Clients:** 1
**Sizes:** 16
**Request-Reply Rates:** 100
**One-Way Rates:** 1000
**Write Rates:** 1000
**Duration Per Run:** 1s
**Repeats Per Case:** 25

This summary reports mean, sample variance, and confidence intervals over repeated executions of each benchmark case.

### Request-Reply Results

Request-reply prioritizes latency statistics. The configured rate is retained as load context, but the primary reported metrics are mean, p50, and p99 latency.

| Clients | Payload | Rate | Repeats | Mean Latency ms | Mean Latency Variance | Mean Latency 95% CI | P50 Latency ms | P50 Latency 95% CI | P99 Latency ms | P99 Latency 95% CI | Sender Mean CPU % | Sender CPU 95% CI | Node Mean CPU % | Node CPU 95% CI | Total Mean CPU % | Total CPU 95% CI | Total Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | 16B | 100 | 25 | 1.00 | 0.00 | [0.98, 1.03] | 1.00 | [0.97, 1.02] | 1.24 | [1.19, 1.30] | 10.82 | [10.61, 11.03] | 3.74 | [3.53, 3.95] | 22.98 | [22.59, 23.36] | 0 |

### Fire-And-Forget Results

| Clients | Payload | Rate | Repeats | Sender Mean msg/sec | Sender Variance | Sender 95% CI | Observed Node Throughput Mean msg/sec | Observed Node Throughput Variance | Observed Node Throughput 95% CI | Sender Mean CPU % | Sender CPU 95% CI | Responder Mean CPU % | Responder CPU 95% CI | Node Mean CPU % | Node CPU 95% CI | Total Mean CPU % | Total CPU 95% CI | Mean Sender Msgs | Mean Observed Msgs | Total Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | 16B | 1000 | 25 | 925.11 | 115.92 | [920.67, 929.56] | 972.24 | 63.22 | [968.95, 975.52] | 18.78 | [18.36, 19.21] | 7.07 | [6.89, 7.25] | 2.38 | [2.16, 2.60] | 28.23 | [27.65, 28.81] | 999.96 | 999.96 | 0 |

### Write Results

| Clients | Payload | Rate | Repeats | Sender Mean msg/sec | Sender Variance | Sender 95% CI | Sender Write Throughput Mean msg/sec | Sender Write Throughput Variance | Sender Write Throughput 95% CI | Sender Mean CPU % | Sender CPU 95% CI | Responder Mean CPU % | Responder CPU 95% CI | Node Mean CPU % | Node CPU 95% CI | Total Mean CPU % | Total CPU 95% CI | Mean Sender Msgs | Mean Observed Msgs | Total Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | 16B | 1000 | 25 | 928.73 | 94.00 | [924.73, 932.74] | 928.73 | 94.00 | [924.73, 932.74] | 19.04 | [18.61, 19.47] | 0.00 | [0.00, 0.00] | 2.39 | [2.17, 2.60] | 21.43 | [20.98, 21.87] | 999.84 | 999.84 | 0 |
