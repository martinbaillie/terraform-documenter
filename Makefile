SHELL 	:= /bin/bash
DIR 	:=$(shell pwd)
TIME 	:=$(shell date '+%Y-%m-%dT%T%z')

SRC 		?=$(shell go list ./...)
SRCFILES	?=$(shell find . -type f -name '*.go' -not -path './vendor/*')

LDFLAGS ?=-s -w -extld ld -extldflags -static
FLAGS	?=-a -installsuffix cgo -ldflags "$(LDFLAGS)"
GOOS 	?=darwin windows linux
GOARCH 	?=amd64

PROJECT	?=terraform-documenter

.DEFAULT_GOAL := build

.PHONY: dep clean fix fmt test build help

clean:: ## Removes binary and generated files
	@echo ">> cleaning"
	@rm -f terraform-documenter*

fix:: ## Runs the Golang fix tool
	@echo ">> fixing"
	@go fix $(SRC)

fmt:: ## Formats source code according to the Go standard
	@echo ">> formatting"
	@gofmt -w -s -l $(SRCFILES)

test: clean dep fix fmt 

build:: test ## Builds for all arch ($GOARCH) and OS ($GOOS)
	@echo ">> building"
	@for arch in ${GOARCH}; do \
		for os in ${GOOS}; do \
			echo ">>>> $${os}/$${arch}"; \
			env GOOS=$${os} GOARCH=$${arch} \
			CGO_ENABLED=0 go build $(FLAGS) \
			-o $(PROJECT)-$${os}-$${arch}; \
		done; \
	done

# A help target including self-documenting targets (see the awk statement)
help: ## This help target
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / \
		{printf "\033[36m%-30s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)
