
COMPILER=go
GOBIN=compress


# These are the values we want to pass for VERSION  and BUILD
VERSION=1.0.0
BUILD=`date +%FT%T%z`
# Setup the -Idflags options for go build here,interpolate the variable values
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

build:
	$(COMPILER) build ${LDFLAGS}

clean:
	$(COMPILER) clean

.PHONY: build clean

