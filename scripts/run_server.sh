#!/bin/bash
set -euo pipefail

ROOT="$BUILD_WORKSPACE_DIRECTORY"
cd "$ROOT"


PG_LOCAL_DATA="$ROOT/.postgres-data"
PG_PROPS_FILE="$PG_LOCAL_DATA/db.properties"

source "./scripts/shared/kv.sh"
init_kv "$PG_PROPS_FILE"

# Build up the DSN manually, for readability.
LOCAL_DSN="user=$(get_val "PG_USER")"
LOCAL_DSN+=" password=$(get_val "PG_PASSWORD")"
LOCAL_DSN+=" host=$(get_val "SOCKET_DIR")"
LOCAL_DSN+=" dbname=$(get_val "PG_DB_NAME")"
LOCAL_DSN+=" sslmode=disable"

declare -a FLAGS=(
  "--config=${ROOT}/cmd/server/configs/local.conf"
  "--local_dsn=${LOCAL_DSN}"
)

bazel run --run_under="cd $ROOT && " //cmd/server -- "${FLAGS[@]}"
