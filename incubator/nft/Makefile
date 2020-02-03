PACKAGES_NOSIMULATION=$(shell go list ./...)
BINDIR ?= $(GOPATH)/bin
SIMAPP = ./app


###############################################################################
###                           Tests & Simulation                            ###
###############################################################################

test: test-unit

test-unit:
	@go test -mod=readonly $(PACKAGES_NOSIMULATION) -tags='ledger test_ledger_mock'

.PHONY: test test-unit

test-sim-nondeterminism:
	@echo "Running non-determinism test..."
	@go test -mod=readonly $(SIMAPP) -run TestAppStateDeterminism -Enabled=true \
		-NumBlocks=100 -BlockSize=200 -Commit=true -Period=0 -v -timeout 24h

test-sim-custom-genesis-fast:
	@echo "Running custom genesis simulation..."
	@echo "By default, ${HOME}/.simapp/config/genesis.json will be used."
	@go test -mod=readonly $(SIMAPP) -run TestFullAppSimulation -Genesis=${HOME}/.simapp/config/genesis.json \
		-Enabled=true -NumBlocks=100 -BlockSize=200 -Commit=true -Seed=99 -Period=5 -v -timeout 24h

test-sim-import-export: runsim
	@echo "Running application import/export simulation. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 50 5 TestAppImportExport

test-sim-after-import: runsim
	@echo "Running application simulation-after-import. This may take several minutes..."
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 50 5 TestAppSimulationAfterImport

test-sim-custom-genesis-multi-seed: runsim
	@echo "Running multi-seed custom genesis simulation..."
	@echo "By default, ${HOME}/.gaiad/config/genesis.json will be used."
	@$(BINDIR)/runsim -Genesis=${HOME}/.simapp/config/genesis.json -SimAppPkg=$(SIMAPP) 400 5 TestFullAppSimulation

test-sim-multi-seed-long: runsim
	@echo "Running long multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 500 50 TestFullAppSimulation

test-sim-multi-seed-short: runsim
	@echo "Running short multi-seed application simulation. This may take awhile!"
	@$(BINDIR)/runsim -Jobs=4 -SimAppPkg=$(SIMAPP) 50 10 TestFullAppSimulation


.PHONY: \
test-sim-nondeterminism \
test-sim-custom-genesis-fast \
test-sim-import-export \
test-sim-after-import \
test-sim-custom-genesis-multi-seed \
test-sim-multi-seed-short \
test-sim-multi-seed-long

lint:
	@echo "--> Running linter"
	@golangci-lint run ./...
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify
.PHONY: lint

###############################################################################
###                          Tools & Dependencies                           ###
###############################################################################

###
# Find OS and Go environment
# GO contains the Go binary
# FS contains the OS file separator
###
ifeq ($(OS),Windows_NT)
  GO := $(shell where go.exe 2> NUL)
  FS := "\\"
else
  GO := $(shell command -v go 2> /dev/null)
  FS := "/"
endif

ifeq ($(GO),)
  $(error could not find go. Is it in PATH? $(GO))
endif

tools: runsim

GOPATH ?= $(shell $(GO) env GOPATH)

TOOLS_DESTDIR  ?= $(GOPATH)/bin
RUNSIM 					= $(TOOLS_DESTDIR)/runsim

runsim: $(RUNSIM)
$(RUNSIM):
	@echo "Installing runsim..."
	@(cd /tmp && go get github.com/cosmos/tools/cmd/runsim@v1.0.0)

tools-clean:
	rm -f $(RUNSIM)
	rm -f tools-stamp