SHELL := /bin/bash  # Use bash shell on Unix
TENK_ENV ?= development
GQLGEN_CMD = github.com/99designs/gqlgen

core:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/core.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/core.air.toml
endif

analytics:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/analytics.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/analytics.air.toml
endif

timetrackings:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/timetrackings.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/timetrackings.air.toml
endif

notifications:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/notifications.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/notifications.air.toml
endif

# Flow
# 1. `make run-all`
# 2. `make gateway`
# 3. `make kill-all`
run-all:
	$(MAKE) core & \
	$(MAKE) analytics & \
	$(MAKE) timetrackings & \
	$(MAKE) notifications & \

kill-all:
	npx kill-port 8080 8082 8083 8084 8070

gateway:
	cd services/gateway && npm run start

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

.PHONY: core analytics timetrackings notifications run-all kill-all gateway tidy gqlgen
