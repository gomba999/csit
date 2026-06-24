# SLIM Benchmark Statistical Summary

**Generated:** 2026-06-24 20:14:53

**Server:** http://127.0.0.1:44081
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
| 1 | 16B | 100 | 25 | 0.74 | 0.00 | [0.72, 0.76] | 0.69 | [0.68, 0.71] | 1.69 | [1.20, 2.19] | 7.23 | [7.06, 7.40] | 2.37 | [2.17, 2.57] | 14.42 | [14.16, 14.69] | 0 |

### Fire-And-Forget Results

| Clients | Payload | Rate | Repeats | Sender Mean msg/sec | Sender Variance | Sender 95% CI | Observed Node Throughput Mean msg/sec | Observed Node Throughput Variance | Observed Node Throughput 95% CI | Sender Mean CPU % | Sender CPU 95% CI | Responder Mean CPU % | Responder CPU 95% CI | Node Mean CPU % | Node CPU 95% CI | Total Mean CPU % | Total CPU 95% CI | Mean Sender Msgs | Mean Observed Msgs | Total Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | 16B | 1000 | 25 | 926.13 | 218.92 | [920.03, 932.24] | 972.91 | 122.05 | [968.35, 977.47] | 13.97 | [13.60, 14.34] | 5.11 | [4.90, 5.31] | 2.06 | [1.83, 2.28] | 21.14 | [20.59, 21.68] | 999.80 | 999.80 | 0 |

### Write Results

| Clients | Payload | Rate | Repeats | Sender Mean msg/sec | Sender Variance | Sender 95% CI | Sender Write Throughput Mean msg/sec | Sender Write Throughput Variance | Sender Write Throughput 95% CI | Sender Mean CPU % | Sender CPU 95% CI | Responder Mean CPU % | Responder CPU 95% CI | Node Mean CPU % | Node CPU 95% CI | Total Mean CPU % | Total CPU 95% CI | Mean Sender Msgs | Mean Observed Msgs | Total Errors |
| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |
| 1 | 16B | 1000 | 25 | 927.99 | 232.82 | [921.69, 934.29] | 927.99 | 232.82 | [921.69, 934.29] | 13.52 | [13.16, 13.88] | 0.00 | [0.00, 0.00] | 1.95 | [1.72, 2.18] | 15.47 | [15.09, 15.86] | 999.88 | 999.88 | 0 |
