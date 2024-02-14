#!/bin/bash
set -euo pipefail

ROOT="$BUILD_WORKSPACE_DIRECTORY"
cd "$ROOT"

VALID_FLAGS=(
  "with_public_endpoint"
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

if ! [ -x "$(command -v sops)" ]; then
  echo 'Error: sops is not installed.' >&2
  exit 1
fi
if ! [ -x "$(command -v jq)" ]; then
  echo 'Error: jq is not installed.' >&2
  exit 1
fi

SOPS_DATA="$(sops -d "${ROOT}/secrets/local.enc.json")"
LOCAL_DOCKER_CREDS="$(echo $SOPS_DATA | jq .localdocker)"

WEBHOOK_CREDS="$(echo $SOPS_DATA | jq .webhook)"
TOPIC_ID="$(echo $WEBHOOK_CREDS | jq -r .topic_id)"
WEBHOOK_PATH="/events"
WEBHOOK_SHARED_SECRET="$(echo $WEBHOOK_CREDS | jq -r .shared_secret)"

FRP="$(echo $SOPS_DATA | jq .frpc)"
FRP_ADDR="$(echo $FRP | jq -r .addr)"

FRPC_PID=""
function cleanup {
  if [[ ! -z "${FRPC_PID}" ]]; then
    echo "Stopping FRP client/proxy..."
    kill $FRPC_PID
  fi
}
trap cleanup EXIT

eval set --$OPTS
declare -a FLAGS=(
  "--local_docker_tenant_id=$(echo $LOCAL_DOCKER_CREDS | jq -r .tenant_id)"
  "--local_docker_client_id=$(echo $LOCAL_DOCKER_CREDS | jq -r .client_id)"
  "--local_docker_client_secret=$(echo $LOCAL_DOCKER_CREDS | jq -r .password)"
  "--secret_azure_webhook_secrets=${WEBHOOK_SHARED_SECRET}"
)

function create_eventgrid_subscription {
  az eventgrid event-subscription create \
    --name "local-webhook-$1" \
    --source-resource-id "$TOPIC_ID" \
    --endpoint-type=webhook \
    --endpoint="https://$1.${FRP_ADDR}${WEBHOOK_PATH}" \
    --delivery-attribute-mapping Authorization static $WEBHOOK_SHARED_SECRET true
}

EG_SUB_NAME=""
while [ ! $# -eq 0 ]
do
  case "$1" in
    --use_azure_runner)
      FLAGS+=("--use_azure_runner")
      ;;
    --with_public_endpoint)
      if ! [ -x "$(command -v frpc)" ]; then
        echo 'Error: frpc is not installed, cannot run the FRP client/proxy.' >&2
        exit 1
      fi
      SUB_NAME="$2"

      # Check if they already have an Event Grid subscription hooked up to their local env.
      set +e # Don't exit on error, this command might fail if the topic doesn't exist
      az eventgrid event-subscription show \
        --name "local-webhook-${SUB_NAME}" \
        --source-resource-id "$TOPIC_ID"
      SUB_CHECK_EXIT_CODE=$?
      set -e # Back to exiting on error

      if [[ $SUB_CHECK_EXIT_CODE -ne 0 ]]; then
        # The check failed, meaning the webhook doesn't exist.
        # Offer to create it.
        while true; do
          read -p "Should we create a webhook subscription for you (y/n)?" yn
          case $yn in
            [Yy]* ) EG_SUB_NAME="$SUB_NAME"; break;;
            [Nn]* ) break;;
            * ) echo "Please answer yes or no.";;
          esac
        done
      fi

      echo "Running FRP proxy at ${FRP_ADDR}..."
      frpc http \
    		--server_addr="$FRP_ADDR" \
    		--server_port="$(echo $FRP | jq -r .port)" \
    		--token="$(echo $FRP | jq -r .token)" \
        --local_port=8081 \
    		--proxy_name="webhook-$2" \
    		--sd="$2" &
      FRPC_PID=$!
      shift # Extra shift for the subdomain parameter
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

if [[ ! -z "$EG_SUB_NAME" ]]; then
  {
    # We wait because Event Grid requires validating the subscription, which will fail if we aren't running the validation endpoint.
    echo "Waiting for server to start before creating EventGrid subscription"
    sleep 10
    echo "Creating subscription..."

    set +e # Don't exit on error, we can just alert for this
    create_eventgrid_subscription "$EG_SUB_NAME"
    set -e
  } &
fi


bazel run --run_under="cd $ROOT && " //cmd/server -- "${FLAGS[@]}"
