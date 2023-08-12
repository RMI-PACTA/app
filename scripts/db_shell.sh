#!/bin/bash
set -euo pipefail

ROOT="$BUILD_WORKSPACE_DIRECTORY"
cd "$ROOT"

PG_LOCAL_DATA="$ROOT/.postgres-data"
PG_PROPS_FILE="$PG_LOCAL_DATA/db.properties"

source "./scripts/shared/kv.sh"
init_kv "$PG_PROPS_FILE"

# Build up the DSN manually, for readability.
SOCKET_DIR="$(get_val "SOCKET_DIR")"

LOCAL_DSN="postgresql://"
LOCAL_DSN+="$(get_val "PG_USER"):$(get_val "PG_PASSWORD")"
LOCAL_DSN+="@/$(get_val "PG_DB_NAME")"
LOCAL_DSN+="?host=${SOCKET_DIR}"
LOCAL_DSN+="&sslmode=disable"

if [ -x "$(command -v psql)" ]; then
  # Use psql directly if it is installed
  psql "$LOCAL_DSN"
elif [ -x "$(command -v docker)" ]; then
  # If psql isn't installed, try running it within Docker. We use this image in
  # //scripts/run_db.sh as well, so its not a bad fit.
  docker run \
    --rm -it \
    --net=host \
    --volume "${SOCKET_DIR}:${SOCKET_DIR}" \
    postgres:14.9 psql "$LOCAL_DSN"
else
  echo 'Error: neither psql or docker is installed.' >&2
  exit 1
fi
