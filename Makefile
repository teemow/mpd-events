PROJECT=mpd-events
ORGANIZATION=teemow

SOURCE := $(shell find . -name '*.go')
VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
GOPATH := $(shell pwd)/.gobuild
PROJECT_PATH := $(GOPATH)/src/github.com/$(ORGANIZATION)

.PHONY: all clean run-tests deps bin install

all: deps $(PROJECT)

ci: clean all run-tests

clean:
	rm -rf $(GOPATH) $(PROJECT)

run-tests:
	@GOPATH=$(GOPATH) go test

# deps
deps: .gobuild
.gobuild:
	mkdir -p $(PROJECT_PATH)
	cd $(PROJECT_PATH) && ln -s ../../../.. $(PROJECT)

	@GOPATH=$(GOPATH) go get -d -v github.com/$(ORGANIZATION)/$(PROJECT)

# build
$(PROJECT): $(SOURCE) VERSION
	@echo Building for $(GOOS)/$(GOARCH)
	@GOPATH=$(GOPATH) go build -a -ldflags "-X main.projectVersion=$(VERSION) -X main.projectBuild=$(COMMIT)" -o $(PROJECT)

install: $(PROJECT)
	cp $(PROJECT) /usr/local/bin/

fmt:
	gofmt -l -w .
