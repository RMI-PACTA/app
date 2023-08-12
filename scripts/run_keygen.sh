#!/bin/bash
set -euo pipefail

ROOT="$BUILD_WORKSPACE_DIRECTORY"
cd "$ROOT"

bazel run --run_under="cd $ROOT && " //cmd/tools/keygen
