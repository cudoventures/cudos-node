#!/usr/bin/make -f

PACKAGES=$(shell go list ./... | grep -v '/simulation')

VERSION := $(shell echo $(shell git describe --tags))
COMMIT := $(shell git log -1 --format='%H')

BUILDDIR ?= $(CURDIR)/build

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=cudos-node \
	-X github.com/cosmos/cosmos-sdk/version.AppName=cudos-noded \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)'
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf

all: install

install: export CGO_LDFLAGS=-Wl,-rpath,$$ORIGIN/../
install: go.sum
		@echo "--> Installing cudos-noded"
		@go install -mod=readonly $(BUILD_FLAGS) -tags "ledger" ./cmd/cudos-noded


build: export CGO_LDFLAGS=-Wl,-rpath,$$ORIGIN/../
build: go.sum
		@echo "--> Building cudos-noded"
		@go build -mod=readonly $(BUILD_FLAGS) -o $(BUILDDIR)/ -tags "ledger" ./cmd/cudos-noded
		
go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify

test:
	@go test -v -mod=readonly $(PACKAGES)


###############################################################################
###                                Protobuf                                 ###
###############################################################################

CONTAINER_PROTO_VER=v0.7
CONTAINER_PROTO_IMAGE=tendermintdev/sdk-proto-gen:$(CONTAINER_PROTO_VER)
CONTAINER_PROTO_FMT=cosmos-sdk-proto-fmt-$(CONTAINER_PROTO_VER)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(CONTAINER_PROTO_IMAGE) sh ./scripts/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	@if docker ps -a --format '{{.Names}}' | grep -Eq "^${CONTAINER_PROTO_FMT}$$"; then docker start -a $(CONTAINER_PROTO_FMT); else docker run --name $(CONTAINER_PROTO_FMT) -v $(CURDIR):/workspace --workdir /workspace tendermintdev/docker-build-proto \
		find ./proto -name "*.proto" -exec clang-format -i {} \; ; fi
	
proto-lint:
	@$(DOCKER_BUF) lint --error-format=json

proto-check-breaking:
	@$(DOCKER_BUF) breaking --against '$(HTTPS_GIT)#branch=main'

.PHONY: proto-all proto-gen proto-format proto-lint proto-check-breaking 
