SHELL := /bin/bash  # Use bash shell on Unix

TENK_ENV ?= development

core:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/core.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/core.air.toml
endif

core_v2:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/core_v2.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/core_v2.air.toml
endif

analytics:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/analytics.air.toml
else
	export TENK_ENV=$(TENK_ENV) && air -c ./tools/air-configs/analytics.air.toml
endif

timetrackings:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=development && air -c ./tools/air-configs/timetrackings.air.toml
else
	export TENK_ENV=development && air -c ./tools/air-configs/timetrackings.air.toml
endif

timetrackings_v2:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=development && air -c ./tools/air-configs/timetrackings_v2.air.toml
else
	export TENK_ENV=development && air -c ./tools/air-configs/timetrackings_v2.air.toml
endif
