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

When testing locally (e.g. without `--use_azure_runner`), you can build and tag a runner image locally and use that. To do that, run `bazel run //scripts:build_and_load_runner`

### Cleaning up old runner containers

By default, we don't auto-remove stopped containers (i.e. finished runner tasks), to give developers a chance to review the logs (e.g. with `docker logs <sha>`). To clean up all completed runs at once, run:

```bash
docker rm $(docker ps -a -q -f "status=exited" -f "ancestor=rmisa.azurecr.io/runner:latest")
```
