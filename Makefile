PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')

test: test-unit

test-unit:
	@go test -mod=readonly $(PACKAGES_NOSIMULATION) -tags='ledger test_ledger_mock'

.PHONY: test test-unit