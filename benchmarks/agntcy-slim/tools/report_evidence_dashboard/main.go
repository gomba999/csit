package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultWorkflowURL = "https://github.com/agntcy/csit/actions/workflows/test-benchmarks-slim.yaml"
)

var modeOrder = []string{"request-reply", "fire-and-forget", "write"}

type modeSpec struct {
	UseCase      string
	SlimFit      string
	ArtifactFile string
	ArtifactName string
}

var modeSpecs = map[string]modeSpec{
	"request-reply": {
		UseCase:      "Agent A calls B and waits for reply",
		SlimFit:      "request-reply mode on named endpoints",
		ArtifactFile: "ci-smoke-report.md",
		ArtifactName: "ci-smoke-report.md",
	},
	"fire-and-forget": {
		UseCase:      "Agent fires event; consumer handles async",
		SlimFit:      "fire-and-forget one-way delivery",
		ArtifactFile: "technical_report.md",
		ArtifactName: "technical_report.md",
	},
	"write": {
		UseCase:      "Publish into mesh without paired responder",
		SlimFit:      "write ingress / publish pattern",
		ArtifactFile: "suite_summary.md",
		ArtifactName: "suite_summary.md",
	},
}

type modeStatus struct {
	Seen   bool
	Failed bool
}

type contractData struct {
	TaxonomyVersion string           `json:"taxonomy_version"`
	Slice           string           `json:"slice"`
	MockData        bool             `json:"mock_data"`
	Meta            contractMeta     `json:"meta"`
	Rows            []contractRow    `json:"rows"`
	Scenario        contractScenario `json:"scenario"`
}

type contractMeta struct {
	LastUpdated string `json:"last_updated"`
	Discussion  string `json:"discussion_url"`
	EpicURL     string `json:"epic_url"`
	WorkflowURL string `json:"workflow_url"`
	LastRunID   string `json:"last_run_id,omitempty"`
	LastRunURL  string `json:"last_run_url,omitempty"`
}

type contractRow struct {
	ID            string `json:"id"`
	Class         string `json:"class"`
	UseCase       string `json:"use_case"`
	SlimFit       string `json:"slim_fit"`
	Status        string `json:"status"`
	LastRun       string `json:"last_run"`
	EvidenceURL   string `json:"evidence_url"`
	EvidenceLabel string `json:"evidence_label"`
	ArtifactURL   string `json:"artifact_url"`
	ArtifactLabel string `json:"artifact_label"`
	RerunCmd      string `json:"rerun_cmd"`
}

type scenarioArtifact struct {
	File string `json:"file"`
	Path string `json:"path"`
}

type contractScenario struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Classes     []string           `json:"classes_served"`
	Status      string             `json:"status"`
	Description string             `json:"description"`
	RerunCmd    string             `json:"rerun_cmd"`
	Artifacts   []scenarioArtifact `json:"artifacts"`
}

type templateView struct {
	ContractJSON     string
	PerfDashboardURL string
}

func main() {
	smokeDir := flag.String("smoke-dir", "", "directory containing smoke artifacts (results.tsv, markdown reports, logs)")
	outputPath := flag.String("output", "./reports/evidence/index.html", "path to the generated evidence dashboard HTML")
	perfURL := flag.String("perf-url", "../index.html", "relative URL to the performance dashboard")
	templatePath := flag.String("template", "./templates/evidence_dashboard.html", "path to the evidence dashboard HTML template")
	runID := flag.String("run-id", "", "optional workflow run id displayed in metadata")
	runURL := flag.String("run-url", "", "optional workflow run url displayed in metadata")
	workflowURL := flag.String("workflow-url", defaultWorkflowURL, "workflow URL displayed in metadata")
	flag.Parse()

	if strings.TrimSpace(*smokeDir) == "" {
		fmt.Fprintln(os.Stderr, "build evidence dashboard: --smoke-dir is required")
		os.Exit(1)
	}

	contract, err := buildContract(*smokeDir, *outputPath, *perfURL, *runID, *runURL, *workflowURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build evidence dashboard: %v\n", err)
		os.Exit(1)
	}

	contractJSON, err := json.MarshalIndent(contract, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal contract: %v\n", err)
		os.Exit(1)
	}

	view := templateView{
		ContractJSON:     string(contractJSON),
		PerfDashboardURL: *perfURL,
	}

	htmlOutput, err := renderTemplate(*templatePath, view)
	if err != nil {
		fmt.Fprintf(os.Stderr, "render evidence dashboard: %v\n", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(filepath.Dir(*outputPath), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "create output directory: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*outputPath, htmlOutput, 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write evidence dashboard: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("wrote %s\n", *outputPath)
}

func buildContract(smokeDir string, outputPath string, perfURL string, runID string, runURL string, workflowURL string) (contractData, error) {
	statusByMode, smokeUpdatedAt, err := readModeStatus(filepath.Join(smokeDir, "results.tsv"))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return contractData{}, err
	}

	lastUpdated := time.Now().UTC()
	if !smokeUpdatedAt.IsZero() {
		lastUpdated = smokeUpdatedAt.UTC()
	}

	contract := contractData{
		TaxonomyVersion: "v0.2-draft",
		Slice:           "c1",
		MockData:        len(statusByMode) == 0,
		Meta: contractMeta{
			LastUpdated: lastUpdated.Format(time.RFC3339),
			Discussion:  "https://github.com/agntcy/csit/discussions/195",
			EpicURL:     "https://github.com/agntcy/csit/issues/200",
			WorkflowURL: workflowURL,
			LastRunID:   runID,
			LastRunURL:  runURL,
		},
		Scenario: contractScenario{
			ID:          "local-messaging-smoke",
			Title:       "Local messaging smoke",
			Classes:     []string{"C1"},
			Status:      "unknown",
			Description: "Single-node CSIT smoke suite exercising request-reply, fire-and-forget, and write workloads.",
			RerunCmd: strings.Join([]string{
				"task benchmarks:slim:deps:slimctl-download SLIMCTL_PATH=\"$HOME/.local/bin/slimctl\"",
				"export PATH=\"$HOME/.local/bin:$PATH\"",
				"task benchmarks:slim:benchmark:ci:suite-smoke",
			}, "\n"),
		},
	}

	allVerified := true
	for _, mode := range modeOrder {
		spec := modeSpecs[mode]
		rowStatus := "unknown"
		if modeState, ok := statusByMode[mode]; ok {
			if modeState.Failed {
				rowStatus = "failed"
				allVerified = false
			} else if modeState.Seen {
				rowStatus = "verified"
			}
		} else {
			allVerified = false
		}

		row := contractRow{
			ID:            mode,
			Class:         "C1",
			UseCase:       spec.UseCase,
			SlimFit:       spec.SlimFit,
			Status:        rowStatus,
			LastRun:       runID,
			EvidenceURL:   perfURL + "#smoke-suite",
			EvidenceLabel: "Published smoke dashboard",
			ArtifactURL:   relFileURL(outputPath, filepath.Join(smokeDir, spec.ArtifactFile)),
			ArtifactLabel: spec.ArtifactName,
			RerunCmd:      "task benchmarks:slim:benchmark:ci:suite-smoke",
		}
		contract.Rows = append(contract.Rows, row)
	}
	if allVerified {
		contract.Scenario.Status = "verified"
	}
	contract.Scenario.Artifacts = collectScenarioArtifacts(smokeDir, outputPath)

	return contract, nil
}

func readModeStatus(resultsPath string) (map[string]modeStatus, time.Time, error) {
	info, err := os.Stat(resultsPath)
	if err != nil {
		return nil, time.Time{}, err
	}

	file, err := os.Open(resultsPath)
	if err != nil {
		return nil, time.Time{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		return nil, time.Time{}, err
	}

	idxMode := indexOf(header, "mode")
	idxSenderErrors := indexOf(header, "sender_runtime_errors")
	idxSinkErrors := indexOf(header, "sink_errors")
	if idxMode < 0 || idxSenderErrors < 0 || idxSinkErrors < 0 {
		return nil, time.Time{}, fmt.Errorf("results.tsv missing required columns")
	}

	statusByMode := make(map[string]modeStatus)
	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if errors.Is(err, csv.ErrFieldCount) {
				continue
			}
			return nil, time.Time{}, err
		}
		if len(record) <= idxSinkErrors {
			continue
		}
		mode := strings.TrimSpace(record[idxMode])
		if mode == "" {
			continue
		}
		current := statusByMode[mode]
		current.Seen = true
		if strings.TrimSpace(record[idxSenderErrors]) != "0" || strings.TrimSpace(record[idxSinkErrors]) != "0" {
			current.Failed = true
		}
		statusByMode[mode] = current
	}

	return statusByMode, info.ModTime(), nil
}

func indexOf(cols []string, key string) int {
	for i, col := range cols {
		if strings.TrimSpace(col) == key {
			return i
		}
	}
	return -1
}

func collectScenarioArtifacts(smokeDir string, outputPath string) []scenarioArtifact {
	names := []string{"index.html", "ci-smoke-report.md", "suite_summary.md", "technical_report.md", "ci-smoke.log", "results.tsv"}
	artifacts := make([]scenarioArtifact, 0, len(names))
	for _, name := range names {
		full := filepath.Join(smokeDir, name)
		if _, err := os.Stat(full); err == nil {
			artifacts = append(artifacts, scenarioArtifact{
				File: name,
				Path: relFileURL(outputPath, full),
			})
		}
	}
	return artifacts
}

func relFileURL(outputPath string, targetPath string) string {
	rel, err := filepath.Rel(filepath.Dir(outputPath), targetPath)
	if err != nil {
		return targetPath
	}
	return filepath.ToSlash(rel)
}

func renderTemplate(templatePath string, view templateView) ([]byte, error) {
	tmpl, err := template.New(filepath.Base(templatePath)).ParseFiles(templatePath)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, view); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
