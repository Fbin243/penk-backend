SHELL := /bin/bash  # Use bash shell on Unix
TENK_ENV ?= development
GQLGEN_CMD = github.com/99designs/gqlgen
PORTS = 8080 8082 8083 8084 8085

core:
	@echo "Starting core service..."
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/core.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/core.air.toml
endif

analytics:
	@echo "Starting analytics service..."
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/analytics.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/analytics.air.toml
endif

timetrackings:
	@echo "Starting timetrackings service..."
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/timetrackings.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/timetrackings.air.toml
endif

notifications:
	@echo "Starting notifications service..."
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/notifications.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/notifications.air.toml
endif

currency:
	@echo "Starting currency service..."
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/currency.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/currency.air.toml
endif

gateway:
	@echo "Checking if all services are ready..."
	@for port in $(PORTS); do \
		while ! curl --silent --fail http://localhost:$$port/health; do sleep 1; done; \
	done
	@echo "All services are ready, starting gateway..."
	cd services/gateway && npm run start

dev:
	mkdir -p tmp
	$(MAKE) -j $(SERVICE)

# Tidy go modules in workspace
tidy:
	@for module in $(shell find . -name 'go.mod' -exec dirname {} \;); do \
		echo "Running go mod tidy in $$module"; \
		(cd $$module && go mod tidy); \
	done

# Generage gqlgen code
gqlgen:
	@echo "Generating gqlgen code for services: $(SERVICE)"
	@for service in $(SERVICE); do \
		echo "Running gqlgen for service: $$service"; \
		go run -C ./services/$$service $(GQLGEN_CMD); \
	done

# Get new JWT token
token:
	@echo "Getting new JWT token..."
	@go run cmd/main.go jwt -u $(UID)

# Import templates from json file to db
templates:
	@echo "Importing templates from file..."
	@go run cmd/main.go import-templates

# Protocol Buffers Compiler
protoc:
	@echo "Generating protobuf code..."
	find ./pkg/proto -name "*.proto" -exec \
	protoc --proto_path=./pkg/proto \
	--go_out=./pkg/proto/pb --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/proto/pb --go-grpc_opt=paths=source_relative \
	{} \;

.PHONY: core analytics timetrackings notifications gateway dev tidy gqlgen token templates
