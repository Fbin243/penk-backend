SHELL := /bin/bash  # Use bash shell on Unix
TENK_ENV ?= development
GQLGEN_CMD = github.com/99designs/gqlgen
ENV_FILE ?= .env.$(TENK_ENV)
DB_URI=mongodb+srv://$(MONGO_USER):$(MONGO_PASSWORD)@$(MONGO_ADDRESS)/$(MONGO_DATABASE_NAME)

# export .env file
-include $(ENV_FILE)

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

payment:
	@echo "Starting payment service..."
	cd services/payment && npm run start

gateway:
	@cd services/gateway && npm run dev

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

# Protocol Buffers Compiler
protoc:
	@echo "Generating protobuf code..."
	@find ./proto -name "*.proto" -exec \
	protoc --proto_path=./proto \
	--go_out=./proto/pb --go_opt=paths=source_relative \
	--go-grpc_out=./proto/pb --go-grpc_opt=paths=source_relative \
	{} \;

protolint:
	@echo "Running protolint..."
	@protolint lint ./proto

protolint-fix:
	@echo "Running protolint with fix..."
	@protolint lint --fix ./proto

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
	@echo "Dumping database [$(DB_URI)]"
	@mongodump --uri=$(DB_URI)

restore:
	@echo "Restoring database [$(DB_URI)]"
	@mongosh $(DB_URI) --eval "db.getSiblingDB('$(MONGO_DATABASE_NAME)').dropDatabase()"
	@mongorestore --uri=$(DB_URI) dump/$(MONGO_DATABASE_NAME)

restore-col:
	@echo "Restoring specific collection [$(COL)] in database [$(DB_URI)]"
	@mongorestore --uri=$(DB_URI) --db $(MONGO_DATABASE_NAME) --collection $(COL)  dump/$(MONGO_DATABASE_NAME)/$(COL).bson --drop

# Docker ----------------------------------------------------------------------
up: 
	@COMPOSE_BAKE=true docker compose -f docker-compose.dev.yml up -d

down: 
	@docker compose -f docker-compose.dev.yml down

down-rmi: 
	@docker compose -f docker-compose.dev.yml down --rmi all

logs:
	@docker compose -f docker-compose.dev.yml logs -f

prune:
	@docker system prune -af

elk-up:
	@docker compose -f elk/docker-compose.dev.yml --env-file .env.development up -d

elk-down:
	@docker compose -f elk/docker-compose.dev.yml --env-file .env.development down
# ------------------------------------------------------------------------------
	
.PHONY: core analytic timetracking notification test dump restore up down logs prune
