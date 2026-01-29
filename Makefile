.PHONY: build test lint fmt clean run

BINARY_NAME=knmi
BUILD_DIR=bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/knmi

build-all:
	./build-all.sh

test:
	go test -v ./...

test-unit:
	go test -v ./tests/unit/...

test-cover:
	go test -coverpkg=./internal/... ./tests/unit/... -coverprofile=coverage.out
	go tool cover -func=coverage.out | tail -1

test-cover-html:
	go test -coverpkg=./internal/... ./tests/unit/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

vet:
	go vet ./...

staticcheck:
	staticcheck ./...

clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

migrate: build
	./$(BUILD_DIR)/$(BINARY_NAME) migrate

sync: build
	./$(BUILD_DIR)/$(BINARY_NAME) sync

all: fmt lint vet test build
