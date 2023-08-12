# PACTA Frontend

This directory contains the web frontend for PACTA, which on [Nuxt 3](https://nuxt.com/), which in turn is built on [Vue 3](https://vuejs.org/).

## Getting Started

[Install npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm), usually via [`nvm`](https://github.com/nvm-sh/nvm), then run `npm install` in this directory.

## Development

To start the development server on http://localhost:3000, run:

```bash
cd frontend
npm run local
```

This will start up a development server, which will hot reload the page as files change.

For typechecking, run:

```bash
npm run typecheck
```

For linting and fixing any fixable lint issues, run:

```bash
# Just lint, no fix
npm run lint

# Lint + fix anything automatically fixable
npm run lint:fix
```

## Compilation 

To build the application for a given environment, run the corresponding `npm run` command:

```bash
npm run build:local
```

## OpenAPI

We use [OpenAPI] to describe the API between this web frontend and the backend PACTA API To take advantage of strong schema typing and TypeScript, we use a [a codegen tool](https://github.com/ferdikoomen/openapi-typescript-codegen), which takes in our [OpenAPI specs](/openapi/) and produces [generated TypeScript files](/frontend/openapi/generated/) 

To update the bindings after changing the OpenAPI schema, run:

```bash
npm run generate:openapi
```
