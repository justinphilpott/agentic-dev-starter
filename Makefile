VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"

.PHONY: build test clean

build:
	go build $(LDFLAGS) -o seed .

test:
	go test -count=1 ./...

clean:
	rm -f seed
