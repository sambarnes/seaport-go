# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# Image URL to use all building/pushing image targets
IMG ?= seaport:latest

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)


# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: all
all: build

##@ Development setup - go

.PHONY: install-dev-tooling-go
install-dev-tooling-go: ## Install dev tooling
	go install github.com/bufbuild/buf/cmd/buf@latest
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

##@ Development

.PHONY: buf
buf: ## Generate golang code from proto files.
	buf generate

.PHONY: tidy
tidy: ## Run go mod tidy.
	go mod tidy

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt buf vet ## Run tests.
	go test ./... -coverprofile cover.out

##@ Build

.PHONY: build
build: fmt buf vet ## Build server binary.
	go build -o bin/main main.go

.PHONY: run
run: fmt buf vet ## Run the server.
	go run main.go
