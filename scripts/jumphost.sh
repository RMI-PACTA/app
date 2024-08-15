#!/bin/bash
set -euo pipefail

ROOT="$BUILD_WORKSPACE_DIRECTORY"
cd "$ROOT"

if [ ! -f ~/.ssh/azure_vm ]; then
  echo "No Azure VM SSH key, creating one"
  ssh-keygen -b 4096 -t rsa -f ~/.ssh/azure_vm -q -N ""
fi

if ! [ -x "$(command -v psql)" ]; then
  echo 'Error: psql is not installed.' >&2
  exit 1
fi

if ! [ -x "$(command -v az)" ]; then
  echo 'Error: The Azure CLI (`az`) is not installed.' >&2
  exit 1
fi

VM_NAME="jumphost-${USER}"
VM_SIZE="Standard_B1ls"
POSTGRES_USER="psqladmin"
POSTGRES_PORT="5432"
case "$1" in
  sa-dev)
    LOCATION="centralus"
    RESOURCE_GROUP="rmi-pacta-dev"
    POSTGRES_HOST="pactadb-dev.postgres.database.azure.com"
    VNET_NAME="pacta-vn-dev"
    SUBNET_NAME="bastion-sn-dev"
    ;;
  *)
    echo "Unknown environment ${1}"
    exit 1
    ;;
esac

echo "Checking for existing jumphost..."

VM_STATE="$(az vm show --name "$VM_NAME" --resource-group "$RESOURCE_GROUP" --query "provisioningState" --output tsv || echo "" 2>/dev/null)"

if [[ $VM_STATE == "Succeeded" ]]; then
  echo "Using existing jumphost"
else

  echo "Creating jumphost VM..."

  # Create a small VM
  az vm create \
    --resource-group "$RESOURCE_GROUP" \
    --name "$VM_NAME" \
    --image Ubuntu2204 \
    --size "$VM_SIZE" \
    --admin-username jumphost \
    --location "$LOCATION" \
    --ssh-key-values @"$HOME/.ssh/azure_vm.pub" \
    --vnet-name "$VNET_NAME" \
    --subnet "$SUBNET_NAME"

  echo "Created jumphost VM, getting IP..."
fi

# Get the public IP address of the VM
PUBLIC_IP=$(az vm show -d --resource-group "$RESOURCE_GROUP" --name "$VM_NAME" --query publicIps -o tsv)

# Set up SSH tunnel to the PostgreSQL instance
printf "Starting tunnel\n"

printf "In a new terminal, connect with:\n\n"
printf "\tpsql --username $POSTGRES_USER --host localhost --port 5433 pactasrv\n"

ssh -N \
  -i ~/.ssh/azure_vm \
  -o ExitOnForwardFailure=yes \
  -L "5433:${POSTGRES_HOST}:${POSTGRES_PORT}" \
  "jumphost@${PUBLIC_IP}"


