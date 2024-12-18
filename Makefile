SHELL := /bin/bash  # Use bash shell on Unix
TENK_ENV ?= development
GQLGEN_CMD = github.com/99designs/gqlgen
GATEWAY_DEPENDENCIES = core analytics timetrackings notifications

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

gateway: $(addprefix tmp/, $(GATEWAY_DEPENDENCIES))
	sleep 3
	@echo "Starting gateway service..."
	cd services/gateway && npm run start

dev:
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

.PHONY: core analytics timetrackings notifications gateway dev tidy gqlgen
