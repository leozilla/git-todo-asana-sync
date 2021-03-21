PROJECT ?= git-todo-asana-sync
IMAGE ?= todosync
BINARY ?= todosync
VERSION ?= $(shell git describe --long --dirty)
LATEST_IMAGE_TAG ?= $(shell docker images | grep "$(PROJECT)/$(IMAGE)" | awk '{print $$2}' | awk 'NR==1')

# directory where this Makefile is located
PROJECT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.DEFAULT_GOAL := help
.PHONY: all help clean lint build deps test test-it run image

all: lint test build image

help: ## (default) prints this help text
	@echo "Supported make tasks:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo "\nSupported parameters:"
	@awk 'BEGIN{FS=".=|##"} /^[a-zA-Z_-]+.*.=.*## / {printf "  \033[32m%-15s\033[0m%-70s\033[33mdefault %s\033[0m\n", $$1,  $$3, $$2}' $(MAKEFILE_LIST)

clean: ## remove compiled binaries
	go clean
	rm -f $(BINARY)
	rm -f $(BINARY)-linux-amd64
	rm -f $(BINARY)-darwin-amd64

LINT := $(shell golangci-lint -v dot 2> /dev/null)

lint: ## runs linters
ifndef LINT
	docker run --rm -v $(SERVICE_DIR):/app -w /app golangci/golangci-lint:v1.34 \
		golangci-lint run -v ./cmd/... ./internal/... ./pkg/... ./test/...
else
	golangci-lint run -v ./cmd/... ./internal/... ./pkg/... ./test/...
endif
	go get -u golang.org/x/lint/golint
	$(shell go list -f {{.Target}} golang.org/x/lint/golint) ./...

build: ## builds the app
	go get github.com/ahmetb/govvv
	govvv build -o $(BINARY) ./cmd/todosync/todosync.go

compile: ## compiles for linux and darwin
	GOOS=linux GOARCH=amd64 go build -o $(BINARY)-linux-amd64 ./cmd/todosync/todosync.go
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY)-darwin-amd64 ./cmd/todosync/todosync.go

deps: ## go gets all deps
	go get -d -v ./...

test: ## runs unit and in-process-component tests
	go test -v -race -coverprofile=coverage.out -run=$(RUN) ./...

ifeq ($(TEST_FILE),)
TEST_FILE := ...
endif

test-it: ## runs integration tests
	cd test/it && go test -v -race -failfast -tags=integration -run=$(RUN) ./$(TEST_FILE)

test-e2e: ## runs e2e tests
	cd test/e2e && go test -v -race -failfast -tags=e2e -run=$(RUN) ./$(TEST_FILE)

run: ## runs the service
	go run ./cmd/todosync/todosync.go

image: ## build docker image
	docker build --cache-from $(PROJECT)/$(IMAGE):latest \
		-t $(PROJECT)/$(IMAGE):$(VERSION) \
		-t $(PROJECT)/$(IMAGE):latest -f build/Dockerfile .