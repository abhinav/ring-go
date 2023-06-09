SHELL = /bin/bash

PROJECT_ROOT = $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

# Setting GOBIN and PATH ensures two things:
# - All 'go install' commands we run
#   only affect the current directory.
# - All installed tools are available on PATH
#   for commands like go generate.
export GOBIN = $(PROJECT_ROOT)/bin
export PATH := $(GOBIN):$(PATH)

TEST_FLAGS ?= -v -race

.PHONY: all
all: lint test

.PHONY: lint
lint: golangci-lint tidy-lint

.PHONY: golangci-lint
golangci-lint:
	@echo "[lint] Checking golangci-lint"
	@golangci-lint run

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: tidy-lint
tidy-lint:
	@echo "[lint] Checking go mod tidy"
	@go mod tidy && \
		git diff --exit-code -- go.mod go.sum || \
		(echo "[$(mod)] go mod tidy changed files" && false)

.PHONY: test
test:
	go test $(TEST_FLAGS) ./...

.PHONY: cover
cover:
	go test $(TEST_FLAGS) -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html
