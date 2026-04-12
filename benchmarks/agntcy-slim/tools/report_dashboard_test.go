package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildDashboardWithSmokeCapacityAndSlimRepo(t *testing.T) {
	tempDir := t.TempDir()
	smokeDir := filepath.Join(tempDir, "smoke")
	capacityDir := filepath.Join(tempDir, "capacity")
	externalDir := filepath.Join(tempDir, "external")
	if err := os.MkdirAll(smokeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(capacityDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(externalDir, 0o755); err != nil {
		t.Fatal(err)
	}

	writeFile(t, filepath.Join(smokeDir, "results.tsv"), strings.Join([]string{
		"mode\tclients\tsize\trate\trepeat\tsender_total_messages\tsender_mps\tsender_mean_latency_ms\tsender_p50_latency_ms\tsender_p90_latency_ms\tsender_p99_latency_ms\tsender_max_latency_ms\tsender_runtime_errors\tsender_duration\tsink_received_messages\tsink_errors\tsink_receive_mps\tsink_active_receive_mps\tsink_elapsed_seconds\tsink_active_receive_seconds\tsender_cpu_seconds\tsender_cpu_percent\tresponder_cpu_seconds\tresponder_cpu_percent\tnode_cpu_seconds\tnode_cpu_percent\ttotal_cpu_seconds\ttotal_cpu_percent",
		"request-reply\t1\t16\t100\t1\t100\t98\t2.10\t1.90\t2.50\t3.30\t4.10\t0\t1s\t100\t0\t100\t99\t1.0\t1.0\t0.2\t12\t0.2\t12\t0.4\t22\t0.8\t46",
		"fire-and-forget\t1\t16\t1000\t1\t1000\t940\t0\t0\t0\t0\t0\t0\t1s\t1000\t0\t930\t925\t1.0\t1.0\t0.3\t10\t0.2\t8\t0.5\t28\t1.0\t46",
		"write\t1\t16\t1000\t1\t1000\t980\t0\t0\t0\t0\t0\t0\t1s\t0\t0\t0\t0\t1.0\t0.0\t0.3\t9\t0.0\t0\t0.4\t18\t0.7\t27",
	}, "\n"))
	writeFile(t, filepath.Join(smokeDir, "suite_summary.md"), "# Smoke Summary\n\n| Column | Value |\n| --- | --- |\n| Modes | 3 |\n")
	writeFile(t, filepath.Join(smokeDir, "technical_report.md"), "# Smoke Technical\n\nRendered markdown body.")

	writeFile(t, filepath.Join(capacityDir, "results-fire-and-forget.tsv"), strings.Join([]string{
		"mode\tclients\tsize\trate\trepeat\tsender_total_messages\tsender_mps\tsender_mean_latency_ms\tsender_p50_latency_ms\tsender_p90_latency_ms\tsender_p99_latency_ms\tsender_max_latency_ms\tsender_runtime_errors\tsender_duration\tsink_received_messages\tsink_errors\tsink_receive_mps\tsink_active_receive_mps\tsink_elapsed_seconds\tsink_active_receive_seconds\tsender_cpu_seconds\tsender_cpu_percent\tresponder_cpu_seconds\tresponder_cpu_percent\tnode_cpu_seconds\tnode_cpu_percent\ttotal_cpu_seconds\ttotal_cpu_percent",
		"fire-and-forget\t1\t16384\t128000\t1\t128000\t118000\t0\t0\t0\t0\t0\t0\t5s\t127000\t0\t110000\t109500\t5.0\t5.0\t1.0\t20\t0.6\t12\t1.8\t42\t3.4\t74",
	}, "\n"))
	writeFile(t, filepath.Join(capacityDir, "capacity_sweep.md"), "# Capacity Sweep\n\n## Sink-Backed Modes\n\n| Mode | Best Offered Rate |\n| --- | --- |\n| fire-and-forget | 128000 |\n")

	writeFile(t, filepath.Join(externalDir, "benchmark-results.csv"), strings.Join([]string{
		"senders,messages_per_sender,payload_bytes,total_messages,run,send_elapsed_s,total_elapsed_s,total_received,send_mps,recv_mps",
		"1,500000,8,500000,1,0.80,0.95,500000,625000,526315",
		"4,500000,64,2000000,1,1.90,2.20,1980000,1052631,900000",
	}, "\n"))

	outputPath := filepath.Join(tempDir, "site", "index.html")
	view, err := buildDashboard(smokeDir, capacityDir, filepath.Join(externalDir, "benchmark-results.csv"), outputPath)
	if err != nil {
		t.Fatalf("buildDashboard returned error: %v", err)
	}
	html, err := renderDashboard(view)
	if err != nil {
		t.Fatalf("renderDashboard returned error: %v", err)
	}
	content := string(html)
	checks := []string{
		"CSIT SLIM Smoke Suite",
		"CSIT SLIM Capacity Sweeps",
		"SLIM Repo Data-Plane Benchmark",
		"Data-Plane Benchmark CSV Summary",
		"Smoke Summary",
		"Capacity Sweep",
		"benchmark-results.csv",
	}
	for _, check := range checks {
		if !strings.Contains(content, check) {
			t.Fatalf("expected HTML to contain %q", check)
		}
	}
}

func TestBuildDashboardWithoutInputsProducesEmptyState(t *testing.T) {
	tempDir := t.TempDir()
	outputPath := filepath.Join(tempDir, "index.html")
	view, err := buildDashboard(filepath.Join(tempDir, "missing-smoke"), filepath.Join(tempDir, "missing-capacity"), filepath.Join(tempDir, "missing.csv"), outputPath)
	if err != nil {
		t.Fatalf("buildDashboard returned error: %v", err)
	}
	html, err := renderDashboard(view)
	if err != nil {
		t.Fatalf("renderDashboard returned error: %v", err)
	}
	if !strings.Contains(string(html), "No benchmark artifacts were found") {
		t.Fatalf("expected empty-state text in generated HTML")
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
