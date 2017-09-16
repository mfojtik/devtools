SOURCEDIR=bz
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=bz

VERSION=0.0.1
BUILD_TIME=`date +%FT%T%z`

LDFLAGS=-ldflags "-X github.com/mfojtik/devtools.version=${VERSION} -X github.com/mfojtik/devtools.buildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
		cd bz && go build ${LDFLAGS} -o ../${BINARY} main.go

.PHONY: install
install:
		cd bz && go install ${LDFLAGS} ./...

.PHONY: clean
clean:
		if [ -f ~/go/bin/${BINARY} ] ; then rm ~/go/bin/${BINARY} ; fi
