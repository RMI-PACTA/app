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
bazel run //scripts:run_server

# In one last terminal, run the frontend
cd frontend
npm run local
```

## Status

This project is at a very early stage, expect things to change rapidly.

## Security

Please report security issues to security@siliconally.org, or by using one of
the contact methods available on our
[Contact Us page](https://siliconally.org/contact/).

## Contributing

Contribution guidelines can be found [on our website](https://siliconally.org/oss/contributor-guidelines).
