PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: test
test: lint
	go test $(PKGS)


BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter ./... --vendor


WDIR := $(shell basename $(shell pwd))
BINARY := $(WDIR)
VERSION ?= vlatest
PLATFORMS := linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p bin
	GOOS=$(os) GOARCH=amd64 go build -o bin/$(BINARY)-$(VERSION)-$(os)-amd64

.PHONY: build
build: linux darwin
