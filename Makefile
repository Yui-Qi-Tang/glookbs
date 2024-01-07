# ENV
BIN=$(notdir $(CURDIR))
BLD=bin


# GO
GO=go
GOOS=GOOS=linux
GOARCH=GOARCH=amd64
GOFLAGS=-ldflags '-w -s'

.PHONY: clean

all: pre_check test test_race build

build:
	$(GOOS) $(GOARCH) $(GO) build $(GOFLAGS) -o $(BLD)/$(BIN)

test:
	$(GO) test ./... -count=1

test_race:
	$(GO) test ./... -race -count=1

pre_check:
	@if ! test -f $(BLD); then mkdir -p $(BLD); fi

clean:
	rm -f $(BLD)/* core

