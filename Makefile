.PHONY: test build clean lint

BIN_DIR = bin

get_version = $(shell git describe --tags 2>/dev/null || echo "dev")
get_build_time = $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
get_commit = $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")

define get_ldflags
    -X github.com/etoneja/go-keeper/internal/buildinfo.Version=$(get_version) \
    -X github.com/etoneja/go-keeper/internal/buildinfo.BuildTime=$(get_build_time) \
    -X github.com/etoneja/go-keeper/internal/buildinfo.Commit=$(get_commit)
endef

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-version:
	go run cmd/migrate/main.go version

genproto:
	protoc --go_out=. --go-grpc_out=. internal/proto/api.proto --go_opt=default_api_level=API_OPAQUE

genmocks:
	mockgen -source=internal/server/token/interfaces.go -destination=internal/server/token/mocks.go -package=token
	mockgen -source=internal/server/repository/interfaces.go -destination=internal/server/repository/mocks.go -package=repository
	mockgen -source=internal/server/interfaces.go -destination=internal/server/mocks.go -package=server

test:
	@go test -v ./...

test-integration:
	@go test -tags=integration ./internal/server/repository/ -v

fmt:
	@go fmt ./...

build-ctl:
	@mkdir -p bin/
	go build -ldflags="$(get_ldflags)" -o bin/keeperctl ./cmd/ctl

build-server:
	@mkdir -p bin
	go build -ldflags="$(get_ldflags)" -o bin/keepersrv ./cmd/server

build-all: build-ctl build-server

lint:
	@golangci-lint run

clean:
	@rm -rf bin/

deps:
	@go mod download
	@go mod verify

all: deps test build
