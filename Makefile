# Makefile for alertmanager-template-preview

.PHONY: build test vet fmt build-ui

build: build-ui
	go build -o bin/server ./cmd/server

build-ui:
	cd ui && npm install && npm run build
	mkdir -p assets/ui && cp -r ui/dist assets/ui/

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...
