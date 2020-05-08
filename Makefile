.EXPORT_ALL_VARIABLES:

ELASTICSEARCH_HOST ?= http://localhost:9201
INDEX_PREFIX ?= strPrefix
ELASTICSEARCH_DEBUG ?= false

TEST_ELASTICSEARCH_HOST_v6 ?= http://localhost:9206
TEST_ELASTICSEARCH_HOST_v7 ?= http://localhost:9207

_PRJ_DIR = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))

apply-dry-run:
	go run cmd/main.go apply --dry-run

apply:
	go run cmd/main.go apply

test-all: test-unit test-integration

test-unit:
	CGO_ENABLED=1; go test -coverprofile=coverage-unit.out -race -covermode=atomic ./...

test-integration:
	CGO_ENABLED=1; go test -tags=integration -coverprofile=coverage-integration.out -race -covermode=atomic ./...

lint:
	golangci-lint run -v

qa: lint test-all
