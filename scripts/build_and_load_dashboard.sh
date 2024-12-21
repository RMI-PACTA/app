#!/bin/bash
set -euo pipefail

ROOT="$BUILD_WORKSPACE_DIRECTORY"
cd "$ROOT"

# Build the image
bazel build  --@io_bazel_rules_go//go/config:pure //cmd/dashboard:image_tarball

# Load it into Docker, capture output
LOAD_OUTPUT=$(docker load < bazel-bin/cmd/dashboard/image_tarball/tarball.tar)

# Extract the SHA
IMAGE_ID=$(echo $LOAD_OUTPUT | grep -oP 'sha256:\K\w+')

# Tag the image
docker tag $IMAGE_ID rmisa.azurecr.io/pactadashboard:latest

echo "Tagged $IMAGE_ID as rmisa.azurecr.io/pactadashboard:latest"

