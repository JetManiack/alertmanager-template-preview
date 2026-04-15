# Makefile for alertmanager-template-preview

.PHONY: build test vet fmt build-ui

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -X main.Version=$(VERSION)

build: build-ui
	go build -ldflags "$(LDFLAGS)" -trimpath -o bin/server ./cmd/server

build-ui:
	cd ui && npm install && npm run build
	mkdir -p assets/ui && cp -r ui/dist assets/ui/

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...
