#!/bin/bash
set -euo pipefail

ROOT="$BUILD_WORKSPACE_DIRECTORY"
cd "$ROOT"

SKIP_MIGRATE=false
while [ ! $# -eq 0 ]
do
	case "$1" in
		--skip-migrate)
			SKIP_MIGRATE=true
			;;
	esac
	shift
done

# Note: This is fixed by the container, changing it here will not work. For
# more info, see https://hub.docker.com/_/postgres/
PG_USER="postgres"
PG_DATA_PATH="/var/lib/postgresql/data"
PG_LOCAL_DATA="$ROOT/.postgres-data"
PG_PROPS_FILE="$PG_LOCAL_DATA/db.properties"
PG_LOCAL_DATA_PATH="$PG_LOCAL_DATA/data"

mkdir -p "$PG_LOCAL_DATA_PATH"

source "./scripts/shared/migrate.sh"
source "./scripts/shared/kv.sh"
init_kv "$PG_PROPS_FILE"
set_kv "PG_USER" "$PG_USER"

SOCKET_DIR=""
CONTAINER_ID=""

PG_PASSWORD="$(get_val "PG_PASSWORD")"
PG_DB_NAME="$(get_val "PG_DB_NAME")"
if [ -z "${PG_PASSWORD}" ]; then
  echo "First run! Generating a random password..."
  PG_PASSWORD="$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c20)"
  set_kv "PG_PASSWORD" "$PG_PASSWORD"
fi
if [ -z "${PG_DB_NAME}" ]; then
  echo "First run! Generating a random DB name..."
  PG_DB_NAME="$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c10)"
  set_kv "PG_DB_NAME" "$PG_DB_NAME"
fi

function cleanup {
  if [[ ! -z "${CONTAINER_ID}" ]]; then
    echo "Stopping Postgres DB..."
    docker stop "$CONTAINER_ID"
  fi
  if [[ ! -z "${SOCKET_DIR}" ]]; then
    echo "Removing socket dir..."
    rm -rf "$SOCKET_DIR"
  fi
}
trap cleanup EXIT

SOCKET_DIR="$(mktemp -d)"
chmod 766 "$SOCKET_DIR"
# We put the socket in a subdirectory, because the postgres Docker image takes
# ownership of this directory, and if we own the parent, we can still manage it
# appropriately.
set_kv "SOCKET_DIR" "$SOCKET_DIR/sub"

# We turn off ports and listen solely on our Unix socket.
docker run \
  --name local-postgres \
  --rm \
  --interactive --tty \
  --detach \
  --env POSTGRES_PASSWORD="$PG_PASSWORD" \
  --env POSTGRES_DB="$PG_DB_NAME" \
  --env PG_DATA="$PG_DATA_PATH" \
  --volume "$SOCKET_DIR/sub:/var/run/postgresql" \
  --volume "$PG_LOCAL_DATA_PATH:$PG_DATA_PATH" \
  postgres:14.9 -c listen_addresses=''
CONTAINER_ID=$(docker ps -aqf "name=local-postgres")

echo "Waiting for database to come up..."

# TODO(brandon): Instead of sleeping here, have the migrator poll for a short
# period of time when a --dsn is passed in.
sleep 10

DSN="user=${PG_USER}"
DSN+=" password=${PG_PASSWORD}"
DSN+=" host=${SOCKET_DIR}/sub"
DSN+=" dbname=${PG_DB_NAME}"
DSN+=" sslmode=disable"

if [ "$SKIP_MIGRATE" = false ]; then
  migrate_db "$DSN"
fi

echo "DB is up!"

docker wait local-postgres
