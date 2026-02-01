LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(PATH):$(LOCAL_BIN)

PKG:=bad-task
SERVER_ENTRYPOINT:=./cmd/server/
SERVER_BIN:=$(LOCAL_BIN)/server

GOLANGCI_LINT_VERSION:=v2.8.0

ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: .install-deps
.install-deps:
	go mod tidy

.PHONY: install-bin
install-bin:
	$(info Installing bin dependencies...)
	@mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: install
install: .install-deps install-bin

.PHONY: run
run:
	@go run $(SERVER_ENTRYPOINT)

.PHONY: clean
clean:
	@rm -i $(LOCAL_BIN)/*

.PHONY: lint
lint:
	@$(LOCAL_BIN)/golangci-lint run ./...
