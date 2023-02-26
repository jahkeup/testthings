all: build test

build: gobuild

test: gotest

lint: golint

fmt: goimports

include ci/go.mk
