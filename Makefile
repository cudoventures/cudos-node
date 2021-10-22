PACKAGES=$(shell go list ./... | grep -v '/simulation')
PACKAGES_UNITTEST=$(shell go list ./... | grep -v '/simulation' | grep -v '/cli_test')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
SIMAPP = ./simapp
BINDIR ?= $(GOPATH)/bin


VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=cudos-node \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=cudos-noded \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) 

BUILD_FLAGS := -ldflags '$(ldflags)'

all: install test

install: go.sum
		@echo "--> Installing cudos-noded"
		@go install -mod=readonly $(BUILD_FLAGS) ./cmd/cudos-noded

go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify


test: test-unit

test-unit:
	@VERSION=$(VERSION) go test -mod=readonly -ldflags '$(ldflags)' ${PACKAGES_UNITTEST}

test-sim-nondeterminism-fast:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=10 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h
