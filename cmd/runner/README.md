# Runner

This directory contains the `runner` binary, which acts as a thin shim around the PACTA portfolio analysis tooling, running tasks created via either Azure Container App Jobs (via the `aztask` package) or local Docker (`dockertask`), loading relevant blobs, and writing relevant outputs.

## Running locally

The `runner` binary doesn't need to be run locally in order to test PACTA processing. By default, the backend API server will execute PACTA runs against a local Docker daemon, testing most of the run-handling code in the process (e.g. file handling, task execution, etc).

If you do want to actually run the full `runner` image on Azure, you can use:

```bash
# Run the backend, tell it to create tasks as real Azure Container Apps Jobs.
bazel run //scripts:run_server -- --use_azure_runner
```

### Creating a new docker image to run locally

When developing the runner, you have two options:

* **Test against local Docker** - Run the server **without** the  `--use_azure_runner`, which means async tasks will run locally, using `docker run ...`. To test local runner changes, you can build and tag a runner image locally with `bazel run //scripts:build_and_load_runner`.
  * After running the script, the updated runner will immediately be available, no need to restart the server.
  * This is the option you'll want to use most of the time.
* **Test against Azure Container Apps Jobs** - Run the server **with** the  `--use_azure_runner`, which means async tasks will be run on Azure, created via the Azure API. To test changes here, you can build and tag a runner image locally with `bazel run //scripts:build_and_load_runner`, and then push it to Azure with `docker push rmisa.azurecr.io/runner:latest`
  * You generally won't need to use this option unless you're testing something very specific about the runner's integration with Azure, as the runner code is identical whether run locally or on Azure.

### Cleaning up old runner containers

By default, we don't auto-remove stopped containers (i.e. finished runner tasks), to give developers a chance to review the logs (e.g. with `docker logs <sha>`). To clean up all completed runs at once, run:

```bash
docker rm $(docker ps -a -q -f "status=exited" -f "ancestor=rmisa.azurecr.io/runner:latest")
```
