BINARY  = socool
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -ldflags "-X main.version=$(VERSION)"

.PHONY: all build install run clean fmt lint test

all: build

build:
	go build $(LDFLAGS) -o $(BINARY) .

install:
	go install $(LDFLAGS) .

run:
	go run .

clean:
	rm -f $(BINARY)

fmt:
	gofmt -w .

lint:
	go vet ./...

test:
	go test ./...
