# General
PWD	:= $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

# Docker
CONTAINER_REGISTRY 	:= docker.io
IMAGE_NAME 			:= r41nwu/telescope

# Git
GIT_TAG 	:= $(shell git describe --always --tags)
GIT_COMMIT 	:= $(shell git rev-parse --short HEAD)

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: lint
lint:
	go mod tidy
	gofmt -w -s .
	docker run \
		--rm \
		-v ${PWD}:/app \
		-w /app \
		golangci/golangci-lint:v1.50.1 \
		golangci-lint run

.PHONY: test
test:
	go test -cover -coverprofile cover.out ./telescope


.PHONY: build
build:
	docker build \
		--platform linux/amd64 \
		--cache-from=${CONTAINER_REGISTRY}/${IMAGE_NAME}:latest \
		-t ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${GIT_COMMIT} \
		.
	docker tag ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${GIT_COMMIT} ${CONTAINER_REGISTRY}/${IMAGE_NAME}:latest

.PHONY: release
release:
	docker pull ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${GIT_COMMIT} || true
	docker tag ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${GIT_COMMIT} ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${GIT_TAG}
	docker push ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${GIT_TAG}

	docker tag ${CONTAINER_REGISTRY}/${IMAGE_NAME}:${GIT_COMMIT} ${CONTAINER_REGISTRY}/${IMAGE_NAME}:latest
	docker push ${CONTAINER_REGISTRY}/${IMAGE_NAME}:latest
