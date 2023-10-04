#!/bin/bash
set -euo pipefail

ROOT="$BUILD_WORKSPACE_DIRECTORY"
cd "$ROOT"

# We keep it around because we'll need it at some point, but it can't be empty.
VALID_FLAGS=(
  "unused"
)

VALID_FLAGS_NO_ARGS=(
  "use_azure_runner"
)

# This argument-parsing monstrosity brought to you by some random GitHub Gist:
# https://gist.github.com/magnetikonline/22c1eb412daa350eeceee76c97519da8
OPTS=$(getopt \
  --longoptions "$(printf "%s:," "${VALID_FLAGS[@]}")" \
  --longoptions "$(printf "%s," "${VALID_FLAGS_NO_ARGS[@]}")" \
  --name "$(basename "$0")" \
  --options "" \
  -- "$@"
)

eval set --$OPTS
declare -a FLAGS=()
while [ ! $# -eq 0 ]
do
  case "$1" in
    --use_azure_runner)
      FLAGS+=("--use_azure_runner")
      ;;
  esac
  shift
done

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

FLAGS+=(
  "--config=${ROOT}/cmd/server/configs/local.conf"
  "--local_dsn=${LOCAL_DSN}"
)

bazel run --run_under="cd $ROOT && " //cmd/server -- "${FLAGS[@]}"
