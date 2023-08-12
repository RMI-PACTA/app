#!/bin/bash
set -euo pipefail

KV_FILE=""

function init_kv {
  KV_FILE="$1"
  touch "$KV_FILE"
}

function set_kv {
  if [[ -z "${KV_FILE}" ]]; then
    echo "init_kv was not called"
    exit 1
  fi

  KEY="$1"
  VAL="$2"
  if ! grep -R "^[#]*\s*${KEY}=.*" "$KV_FILE" > /dev/null; then
    # Appending because the key wasn't found
    echo "$KEY=$VAL" >> "$KV_FILE"
  else
    # Replacing because the key was found.
    # We use tildes here and below as our 'regex delimiter' because we expect
    # some values to have slashes in them (e.g. paths), see
    # https://stackoverflow.com/a/27787551 for more info.
    sed -ir "s~^[#]*\s*${KEY}=.*~$KEY=$VAL~" "$KV_FILE"
  fi
}

function get_val {
  if [[ -z "${KV_FILE}" ]]; then
    echo "init_kv was not called"
    exit 1
  fi

  KEY="$1"
  sed -rn "s~^${KEY}=([^\n]+)$~\1~p" $KV_FILE
}
