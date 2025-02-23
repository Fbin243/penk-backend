SHELL := /bin/bash  # Use bash shell on Unix
TENK_ENV ?= development
GQLGEN_CMD = github.com/99designs/gqlgen
ENV_FILE ?= .env.$(TENK_ENV)
DB_URI=mongodb+srv://$(MONGO_USER):$(MONGO_PASSWORD)@$(MONGO_ADDRESS)/$(MONGO_DATABASE_NAME)

# export .env file
-include $(ENV_FILE)
export

core:
	@echo "Starting core service..."
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/core.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/core.air.toml
endif

analytic:
	@echo "Starting analytic service..."
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/analytic.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/analytic.air.toml
endif

timetracking:
	@echo "Starting timetracking service..."
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/timetracking.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/timetracking.air.toml
endif

notification:
	@echo "Starting notification service..."
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/notification.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/notification.air.toml
endif

currency:
	@echo "Starting currency service..."
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/currency.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/currency.air.toml
endif

penk:
	@echo "Starting penk service..."
	cd services/penk && npm run start

gateway:
	@echo "Checking if all services are ready..."

	@echo "All services are ready, starting gateway..."
	cd services/gateway && npm run start

dev:
	mkdir -p tmp
	$(MAKE) -j $(SERVICE)

test:
	mkdir -p tmp
	$(MAKE) TENK_ENV=test -j $(SERVICE)

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
		go run -C ./services/$$service/transport $(GQLGEN_CMD); \
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
	@find ./pkg/proto -name "*.proto" -exec \
	protoc --proto_path=./pkg/proto \
	--go_out=./pkg/proto/pb --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/proto/pb --go-grpc_opt=paths=source_relative \
	{} \;

lint:
	@echo "Running linters..."
	@for module in $(shell find . -name 'go.mod' -exec dirname {} \;); do \
		echo "Running linters in $$module"; \
		(cd $$module && golangci-lint run); \
	done

lint-fix:
	@echo "Running linters with fix..."
	@for module in $(shell find . -name 'go.mod' -exec dirname {} \;); do \
		echo "Running linters in $$module"; \
		(cd $$module && golangci-lint run --fix); \
	done

unit-test:
	@echo "Running unit tests..."
	@for module in $(shell find . -name 'go.mod' -exec dirname {} \; | grep -v './test'); do \
		echo "Running unit tests in $$module"; \
		(cd $$module && go test ./... -v) && \
		echo "SUCCESS: Tests passed in $$module" || \
		{ echo "FAIL: Tests failed in $$module"; exit 1; }; \
	done

api-test:
	@echo "Running API tests..."
	@go run cmd/main.go api-test -f profile,character,timetracking,goal

dump:
	@echo "Dumping database $(DB_URI)"
	@mongodump --uri=$(DB_URI)

restore:
	@echo "Restoring database $(DB_URI)"
	@mongorestore --uri=$(DB_URI) dump/$(MONGO_DATABASE_NAME)

.PHONY: core analytic timetracking notification test dump restore
