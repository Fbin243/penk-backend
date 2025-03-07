# TenK Backend

## Setup

### Environment

1. Create a `.env.development` at the root level
2. Go to our discord server, look for `secret` channel in category `Development`, find and copy the `.env.development` content
3. Paste that content into the `.env.development`

### Air for Live Reload

Check if your machine has installed Air already or not

```sh
air -v
```

- [Install Air on your local machine](https://github.com/cosmtrek/air)

### Make

Check if your machine has installed Make already or not

```sh
make -v
```

- [Install Make for Windows](https://stackoverflow.com/questions/32127524/how-to-install-and-use-make-in-windows)
- [Install Make for MacOS](https://stackoverflow.com/questions/10265742/how-to-install-make-and-gcc-on-a-mac)

### VSCode tools

1. Install `Go` extension
2. Install `gofumpt` in your local machine: `go install mvdan.cc/gofumpt@latest`
3. Press `Ctrl/Cmd + Shift + P`, search `Preferences: Open User Settings (JSON)` and paste these line at the end of the json

```json
"go.useLanguageServer": true,
"gopls": {
 "formatting.gofumpt": true,
},
```

## Development

To run a specific service with Air

```sh
# Run core service with air
make core

# Run analytic service with air
make analytic
```

## MongoDB Migration

Please follow this guide: <https://www.npmjs.com/package/migrate-mongo>

```sh
npm i -g migrate-mongo

# Manage migrations
migrate-mongo up
migrate-mongo down
migrate-mongo status
```

## Mongo backup and restore

Install the Database Tools: <https://www.mongodb.com/docs/database-tools/installation>

```sh
# Backup database
make dump

# Restore database
make restore

# Restore specific collection
make restore-col COL=<colection-name>
```

## Start services, run unit tests and API tests

```sh
# Microservices (Test DB)
make test SERVICE="core analytic timetracking currency penk"

# Microservices (Dev DB)
make dev SERVICE="core analytic timetracking currency penk"

# Gateway
make gateway

# Run unit tests
make unit-test

# Run api tests
make api-test
```

## Linters with golangci-lint and protolint

- Install [golangci-lint](https://golangci-lint.run/welcome/install/) for linting the Golang source code based on the rules defined in `.golangci.yml`

```sh
# Run golangci-lint to check issues for all .go files
make lint

# Run golangci-lint and automatically fix if possible
make lint-fix
```

- Install [protolint](https://github.com/yoheimuta/protolint) for linting Protocol Buffer (.proto) files to ensure they follow best practices.

```sh
# Run protolint to check issues for all .proto files
make protolint

# Run protolint and automatically fix if possible
make protolint-fix
```

## Docker

```sh
# Build docker images from Dockerfile
# Run containers with docker compose
make up

# Stop containers
make down

# Stop containers and remove images
make down-rmi

# Clean up unused images, containers, volumes, and networks
make clean

# Capture logs of containers
make logs
```
