VERSION := $(shell git describe --always --long --dirty)

.DEFAULT_GOAL := build
.PHONY: fmt vet build install

fmt:
	go fmt ./...

vet:
	go vet ./...

build:
	go build -ldflags "-s -w -X 'main.version=${VERSION}'"

install:
	go install -ldflags "-s -w -X 'main.version=${VERSION}'"
	