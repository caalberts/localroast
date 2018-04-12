.PHONY: all
all: build-deps build fmt vet lint test

ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

APP_MAIN="./cmd/localghost"
APP_EXECUTABLE="localghost"

setup:
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/golang/lint/golint

build-deps:
	dep ensure

update-deps:
	dep ensure

compile:
	go build -o $(APP_EXECUTABLE) $(APP_MAIN)

build: build-deps compile fmt vet lint

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

lint:
	@for p in $(ALL_PACKAGES); do \
		echo "==> Linting $$p"; \
		golint $$p | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; } \
	done

test: build-deps
	ENVIRONMENT=test go test $(ALL_PACKAGES) -race
