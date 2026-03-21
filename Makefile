BINARY := google-sheets-mcp

.PHONY: build run test lint clean

build:
	go build -o $(BINARY) .

run:
	go run .

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -f $(BINARY)
