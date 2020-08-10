COMPILER=go
TARGET=zlip
OUT_DIR=bin
MAIN_DIR=cmd

# Available cpus for compiling, please refer to https://github.com/caicloud/engineering/issues/8186#issuecomment-518656946 for more information.
CPUS ?= $(shell /bin/bash hack/read_cpus_available.sh)
BIN_DIR := $(GOPATH)/bin
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint

# These are the values we want to pass for VERSION  and BUILD
VERSION=v0.1.1
# Setup the -Idflags options for go build here,interpolate the variable values
LDFLAGS=-ldflags "-X pkg/version/version.VERSION=${VERSION}"


# All targets.
.PHONY: lint build clean test


# more info about `GOGC` env: https://github.com/golangci/golangci-lint#memory-usage-of-golangci-lint
lint: $(GOLANGCI_LINT)
	@$(GOLANGCI_LINT) run

$(GOLANGCI_LINT):
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.23.6

build:
	$(COMPILER) build -o ${OUT_DIR}/${TARGET}  ${LDFLAGS} ${MAIN_DIR}/${TARGET}/main.go

clean:
	rm -rf $(OUT_DIR)/*

test:
	@go test -p $(CPUS) $$(go list ./... | grep -v /vendor | grep -v /test) -coverprofile=coverage.out
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'

.PHONY: build clean

