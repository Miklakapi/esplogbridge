SHELL := /usr/bin/bash
.ONESHELL:
.SHELLFLAGS := -euo pipefail -c

GO ?= go
BIN ?= esplogbridge
CMD ?= ./cmd/esplogbridge/main.go
OUT ?= ./bin/$(BIN)

.PHONY: help run build clean

help:
	@echo "Targets:"
	@echo "  make run -- [flags]        Run app (passes flags after --)"
	@echo "  make build                 Build binary into ./bin/"
	@echo "  make clean                 Remove ./bin/"
	@echo ""
	@echo "Examples:"
	@echo "  make run -- --help"
	@echo "  make run -- --config ./esplogbridge.yaml --once"
	@echo "  make run -- --dry-run --log-level debug"
	@echo "  make build && ./bin/$(BIN) --help"

run:
	$(GO) run $(CMD) $(filter-out $@,$(MAKECMDGOALS))

build:
	mkdir -p ./bin
	$(GO) build -o $(OUT) $(CMD)

clean:
	rm -rf ./bin

%:
	@: