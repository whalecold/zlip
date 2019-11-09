
COMPILER=go
TARGET=zlip
OUT_DIR=bin
MAIN_DIR=cmd

# Available cpus for compiling, please refer to https://github.com/caicloud/engineering/issues/8186#issuecomment-518656946 for more information.
CPUS ?= $(shell /bin/bash hack/read_cpus_available.sh)

# These are the values we want to pass for VERSION  and BUILD
VERSION=v0.1.1
# Setup the -Idflags options for go build here,interpolate the variable values
LDFLAGS=-ldflags "-X pkg/version/version.VERSION=${VERSION}"

build:
	$(COMPILER) build -o ${OUT_DIR}/${TARGET}  ${LDFLAGS} ${MAIN_DIR}/${TARGET}/main.go

clean:
	rm -rf $(OUT_DIR)/*

test:
	@go test -p $(CPUS) $$(go list ./... | grep -v /vendor | grep -v /test) -coverprofile=coverage.out
	@go tool cover -func coverage.out | tail -n 1 | awk '{ print "Total coverage: " $$3 }'

.PHONY: build clean

