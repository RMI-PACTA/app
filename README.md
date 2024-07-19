# PACTA

This repository contains code for the Paris Agreement Capital Transition Assessment (PACTA) project, which consists of an OpenAPI v3-based API and a Nuxt-based frontend.

## Running

```bash
# First, run a credential service, which you'll need if you want to log in.
# Otherwise, you can manually create a token with genjwt and use the API directly.

cd <path to credential service>

# Run the credential service
bazel run //scripts:run_server -- --use_azure_auth

# In a new terminal, from this directory, run the PACTA database
bazel run //scripts:run_db

# In another terminal, run the PACTA server
bazel run //scripts:run_server -- --with_public_endpoint=$USER

# In one last terminal, run the frontend
cd frontend
npm run local
```

## Status

This project is at a very early stage, expect things to change rapidly.

## Testing the PACTA workflow

To run the PACTA workflow code (e.g. [from this repo](https://github.com/RMI-PACTA/workflow.pacta.webapp)), first create the relevant directories:

```bash
# From the repo root
mkdir workflow-data
cd workflow-data

mkdir -p analysis-output pacta-data real-estate score-card survey benchmarks portfolios report-output summary-output
```

And then load in the relevant files:

* `pacta-data` - Should contain timestamped directories (one per year or quarter or something) that contain the actual data
* `benchmarks` - Should contain timestamped directories containing pre-rendered result sets for comparison to outputs
* `portfolios` - Should contain a single `default_portfolio.csv`, [can be seen here](https://github.com/RMI-PACTA/workflow.pacta.webapp/blob/e02e944b9e94f8af58a83a0210fb0737b9bb908d/tests/portfolios/default_portfolio.csv)

Look at `scripts/run_workflow.sh` for more details. Once all the files are in the correct location, start a run with:

```bash
bazel run //scripts:run_workflow
```

You should see output like:

```
DEBUG [...] Checking configuration.
INFO [...] Running PACTA
INFO [...] Starting portfolio audit
...
```

## Security

Please report security issues to security@siliconally.org, or by using one of
the contact methods available on our
[Contact Us page](https://siliconally.org/contact/).

## Contributing

Contribution guidelines can be found [on our website](https://siliconally.org/oss/contributor-guidelines).
