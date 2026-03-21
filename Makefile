BINARY := comyms

.PHONY: build run test lint format clean

all: build test lint format

build:
	go build -o $(BINARY) .

run:
	go run .

test:
	go test ./...

lint:
	go vet ./...
	golangci-lint run ./...

format:
	gofmt -w .
	npx prettier --write "**/*.md"

clean:
	rm -f $(BINARY)
