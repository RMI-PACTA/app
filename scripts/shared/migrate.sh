#!/bin/bash
set -euo pipefail

function migrate_db {
  echo "Building migrator..."
  bazel build //cmd/tools/migratesqldb

  DSN="$1"
  SUBCOMMAND="${2:-apply}"

  bazel run //cmd/tools/migratesqldb -- --dsn="$DSN" "$SUBCOMMAND"
}
