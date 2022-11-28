ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: lint
lint:
	go mod tidy
	gofmt -w -s .

.PHONY: build
build:
	go build -o builds/telescope .
