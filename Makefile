SHELL := /bin/bash  # Use bash shell on Unix

core:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=development && air -c ./tools/air-configs/core.air.toml
else
	export TENK_ENV=development && air -c ./tools/air-configs/core.air.toml
endif

analytics:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=development && air -c ./tools/air-configs/analytics.air.toml
else
	export TENK_ENV=development && air -c ./tools/air-configs/analytics.air.toml
endif

timetrackings:
ifeq ($(OS),Windows_NT)
	set TENK_ENV=development && air -c ./tools/air-configs/timetrackings.air.toml
else
	export TENK_ENV=development && air -c ./tools/air-configs/timetrackings.air.toml
endif
