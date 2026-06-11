#!/usr/bin/env bash
# Copyright AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0
#
# Resolves GitHub Actions run IDs that produced Pages source artifacts.
# Writes key=value lines to GITHUB_OUTPUT when set, otherwise stdout.
#
# Required env:
#   GITHUB_REPOSITORY
#   GITHUB_TOKEN
#   DEFAULT_BRANCH
#
# Optional env:
#   INPUT_INTEGRATIONS_RUN_ID
#   INPUT_BENCHMARK_RUN_ID

set -euo pipefail

DEFAULT_BRANCH="${DEFAULT_BRANCH:?DEFAULT_BRANCH is required}"
GITHUB_REPOSITORY="${GITHUB_REPOSITORY:?GITHUB_REPOSITORY is required}"
GITHUB_TOKEN="${GITHUB_TOKEN:?GITHUB_TOKEN is required}"
INPUT_INTEGRATIONS_RUN_ID="${INPUT_INTEGRATIONS_RUN_ID:-}"
INPUT_BENCHMARK_RUN_ID="${INPUT_BENCHMARK_RUN_ID:-}"

output() {
  local key="$1"
  local value="$2"
  if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
    echo "${key}=${value}" >> "$GITHUB_OUTPUT"
  else
    echo "${key}=${value}"
  fi
}

api_get() {
  curl -fsSL \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    -H "Accept: application/vnd.github+json" \
    "$1"
}

run_artifacts() {
  local run_id="$1"
  api_get "https://api.github.com/repos/${GITHUB_REPOSITORY}/actions/runs/${run_id}/artifacts?per_page=100"
}

has_matching_artifact() {
  local artifacts_json="$1"
  local pattern="$2"
  local count
  count=$(echo "$artifacts_json" | jq --arg pattern "$pattern" '[.artifacts[].name | select(test($pattern))] | length')
  [[ "$count" -gt 0 ]]
}

resolve_run_with_artifact() {
  local workflow_file="$1"
  local pattern="$2"
  local runs_json run_ids run_id artifacts_json

  runs_json=$(api_get "https://api.github.com/repos/${GITHUB_REPOSITORY}/actions/workflows/${workflow_file}/runs?branch=${DEFAULT_BRANCH}&status=completed&per_page=15")
  run_ids=$(echo "$runs_json" | jq -r '.workflow_runs[].id')

  for run_id in $run_ids; do
    artifacts_json=$(run_artifacts "$run_id")
    if has_matching_artifact "$artifacts_json" "$pattern"; then
      echo "$run_id"
      return 0
    fi
  done
  return 1
}

integrations_run_id="${INPUT_INTEGRATIONS_RUN_ID}"
if [[ -z "$integrations_run_id" ]]; then
  integrations_run_id="$(resolve_run_with_artifact test-integrations.yaml '^a2a-interop-test-result-' || true)"
fi

benchmark_run_id="${INPUT_BENCHMARK_RUN_ID}"
if [[ -z "$benchmark_run_id" ]]; then
  benchmark_run_id="$(resolve_run_with_artifact test-benchmarks-slim.yaml '^slim-benchmark-smoke-report$' || true)"
fi

has_a2a=false
has_directory_conformance=false
directory_conformance_run_id=""
if [[ -n "$integrations_run_id" ]]; then
  integrations_artifacts_json=$(run_artifacts "$integrations_run_id")
  if has_matching_artifact "$integrations_artifacts_json" '^a2a-interop-test-result-'; then
    has_a2a=true
  fi
  if has_matching_artifact "$integrations_artifacts_json" '^directory-conformance-test-result$'; then
    has_directory_conformance=true
  fi
fi

if [[ "$has_directory_conformance" != "true" && -z "${INPUT_INTEGRATIONS_RUN_ID}" ]]; then
  directory_conformance_run_id="$(resolve_run_with_artifact test-integrations.yaml '^directory-conformance-test-result$' || true)"
  if [[ -n "$directory_conformance_run_id" ]]; then
    has_directory_conformance=true
  fi
else
  directory_conformance_run_id="$integrations_run_id"
fi

has_benchmark_smoke=false
has_benchmark_capacity=false
has_benchmark_basic=false
if [[ -n "$benchmark_run_id" ]]; then
  benchmark_artifacts_json=$(run_artifacts "$benchmark_run_id")
  if has_matching_artifact "$benchmark_artifacts_json" '^slim-benchmark-smoke-report$'; then
    has_benchmark_smoke=true
  fi
  if has_matching_artifact "$benchmark_artifacts_json" '^slim-benchmark-capacity-report$'; then
    has_benchmark_capacity=true
  fi
  if has_matching_artifact "$benchmark_artifacts_json" '^slim-benchmark-basic-report$'; then
    has_benchmark_basic=true
  fi
fi

output "integrations-run-id" "$integrations_run_id"
output "directory-conformance-run-id" "$directory_conformance_run_id"
output "benchmark-run-id" "$benchmark_run_id"
output "has-a2a" "$has_a2a"
output "has-directory-conformance" "$has_directory_conformance"
output "has-benchmark-smoke" "$has_benchmark_smoke"
output "has-benchmark-capacity" "$has_benchmark_capacity"
output "has-benchmark-basic" "$has_benchmark_basic"
