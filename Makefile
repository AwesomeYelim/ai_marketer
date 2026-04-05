.PHONY: build test lint clean run

BINARY=ai-marketer
GO=/usr/local/go/bin/go

build:
	$(GO) build -o $(BINARY) ./...

test:
	$(GO) test ./... -v -race -count=1

test-cover:
	$(GO) test ./... -v -race -coverprofile=coverage.txt -covermode=atomic
	$(GO) tool cover -html=coverage.txt -o coverage.html

lint:
	$(GO) vet ./...

clean:
	rm -f $(BINARY) coverage.txt coverage.html

run:
	$(GO) run main.go run "$(ARGS)"
