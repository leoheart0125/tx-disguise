GOBIN ?= $(shell go env GOBIN)
ifeq ($(GOBIN),)
  GOBIN := $(shell go env GOPATH)/bin
endif

AIR_BIN := $(GOBIN)/air
AIR_PKG := github.com/air-verse/air@latest

GOLANGCI_LINT_BIN := $(GOBIN)/golangci-lint
GOLANGCI_LINT_PKG := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

$(AIR_BIN):
	@go install $(AIR_PKG)

$(GOLANGCI_LINT_BIN):
	@go install $(GOLANGCI_LINT_PKG)

.PHONY: dev
dev: $(AIR_BIN)
	$(AIR_BIN) -c .air.toml

.PHONY: lint
lint: $(GOLANGCI_LINT_BIN)
	$(GOLANGCI_LINT_BIN) fmt
	$(GOLANGCI_LINT_BIN) run --fix
