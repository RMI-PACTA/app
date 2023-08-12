#!/bin/bash
set -euo pipefail

ROOT="$BUILD_WORKSPACE_DIRECTORY"
cd "$ROOT"

PG_LOCAL_DATA="$ROOT/.postgres-data"
PG_PROPS_FILE="$PG_LOCAL_DATA/db.properties"

source "./scripts/shared/migrate.sh"
source "./scripts/shared/kv.sh"
init_kv "$PG_PROPS_FILE"

# Note: This is fixed by the container, changing it here will not work. For
# more info, see https://hub.docker.com/_/postgres/
PG_USER="postgres"
PG_PASSWORD="$(get_val "PG_PASSWORD")"
SOCKET_DIR="$(get_val "SOCKET_DIR")"
PG_DB_NAME="$(get_val "PG_DB_NAME")"

DSN="user=${PG_USER}"
DSN+=" password=${PG_PASSWORD}"
DSN+=" host=${SOCKET_DIR}"
DSN+=" dbname=${PG_DB_NAME}"
DSN+=" sslmode=disable"

migrate_db "$DSN" rollback
