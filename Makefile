PACKAGES_NOSIMULATION=$(shell go list ./...)

test: test-unit

test-unit:
	@go test -mod=readonly $(PACKAGES_NOSIMULATION) -tags='ledger test_ledger_mock'

.PHONY: test test-unit

lint:
	@echo "--> Running linter"
	@golangci-lint run ./...
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify
.PHONY: lint