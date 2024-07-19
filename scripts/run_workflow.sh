#!/bin/bash

# run_workflow.sh is a script for testing out the analysis code from
# https://github.com/RMI-PACTA/workflow.pacta.webapp

set -euo pipefail

ROOT="$BUILD_WORKSPACE_DIRECTORY"
cd "$ROOT"

PORTFOLIO_FOLDER=""
BENCHMARKS_FOLDER=""
YEAR="2022"
case $YEAR in
  "2022")
    PORTFOLIO_FOLDER="2022Q4_20240426T113151Z"
    BENCHMARKS_FOLDER="2022Q4_20240529T002407Z"
    ;;
  "2023")
    PORTFOLIO_FOLDER="2023Q4_20240424T120055Z"
    BENCHMARKS_FOLDER="2023Q4_20240529T002355Z"
    ;;
  *)
    echo "unexpected year $YEAR"
    exit 1
    ;;
esac

IS_PODMAN=false
if docker --version | grep -q 'podman'; then
  IS_PODMAN=true
  echo "Running against podman"
fi

declare -a DOCKER_FLAGS=(
  # Read-only mounts and corresponding env vars
  "-v" "$ROOT/workflow-data/pacta-data/$PORTFOLIO_FOLDER:/mnt/pacta-data:ro"
  "-v" "$ROOT/workflow-data/benchmarks/$BENCHMARKS_FOLDER:/mnt/benchmarks:ro"
  "-v" "$ROOT/workflow-data/portfolios:/mnt/portfolios:ro"
  "-v" "$ROOT/workflow-data/real-estate:/mnt/real-estate:ro"
  "-v" "$ROOT/workflow-data/score-card:/mnt/score-card:ro"
  "-v" "$ROOT/workflow-data/survey:/mnt/survey:ro"
  "-e" "BENCHMARKS_DIR=/mnt/benchmarks"
  "-e" "PACTA_DATA_DIR=/mnt/pacta-data"
  "-e" "PORTFOLIO_DIR=/mnt/portfolios"
  "-e" "REAL_ESTATE_DIR=/mnt/real-estate"
  "-e" "SCORE_CARD_DIR=/mnt/score-card"
  "-e" "SURVEY_DIR=/mnt/survey"

  # Write mounts and corresponding env vars
  "-v" "$ROOT/workflow-data/analysis-output:/mnt/analysis-output"
  "-v" "$ROOT/workflow-data/report-output:/mnt/report-output"
  "-v" "$ROOT/workflow-data/summary-output:/mnt/summary-output"
  "-e" "ANALYSIS_OUTPUT_DIR=/mnt/analysis-output"
  "-e" "REPORT_OUTPUT_DIR=/mnt/report-output"
  "-e" "SUMMARY_OUTPUT_DIR=/mnt/summary-output"

  # Misc
  "-e" "LOG_LEVEL=DEBUG"
)

# TODO: Unclear if this will work in Docker as is. For 'normal' root-running
# Docker daemons, it should probably just create files/directories on the host
# owned by 1000:1000, which is fine.
if [ "$IS_PODMAN" = true ]; then
  DOCKER_FLAGS+=("--userns" "keep-id:uid=1000,gid=1000")
fi


JSON_INPUT="{
  \"portfolio\": {
    \"files\": \"default_portfolio.csv\",
    \"holdingsDate\": \"2023-12-31\",
    \"name\": \"FooPortfolio\"
  },
  \"inherit\": \"GENERAL_2023Q4\"
}"

docker run --rm -it \
  "${DOCKER_FLAGS[@]}" \
  ghcr.io/rmi-pacta/workflow.pacta.webapp:nightly  "$JSON_INPUT"
