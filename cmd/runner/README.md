# Runner

This directory contains the `runner` binary, which acts as a thin shim around the PACTA portfolio analysis tooling, running tasks created via either Azure Container App Jobs (via the `aztask` package) or local Docker (`localRunner`), loading relevant blobs, and writing relevant outputs.

## Running locally

The `runner` binary doesn't need to be run locally in order to test PACTA processing. By default, the backend API server will execute PACTA runs against a local Docker daemon, testing most of the run-handling code in the process (e.g. file handling, task execution, etc).

If you do want to actually run the full `runner` image on Azure, you can use:

```bash
# Run the backend, tell it to create tasks as real Azure Container Apps Jobs.
bazel run //scripts:run_server -- --use_azure_runner
```

### Creating a new docker image to run locally

```bash
# Build the runner binary
bazel build --@io_bazel_rules_go//go/config:pure //cmd/runner:image_tarball
# Load the new image into docker, which will output a SHA256 value
docker load < bazel-bin/cmd/runner/image_tarball/tarball.tar
# Tag the runner image in order for it to be picked up locally. Don't push this to the registry!
docker tag <SHA from previous step> rmisa.azurecr.io/runner
```