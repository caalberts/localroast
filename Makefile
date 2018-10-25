.PHONY: all
all: build fmt vet lint test

APP_MAIN="./"
APP_EXECUTABLE="localroast"

setup:
	go get -u github.com/golang/lint/golint

compile:
	go build -o $(APP_EXECUTABLE) $(APP_MAIN)

build: compile fmt vet lint

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	@for p in $(go list ./...); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done

test:
	go test ./... -race

test-cover:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...