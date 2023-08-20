.PHONY: build run test lint

build:
	go build

run:
	go run ./...

test:
	go test -race ./...

lint:
	golangci-lint run
