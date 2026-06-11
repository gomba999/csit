#!/usr/bin/env bash
# Copyright AGNTCY Contributors (https://github.com/agntcy)
# SPDX-License-Identifier: Apache-2.0
#
# Downloads CI report artifacts and assembles the GitHub Pages site tree.
#
# Required env:
#   GITHUB_REPOSITORY
#   GITHUB_TOKEN
#
# Optional env (from resolve-pages-source-runs.sh):
#   INTEGRATIONS_RUN_ID
#   DIRECTORY_CONFORMANCE_RUN_ID
#   BENCHMARK_RUN_ID
#   HAS_A2A, HAS_DIRECTORY_CONFORMANCE
#   HAS_BENCHMARK_SMOKE, HAS_BENCHMARK_CAPACITY, HAS_BENCHMARK_BASIC
#   SITE_DIR (default: site)
#   SLIM_REPO_BENCHMARK_RUN_ID
#   SLIM_REPO_ACTIONS_READ_TOKEN

set -euo pipefail

GITHUB_REPOSITORY="${GITHUB_REPOSITORY:?GITHUB_REPOSITORY is required}"
GITHUB_TOKEN="${GITHUB_TOKEN:?GITHUB_TOKEN is required}"
SITE_DIR="${SITE_DIR:-site}"
REPO_ROOT="${REPO_ROOT:-$(cd "$(dirname "$0")/../.." && pwd)}"

INTEGRATIONS_RUN_ID="${INTEGRATIONS_RUN_ID:-}"
DIRECTORY_CONFORMANCE_RUN_ID="${DIRECTORY_CONFORMANCE_RUN_ID:-}"
BENCHMARK_RUN_ID="${BENCHMARK_RUN_ID:-}"
HAS_A2A="${HAS_A2A:-false}"
HAS_DIRECTORY_CONFORMANCE="${HAS_DIRECTORY_CONFORMANCE:-false}"
HAS_BENCHMARK_SMOKE="${HAS_BENCHMARK_SMOKE:-false}"
HAS_BENCHMARK_CAPACITY="${HAS_BENCHMARK_CAPACITY:-false}"
HAS_BENCHMARK_BASIC="${HAS_BENCHMARK_BASIC:-false}"
SLIM_REPO_BENCHMARK_RUN_ID="${SLIM_REPO_BENCHMARK_RUN_ID:-}"
SLIM_REPO_ACTIONS_READ_TOKEN="${SLIM_REPO_ACTIONS_READ_TOKEN:-}"

api_get() {
  local token="$1"
  local url="$2"
  curl -fsSL \
    -H "Authorization: Bearer ${token}" \
    -H "Accept: application/vnd.github+json" \
    "$url"
}

download_named_artifact() {
  local run_id="$1"
  local artifact_name="$2"
  local dest_dir="$3"
  local artifacts_json download_url archive_path

  artifacts_json=$(api_get "$GITHUB_TOKEN" \
    "https://api.github.com/repos/${GITHUB_REPOSITORY}/actions/runs/${run_id}/artifacts?per_page=100")

  download_url=$(echo "$artifacts_json" | jq -r --arg name "$artifact_name" \
    '.artifacts[] | select(.name == $name) | .archive_download_url' | head -n1)

  if [[ -z "$download_url" || "$download_url" == "null" ]]; then
    echo "Artifact ${artifact_name} not found in run ${run_id}" >&2
    return 1
  fi

  mkdir -p "$dest_dir"
  archive_path=$(mktemp "${RUNNER_TEMP:-/tmp}/pages-artifact-XXXXXX.zip")
  trap 'rm -f "$archive_path"' RETURN

  curl -fsSL -L \
    -H "Authorization: Bearer ${GITHUB_TOKEN}" \
    -H "Accept: application/vnd.github+json" \
    "$download_url" \
    -o "$archive_path"

  unzip -qo "$archive_path" -d "$dest_dir"
}

download_pattern_artifacts() {
  local run_id="$1"
  local pattern="$2"
  local dest_dir="$3"
  local artifacts_json names name download_url archive_path

  artifacts_json=$(api_get "$GITHUB_TOKEN" \
    "https://api.github.com/repos/${GITHUB_REPOSITORY}/actions/runs/${run_id}/artifacts?per_page=100")

  names=$(echo "$artifacts_json" | jq -r --arg pattern "$pattern" \
    '.artifacts[] | select(.name | test($pattern)) | .name')

  mkdir -p "$dest_dir"
  for name in $names; do
    download_url=$(echo "$artifacts_json" | jq -r --arg name "$name" \
      '.artifacts[] | select(.name == $name) | .archive_download_url')
    archive_path=$(mktemp "${RUNNER_TEMP:-/tmp}/pages-artifact-XXXXXX.zip")
    curl -fsSL -L \
      -H "Authorization: Bearer ${GITHUB_TOKEN}" \
      -H "Accept: application/vnd.github+json" \
      "$download_url" \
      -o "$archive_path"
    unzip -qo "$archive_path" -d "$dest_dir"
    rm -f "$archive_path"
  done
}

rm -rf "$SITE_DIR"
mkdir -p "$SITE_DIR"

if [[ "$HAS_A2A" == "true" && -n "$INTEGRATIONS_RUN_ID" ]]; then
  download_pattern_artifacts "$INTEGRATIONS_RUN_ID" '^a2a-interop-test-result-' "${SITE_DIR}/a2a"
  compgen -G "${SITE_DIR}/a2a/*.json" > /dev/null
  (
    cd "${REPO_ROOT}/integrations"
    go run ./agntcy-a2a/tools/report_dashboard.go \
      --reports-dir "../${SITE_DIR}/a2a" \
      --output "../${SITE_DIR}/a2a/index.html"
  )
fi

if [[ "$HAS_DIRECTORY_CONFORMANCE" == "true" ]]; then
  local_run_id="${DIRECTORY_CONFORMANCE_RUN_ID:-$INTEGRATIONS_RUN_ID}"
  tmp_dir="${SITE_DIR}/directory-conformance-tmp"
  download_named_artifact "$local_run_id" "directory-conformance-test-result" "$tmp_dir"
  mkdir -p "${SITE_DIR}/directory"
  if [[ -f "${tmp_dir}/summary.html" ]]; then
    cp "${tmp_dir}/summary.html" "${SITE_DIR}/directory/index.html"
  else
    echo "summary.html not found in directory conformance artifact" >&2
    exit 1
  fi
  rm -rf "$tmp_dir"
fi

if [[ "$HAS_BENCHMARK_SMOKE" == "true" && -n "$BENCHMARK_RUN_ID" ]]; then
  download_named_artifact "$BENCHMARK_RUN_ID" "slim-benchmark-smoke-report" "${SITE_DIR}/benchmarks/slim/smoke"
fi

if [[ "$HAS_BENCHMARK_CAPACITY" == "true" && -n "$BENCHMARK_RUN_ID" ]]; then
  download_named_artifact "$BENCHMARK_RUN_ID" "slim-benchmark-capacity-report" "${SITE_DIR}/benchmarks/slim/capacity"
fi

if [[ "$HAS_BENCHMARK_BASIC" == "true" && -n "$BENCHMARK_RUN_ID" ]]; then
  download_named_artifact "$BENCHMARK_RUN_ID" "slim-benchmark-basic-report" "${SITE_DIR}/benchmarks/slim/basic"
fi

has_slim_external=false
if [[ -n "$SLIM_REPO_BENCHMARK_RUN_ID" ]]; then
  if [[ -z "$SLIM_REPO_ACTIONS_READ_TOKEN" ]]; then
    echo "SLIM_REPO_ACTIONS_READ_TOKEN is not configured; skipping external slim benchmark import."
  else
    artifacts_json=$(api_get "$SLIM_REPO_ACTIONS_READ_TOKEN" \
      "https://api.github.com/repos/agntcy/slim/actions/runs/${SLIM_REPO_BENCHMARK_RUN_ID}/artifacts?per_page=100")
    download_url=$(echo "$artifacts_json" | jq -r \
      '.artifacts[] | select(.name | startswith("benchmark-results-")) | .archive_download_url' | head -n1)
    if [[ -n "$download_url" && "$download_url" != "null" ]]; then
      mkdir -p "${SITE_DIR}/benchmarks/slim/external"
      archive_path=$(mktemp "${RUNNER_TEMP:-/tmp}/slim-benchmark-XXXXXX.zip")
      curl -fsSL -L \
        -H "Authorization: Bearer ${SLIM_REPO_ACTIONS_READ_TOKEN}" \
        -H "Accept: application/vnd.github+json" \
        "$download_url" \
        -o "$archive_path"
      unzip -qo "$archive_path" -d "${SITE_DIR}/benchmarks/slim/external"
      rm -f "$archive_path"
      csv_path=$(find "${SITE_DIR}/benchmarks/slim/external" -type f -name 'benchmark-results*.csv' | head -n1 || true)
      if [[ -n "$csv_path" && "$(basename "$csv_path")" != "benchmark-results.csv" ]]; then
        cp "$csv_path" "${SITE_DIR}/benchmarks/slim/external/benchmark-results.csv"
      fi
      if [[ -n "$csv_path" ]]; then
        has_slim_external=true
      fi
    else
      echo "No retained benchmark-results artifact was found for agntcy/slim run ${SLIM_REPO_BENCHMARK_RUN_ID}."
    fi
  fi
fi

if [[ "$HAS_BENCHMARK_SMOKE" == "true" || "$HAS_BENCHMARK_CAPACITY" == "true" || "$HAS_BENCHMARK_BASIC" == "true" || "$has_slim_external" == "true" ]]; then
  mkdir -p "${SITE_DIR}/benchmarks/slim/smoke" "${SITE_DIR}/benchmarks/slim/capacity"
  args=(go run ./agntcy-slim/tools/report_dashboard.go
    --smoke-dir "../${SITE_DIR}/benchmarks/slim/smoke"
    --capacity-dir "../${SITE_DIR}/benchmarks/slim/capacity"
    --output "../${SITE_DIR}/benchmarks/slim/index.html")
  if [[ -f "${SITE_DIR}/benchmarks/slim/basic/basic-benchmark-results.csv" ]]; then
    args+=(--basic-csv "../${SITE_DIR}/benchmarks/slim/basic/basic-benchmark-results.csv")
  fi
  if [[ -f "${SITE_DIR}/benchmarks/slim/external/benchmark-results.csv" ]]; then
    args+=(--slim-csv "../${SITE_DIR}/benchmarks/slim/external/benchmark-results.csv")
  fi
  (
    cd "${REPO_ROOT}/benchmarks"
    "${args[@]}"
  )
fi

HAS_A2A="$([ -f "${SITE_DIR}/a2a/index.html" ] && echo true || echo false)"
HAS_BENCHMARKS="$([ -f "${SITE_DIR}/benchmarks/slim/index.html" ] && echo true || echo false)"
HAS_DIRECTORY="$([ -f "${SITE_DIR}/directory/index.html" ] && echo true || echo false)"
OUTPUT="${SITE_DIR}/index.html" \
  HAS_A2A="$HAS_A2A" \
  HAS_BENCHMARKS="$HAS_BENCHMARKS" \
  HAS_DIRECTORY="$HAS_DIRECTORY" \
  INTEGRATIONS_RUN_ID="$INTEGRATIONS_RUN_ID" \
  BENCHMARK_RUN_ID="$BENCHMARK_RUN_ID" \
  bash "${REPO_ROOT}/.github/scripts/build-site-landing-page.sh"

echo "assembled ${SITE_DIR}/"
