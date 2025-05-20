.PHONY: build run test clean

APP_NAME=api-gateway

build:
	go build -o $(APP_NAME) main.go

run:
	go run main.go

test:
	go test ./... -cover

clean:
	rm -f $(APP_NAME)

deps:
	go mod tidy
	go mod download

lint:
	golangci-lint run

dev:
	air 