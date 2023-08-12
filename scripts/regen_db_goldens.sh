#!/bin/bash
set -euo pipefail

ROOT="$BUILD_WORKSPACE_DIRECTORY"
cd "$ROOT"

bazel run --run_under="cd $ROOT && " //db/sqldb/golden/regen/humanreadableschema > "$ROOT/db/sqldb/golden/human_readable_schema.sql"
bazel run --run_under="cd $ROOT && " //db/sqldb/golden/regen/schemadump > "$ROOT/db/sqldb/golden/schema_dump.sql"
