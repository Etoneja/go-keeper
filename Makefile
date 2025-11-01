.PHONY: test build clean lint

test:
	@go test -v ./...

fmt:
	@go fmt ./...

build:
	@go build -o bin/server ./cmd/server/

run:
	@go run ./cmd/my-app/

lint:
	@golangci-lint run

clean:
	@rm -rf bin/

deps:
	@go mod download
	@go mod verify

all: deps test build
