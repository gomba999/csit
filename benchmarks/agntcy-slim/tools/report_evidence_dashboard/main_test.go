package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildContractVerifiedFromSmokeResults(t *testing.T) {
	tempDir := t.TempDir()
	smokeDir := filepath.Join(tempDir, "smoke")
	if err := os.MkdirAll(smokeDir, 0o755); err != nil {
		t.Fatal(err)
	}

	results := strings.Join([]string{
		"mode\tclients\tsize\trate\trepeat\tsender_runtime_errors\tsink_errors",
		"request-reply\t1\t16\t100\t1\t0\t0",
		"fire-and-forget\t1\t16\t1000\t1\t0\t0",
		"write\t1\t16\t1000\t1\t0\t0",
	}, "\n")
	if err := os.WriteFile(filepath.Join(smokeDir, "results.tsv"), []byte(results), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"ci-smoke-report.md", "suite_summary.md", "technical_report.md", "ci-smoke.log"} {
		if err := os.WriteFile(filepath.Join(smokeDir, name), []byte("# ok"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	outputPath := filepath.Join(tempDir, "reports", "evidence", "index.html")
	contract, err := buildContract(smokeDir, outputPath, "../index.html", "123", "https://example.test/runs/123", defaultWorkflowURL)
	if err != nil {
		t.Fatalf("buildContract returned error: %v", err)
	}

	if contract.MockData {
		t.Fatalf("expected mock_data false when results.tsv exists")
	}
	if contract.Scenario.Status != "verified" {
		t.Fatalf("scenario status = %q, want verified", contract.Scenario.Status)
	}
	if len(contract.Rows) != 3 {
		t.Fatalf("rows = %d, want 3", len(contract.Rows))
	}
	for _, row := range contract.Rows {
		if row.Status != "verified" {
			t.Fatalf("mode %s status = %q, want verified", row.ID, row.Status)
		}
		if row.EvidenceURL != "../index.html#smoke-suite" {
			t.Fatalf("mode %s evidence_url = %q", row.ID, row.EvidenceURL)
		}
		if row.ArtifactURL == "" {
			t.Fatalf("mode %s artifact_url is empty", row.ID)
		}
	}
}

func TestTemplateRenderContainsEvidenceGrid(t *testing.T) {
	contract := contractData{
		TaxonomyVersion: "v0.2-draft",
		Slice:           "c1",
		Meta: contractMeta{
			LastUpdated: "2026-06-01T12:00:00Z",
			WorkflowURL: defaultWorkflowURL,
		},
		Rows: []contractRow{
			{
				ID:            "request-reply",
				Class:         "C1",
				UseCase:       "Agent A calls B and waits for reply",
				SlimFit:       "request-reply mode on named endpoints",
				Status:        "verified",
				EvidenceURL:   "../index.html#smoke-suite",
				EvidenceLabel: "Published smoke dashboard",
				ArtifactURL:   "ci-smoke-report.md",
				ArtifactLabel: "ci-smoke-report.md",
				RerunCmd:      "task benchmarks:slim:benchmark:ci:suite-smoke",
			},
		},
		Scenario: contractScenario{
			Status:   "verified",
			RerunCmd: "task benchmarks:slim:benchmark:ci:suite-smoke",
		},
	}

	contractJSON, err := json.MarshalIndent(contract, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	view := templateView{
		ContractJSON:     string(contractJSON),
		PerfDashboardURL: "../index.html",
	}
	templatePath := filepath.Join("..", "..", "templates", "evidence_dashboard.html")
	rendered, err := renderTemplate(templatePath, view)
	if err != nil {
		t.Fatalf("renderTemplate returned error: %v", err)
	}

	html := string(rendered)
	for _, expected := range []string{
		"C1 evidence grid",
		"slim-dashboard-contract",
		"request-reply",
		"benchmark:ci:suite-smoke",
		"../index.html",
	} {
		if !strings.Contains(html, expected) {
			t.Fatalf("expected rendered HTML to contain %q", expected)
		}
	}
}
