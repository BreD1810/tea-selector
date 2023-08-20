.PHONY: build run test lint

build:
	go build -o api cmd/tea-selector/main.go

run:
	go run cmd/tea-selector/main.go

test:
	go test -race ./...

lint:
	golangci-lint run
