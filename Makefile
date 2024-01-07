# ENV
BIN=$(notdir $(CURDIR))
BLD=bin


# GO
GO=go
GOOS=GOOS=linux
GOARCH=GOARCH=amd64
GOFLAGS=-ldflags '-w -s'
GOSWAG=$$(go env GOPATH)/bin/swag

.PHONY: clean

all: pre_check test test_race build

build:
	CGO_ENABLED=0 $(GOOS) $(GOARCH) $(GO) build $(GOFLAGS) -o $(BLD)/$(BIN)

test:
	$(GO) test ./... -count=1

test_race:
	$(GO) test ./... -race -count=1

pre_check:
	@if ! test -f $(BLD); then mkdir -p $(BLD); fi

api_doc:
	$(GOSWAG) init -g api/httphandler/*.go

run:
	$(GO) run main.go runserver

clean:
	rm -f $(BLD)/* core

