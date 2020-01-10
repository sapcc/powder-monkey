DATE      = $(shell date +%Y%m%d%H%M)
VERSION   = v$(DATE)
GOOS      ?= $(shell go env | grep GOOS | cut -d'"' -f2)
GOARCH    ?= $(shell go env GOARCH)
REPO_PATH = github.com/sapcc/powder-monkey

BINARIES := powder-monkey

LDFLAGS := -X github.com/sapcc/powder-monkey/cmd.Version=$(VERSION)

GOFLAGS := -mod=vendor -ldflags "-s -w $(LDFLAGS)"

SRCDIRS  := .
PACKAGES := $(shell find $(SRCDIRS) -type d)
GOFILES  := $(addsuffix /*.go,$(PACKAGES))
GOFILES  := $(wildcard $(GOFILES))

.PHONY: all clean bin/version

all: $(BINARIES:%=bin/$(GOOS)/$(GOARCH)/%) bin/version

bin/$(GOOS)/$(GOARCH)/%: $(GOFILES) Makefile
	GOPATH=$(PWD)/.gopath GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 go build $(GOFLAGS) -v -i -o bin/$(GOOS)/$(GOARCH)/$* $(REPO_PATH)

bin/version:
	echo $(VERSION) > bin/version

clean:
	rm -rf bin/*

.PHONY: vendor
vendor:
	@go mod tidy -v
	@go mod vendor -v
	@go mod verify
