
COMPILER=go
GOBIN=compress
TARGET=edm
OUT_DIR=bin
MAIN_DIR=cmd

# These are the values we want to pass for VERSION  and BUILD
VERSION=v0.1.1
# Setup the -Idflags options for go build here,interpolate the variable values
LDFLAGS=-ldflags "-X pkg/version/version.VERSION=${VERSION}"

build:
	$(COMPILER) build -o ${OUT_DIR}/${TARGET}  ${LDFLAGS} ${MAIN_DIR}/${TARGET}/main.go

clean:
	rm -rf $(OUT_DIR)/*

.PHONY: build clean

