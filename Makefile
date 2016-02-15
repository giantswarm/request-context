PROJECT=request-context

BUILD_PATH := $(shell pwd)/.gobuild

GS_PATH := "$(BUILD_PATH)/src/github.com/giantswarm"

BIN=request-context

.PHONY:clean get-deps fmt run-tests

GOPATH := $(BUILD_PATH)

SOURCE=$(shell find . -name '*.go')

all: get-deps $(BIN)

ci: clean all run-tests

clean:
	rm -rf $(BUILD_PATH) $(BIN)

get-deps: .gobuild

.gobuild:
	mkdir -p $(GS_PATH)
	cd "$(GS_PATH)" && ln -s ../../../.. $(PROJECT)
	
	# Pin versions of certain libs
	@builder get dep -b b2cb9fa56473e98db8caba80237377e83fe44db5 https://github.com/op/go-logging.git $(GOPATH)/src/github.com/op/go-logging

	#
	# Fetch public dependencies via `go get`
	GOPATH=$(GOPATH) go get -d -v github.com/giantswarm/$(PROJECT)

$(BIN): $(SOURCE)
	GOPATH=$(GOPATH) go build -o $(BIN)

run-tests:
	GOPATH=$(GOPATH) go test ./...

fmt:
	gofmt -l -w .
