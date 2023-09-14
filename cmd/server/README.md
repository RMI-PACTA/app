# PACTA API Server

The PACTA API Server is the main API-serving binary in the PACTA ecosystem. All endpoints are defined in [OpenAPI 3.0](https://spec.openapis.org/oas/v3.0.0) definitions, which live in the [`/openapi`](/openapi) directory. This binary serves the PACTA service (in [pacta.yaml](/openapi/pacta.yaml)), and may add other related services in the future.

## Running the API server

Run the server:

```bash
# Run the backend
bazel run //scripts:run_server
```

There are two ways to access the PACTA API endpoints, both require an RMI-signed JWT token:

1. **With the frontend** - Using the frontend, you can login with Azure and exchange the token for an RMI JWT, which will then be available in the `jwt` cookie in the browser.
2. **With `genjwt`** - The genjwt tool can generate tokens that can be used directly with the PACTA API, make sure to have the `test_server.key` private key generated from your credential service (using the `keygen` tool) in the root of the PACTA repo directory, then run:

```bash
bazel run //scripts:run_genjwt

# This will output something like:
# Token: <header>.<payload>.<sig>
```

You can use this token to query the PACTA (currently just the Petstore example) API:

```bash
APIKEY='<the token from genjwt>'
# Get pets
curl -H "Authorization: BEARER $APIKEY" -X GET localhost:8081/pacta-versions

# []

# Add a pet
curl \
  -H "Authorization: BEARER $APIKEY" \
  -X POST \
  --data '{"name": "Scruffles", "tag": "good dog"}' \
  -H 'Content-Type: application/json' \
  localhost:8081/pets

# {"id":1,"name":"Scruffles","tag":"good dog"}

# Get pets again
curl -H "Authorization: BEARER $APIKEY" -X GET localhost:8081/pets

# [{"id":1,"name":"Scruffles","tag":"good dog"}]
```

## Building and running the Docker container locally

To build and run the image locally:

```bash
# Build the image
bazel build --@io_bazel_rules_go//go/config:pure //cmd/server:image_tarball

# Load it into Docker. This will print out something like:
# Loaded image ID: sha256:<image SHA>
docker load < bazel-bin/cmd/server/image_tarball/tarball.tar

docker run --rm -it sha256:<image SHA from previous step> --config=/configs/local.conf
```

If you get an error like:

```
/server: /lib/x86_64-linux-gnu/libc.so.6: version `GLIBC_2.32' not found (required by /server)
/server: /lib/x86_64-linux-gnu/libc.so.6: version `GLIBC_2.34' not found (required by /server)
```

Make sure you included the `--@io_bazel_rules_go//go/config:pure` flag in `bazel build`, see [`pure` docs](https://github.com/bazelbuild/rules_go/blob/master/go/modes.rst#pure). The problem is that without it, the compiled binary dynamically links glibc against your system, which may use a different version of glibc than the Docker container, which currently uses Debian 11 + glibc 2.28
