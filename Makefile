SHELL = /bin/bash

all: build test

build:
	go build ./...

test:
	go test -race ./...

.PHONY: all build test
